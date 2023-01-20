package raftstore

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/f1shl3gs/manta/raftstore/pb"
	"github.com/f1shl3gs/manta/raftstore/wal"

	"go.etcd.io/raft/v3"
	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
)

const (
	heartbeat = time.Duration(100) * time.Millisecond

	readIndexRetryTime = 500 * time.Millisecond
)

var (
	ErrLeaderChanged = errors.New("leader changed")

	ErrTimeout = errors.New("request timed out")

	ErrStopped = errors.New("server stopped")
)

// toApply contains entries, snapshot to be applied. Once
// an toApply is consumed, the entries will be persisted to
// to raft storage concurrently; the application must read
// raftDone before assuming the raft messages are stable.
type toApply struct {
	entries  []raftpb.Entry
	snapshot raftpb.Snapshot
	// notifyCh synchronizes etcd server applies with the raft node
	notifyCh chan struct{}
}

// advanceTicks advances ticks of raft node.
// This can be used for fast-forwarding election
// ticks in multi data-center deployments, thus
// speeding up election process.
func (s *Store) advanceTicks(ticks int) {
	for i := 0; i < ticks; i++ {
		s.tick()
	}
}

// raft.Node does not have locks in Raft package
func (s *Store) tick() {
	s.tickMu.Lock()
	s.raftNode.Tick()
	s.tickMu.Unlock()
}

func (s *Store) checkSnapshot(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	lastSnapshot := time.Now()
	snapshotAfterEntries := uint64(10000)
	// snapshotFrequence = 30 * time.Minute

	propagateSnapshot := func(index uint64) {
		pCtx, cancel := context.WithTimeout(ctx, s.reqTimeout())
		err := s.propose(pCtx, pb.InternalRequest{
			Snapshot: &pb.Snapshot{
				Index: index,
			},
		})
		cancel()

		if err != nil {
			s.logger.Error("propose snapshot failed",
				zap.Error(err))
		} else {
			lastSnapshot = time.Now()

			s.logger.Debug("propagate snapshot success",
				zap.Time("lastTime", lastSnapshot),
				zap.Uint64("toIndex", index))
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		applied := s.appliedIndex.Load()
		s.raftStorage.SetUint(wal.CheckpointIndex, applied)

		if s.lead.Load() != s.self.ID {
			continue
		}

		// If leader doesn't have a snapshot, we should create one immediately.
		// This is very useful when you bring up the cluster. If you remove or
		// ad a node, the new follower won't get a snapshot if the leader doesn't
		// have one.
		snap, err := s.raftStorage.Snapshot()
		if err != nil {
			s.logger.Error("while retrieving snapshot from raft storage",
				zap.Error(err))
			continue
		}

		// If we don't have a snapshot, or if there are too many log files in
		// Raft.
		if raft.IsEmptySnap(snap) || s.raftStorage.NumLogFiles() > 4 {
			propagateSnapshot(applied)
			continue
		}

		first, err := s.raftStorage.FirstIndex()
		if err == nil {
			applied := s.appliedIndex.Load()
			if applied-snapshotAfterEntries > first {
				applied = applied - snapshotAfterEntries
			}
		}

		first, err = s.raftStorage.FirstIndex()
		if err == nil {
			if applied-first > snapshotAfterEntries {
				to := first + snapshotAfterEntries
				s.logger.Info("trying to take snapshot",
					zap.Uint64("first", first),
					zap.Uint64("applied", applied),
					zap.Uint64("to", to),
					zap.Uint64("snapAfterEntries", applied-first-snapshotAfterEntries))
				propagateSnapshot(to)
			}
		}
	}
}

// startRaftLoop prepares and starts Node in a new goroutine. It is no longer safe
// to modify the fields after it has been started.
func (s *Store) startRaftLoop() {
	isLead := false
	internalTimeout := time.Second

	defer s.onStop()

	for {
		select {
		case <-s.ticker.C:
			s.tick()

		case rd := <-s.raftNode.Ready():
			if rd.SoftState != nil {
				newLeader := rd.SoftState.Lead != raft.None && s.getLead() != rd.SoftState.Lead
				if newLeader {
					s.leaderChanges.Inc()
				}

				// TODO: leader state metric
				if rd.SoftState.Lead == raft.None {
					s.hasLeader.Set(0)
				} else {
					s.hasLeader.Set(1)
				}

				s.updateLead(rd.SoftState.Lead)
				isLead = rd.RaftState == raft.StateLeader
				if isLead {
					s.isLeader.Set(1)
				} else {
					s.isLeader.Set(0)
				}

				s.updateLeadership(newLeader)
				s.td.Reset()
			}

			if len(rd.ReadStates) != 0 {
				select {
				case s.readStateCh <- rd.ReadStates[len(rd.ReadStates)-1]:
				case <-time.After(internalTimeout):
					s.logger.Warn("timed out sending read state",
						zap.Duration("timeout", internalTimeout))
				case <-s.stopped:
					return
				}
			}

			notifyCh := make(chan struct{}, 1)
			ap := toApply{
				entries:  rd.CommittedEntries,
				snapshot: rd.Snapshot,
				notifyCh: notifyCh,
			}

			// update commited index
			{
				var ci uint64
				if len(ap.entries) != 0 {
					ci = ap.entries[len(ap.entries)-1].Index
				}
				if ap.snapshot.Metadata.Index > ci {
					ci = ap.snapshot.Metadata.Index
				}
				if ci != 0 {
					s.updateCommittedIndex(ci)
				}
			}

			select {
			case s.applyCh <- ap:
			case <-s.stopped:
				return
			}

			// the leader can write to its disk in parallel with replicating to the
			// followers and them writing to their disks.
			// For more details, check raft thesis 10.2.1
			if isLead {
				s.transport.Send(s.processMessages(rd.Messages))
			}

			// Must save the snapshot file and WAL snapshot entry before saving any other
			// entries or hardstate to ensure that recovery after a snapshot restore is
			// possible.
			if !raft.IsEmptySnap(rd.Snapshot) {
				// we may not need this
				s.logger.Fatal("save snapshot is not implement")
			}

			if err := s.raftStorage.Save(&rd.HardState, rd.Entries, &rd.Snapshot); err != nil {
				s.logger.Fatal("failed to save raft hard state and entries", zap.Error(err))
			}
			if !raft.IsEmptyHardState(rd.HardState) {
				// TODO: proposalsCommitted metrics
			}

			if !raft.IsEmptySnap(rd.Snapshot) {
				// Force WAL to fsync its hard state before Release() releases
				// old data from the WAL. Otherwise could get an error like:
				// panic: tocommit(107) is out of range [lastIndex(84)]. Was the raft log corrupted, truncated, or lost?
				// See https://github.com/etcd-io/etcd/issues/10219 for more details.
				if err := s.raftStorage.Sync(); err != nil {
					s.logger.Fatal("failed to sync raft snapshot", zap.Error(err))
				}

				// now claim the snapshot has been persisted onto the disk
				notifyCh <- struct{}{}

				/*
				   TODO:

				   s.storage.ApplySnapshot(rd.Snapshot)
				   s.logger.Info("applied incoming raft snapshot",
				   zap.Uint64("snapshot-index", rd.Snapshot.Metadata.Index))

				   if err := s.storage.Release(rd.Snapshot); err != nil {
				   s.logger.Fatal("failed to release raft wal", zap.Error(err))
				   }
				*/
			}

			// s.storage.Append(rd.Entries)

			if !isLead {
				// finish processing incoming messages before we signal raftdone chan
				msgs := s.processMessages(rd.Messages)

				// now unblocks 'applyAll' that waits on Raft log disk writes before
				// triggering snapshots
				notifyCh <- struct{}{}

				// Candidate or follower needs to wait for all pending configuration
				// changes to be applied before sending messages.
				// Otherwise we might incorrectly count votes (e.g. votes from removed members).
				// Also slow machine's follower raft-layer could proceed to become the leader
				// on its own single-node cluster, before toApply-layer applies the config change.
				// We simply wait for ALL pending entries to be applied for now.
				// We might improve this later on if it causes unnecessary long blocking issues.
				waitApply := false
				for _, ent := range rd.CommittedEntries {
					if ent.Type == raftpb.EntryConfChange {
						waitApply = true
						break
					}
				}
				if waitApply {
					// blocks until 'applyAll' calls 'applyWait.Trigger'
					// to be in sync with scheduled config-change job
					// (assume notifyCh has cpa of 1)
					select {
					case notifyCh <- struct{}{}:
					case <-s.stopped:
						return
					}
				}

				s.transport.Send(msgs)
			} else {
				// leader already processed 'MsgSnap' and signaled
				notifyCh <- struct{}{}
			}

			s.raftNode.Advance()

		case <-s.stopped:
			return
		}
	}
}

func (s *Store) updateLead(lead uint64) {
	s.lead.Store(lead)
}

func (s *Store) getLead() uint64 {
	return s.lead.Load()
}

func (s *Store) updateLeadership(newLeader bool) {
	// we can start or stop some backgroud service here
	s.logger.Info("update leadership", zap.Bool("new leader", newLeader), zap.String("id", strconv.FormatUint(s.getLead(), 16)))
}

func (s *Store) updateCommittedIndex(ci uint64) {
	cci := s.committedIndex.Load()
	if ci > cci {
		s.committedIndex.Store(ci)
	}
}

func (s *Store) isIDRemoved(id uint64) bool {
	return s.transport.Removed(id)
}

func (s *Store) processMessages(msgs []raftpb.Message) []raftpb.Message {
	sentAppResp := false

	for i := len(msgs) - 1; i >= 0; i-- {
		if s.isIDRemoved(msgs[i].To) {
			msgs[i].To = 0
		}

		if msgs[i].Type == raftpb.MsgAppResp {
			if sentAppResp {
				msgs[i].To = 0
			} else {
				sentAppResp = true
			}
		}

		if msgs[i].Type == raftpb.MsgSnap {
			// The msgSnap only contains the most recent snapshot of store without KV.
			// So we need to redirect the msgSNap to etcd server main loop for mergin
			// in the current store snapshot and KV snapshot.
			select {
			case s.msgSnapCh <- msgs[i]:
			default:
				// drop msgSnap if the inflight chan is full.
			}

			msgs[i].To = 0
		}

		if msgs[i].Type == raftpb.MsgHeartbeat {
			ok, exceed := s.td.Observe(msgs[i].To)
			if !ok {
				// TODO: limit request rate!?
				s.logger.Warn(
					"leader failed to send out heartbeat on time; took too long, leader is overloaded likely from slow disk",
					zap.Uint64("to", msgs[i].To),
					zap.Duration("heartbeat-interval", heartbeat),
					zap.Duration("expected-duration", 2*heartbeat),
					zap.Duration("exceeded-duration", exceed),
				)
			}
		}
	}

	return msgs
}

func uint64ToBigEndianBytes(number uint64) []byte {
	byteResult := make([]byte, 8)
	binary.BigEndian.PutUint64(byteResult, number)
	return byteResult
}

func (s *Store) reqTimeout() time.Duration {
	// 5s for queue waiting, computation and disk IO delay
	// + 2 * election timeout for possible leader election
	return 5*time.Second + 2*time.Duration(electionTick*tickMs)*time.Millisecond
}

func (s *Store) sendReadIndex(ctx context.Context, reqIndex uint64) error {
	toSend := uint64ToBigEndianBytes(reqIndex)

	ctx, cancel := context.WithTimeout(ctx, s.reqTimeout())
	err := s.raftNode.ReadIndex(ctx, toSend)
	cancel()

	if err == ErrStopped {
		return err
	}

	if err != nil {
		s.readIndexFailed.Inc()
		s.logger.Warn("failed to get read index from Raft",
			zap.Error(err))
	}

	return err
}

func (s *Store) requestCurrentIndex(ctx context.Context, leaderChangeNotifier <-chan struct{}, reqID uint64) (uint64, error) {
	err := s.sendReadIndex(ctx, reqID)
	if err != nil {
		return 0, err
	}

	errTimer := time.NewTimer(s.reqTimeout())
	defer errTimer.Stop()
	retryTimer := time.NewTimer(readIndexRetryTime)
	defer retryTimer.Stop()

	firstCommitInTermNotifier := s.firstCommitInTerm.receive()

	for {
		select {
		case rs := <-s.readStateCh:
			reqIdBytes := uint64ToBigEndianBytes(reqID)
			gotOwnResp := bytes.Equal(rs.RequestCtx, reqIdBytes)
			if !gotOwnResp {
				// a previous request might time out. now we should ignore the
				// response of it and continue waiting for the response of the
				// current requests.
				respId := uint64(0)
				if len(rs.RequestCtx) == 8 {
					respId = binary.BigEndian.Uint64(rs.RequestCtx)
				}

				s.logger.Warn("ignored out-of-date read index response; "+
					"local node read indexes queueing up and waitting to be in sync with leader",
					zap.Uint64("sent-request-id", reqID),
					zap.Uint64("received-request-id", respId))
				s.slowReadInex.Inc()
				continue
			}

			return rs.Index, nil

		case <-leaderChangeNotifier:
			s.readIndexFailed.Inc()
			return 0, ErrLeaderChanged

		case <-firstCommitInTermNotifier:
			firstCommitInTermNotifier = s.firstCommitInTerm.receive()
			s.logger.Info("first commit in current term: resending ReadIndex request")
			err := s.sendReadIndex(ctx, reqID)
			if err != nil {
				return 0, err
			}

			retryTimer.Reset(readIndexRetryTime)
			continue

		case <-retryTimer.C:
			s.logger.Warn("waiting for ReadIndex response took too long, retrying",
				zap.Uint64("sent-request-id", reqID),
				zap.Duration("retry-timeout", readIndexRetryTime))

			err := s.sendReadIndex(ctx, reqID)
			if err != nil {
				return 0, err
			}

			retryTimer.Reset(readIndexRetryTime)
			continue

		case <-errTimer.C:
			s.logger.Warn("timed out waiting for read inex response (local node might have slow network)",
				zap.Duration("timeout", s.reqTimeout()))
			s.slowReadInex.Inc()
			return 0, ErrTimeout

		case <-s.stopped:
			return 0, ErrStopped
		}
	}
}

func (s *Store) linearizableReadLoop(ctx context.Context) {
	for {
		reqID := s.idGen.Next()
		leaderChangeNotifier := s.leaderChanged.receive()

		select {
		case <-ctx.Done():
			return
		case <-leaderChangeNotifier:
			continue
		case <-s.readWaitCh:
		}

		nextNr := newErrNotifier()
		s.readMtx.Lock()
		nr := s.readNotifier
		s.readNotifier = nextNr
		s.readMtx.Unlock()

		confirmedIndex, err := s.requestCurrentIndex(ctx, leaderChangeNotifier, reqID)
		if err == ErrStopped {
			return
		}
		if err != nil {
			nr.notify(err)
			continue
		}

		appliedInex := s.appliedIndex.Load()
		if appliedInex < confirmedIndex {
			select {
			case <-s.applyWait.Wait(confirmedIndex):
			case <-s.stopped:
				return
			}
		}

		nr.notify(nil)
	}
}

type errNotifier struct {
	ch  chan struct{}
	err error
}

func (n *errNotifier) notify(err error) {
	n.err = err
	close(n.ch)
}

func newErrNotifier() *errNotifier {
	return &errNotifier{
		ch:  make(chan struct{}),
		err: nil,
	}
}

// notify is a thread safe struct that can be used to send notification
// about some event to multiple consumer.
type notifier struct {
	mtx sync.RWMutex
	ch  chan struct{}
}

func newNotifier() *notifier {
	return &notifier{
		ch: make(chan struct{}),
	}
}

// receive returns channel that can be used to wait for notification.
// consumers will be informed by closing the channel.
func (n *notifier) receive() <-chan struct{} {
	n.mtx.RLock()
	ch := n.ch
	n.mtx.RUnlock()

	return ch
}

// notifiy closes the channel passed to consumers and creates new
// channel to used for next notification.
func (n *notifier) notify() {
	newCh := make(chan struct{})
	n.mtx.Lock()
	toClose := n.ch
	n.ch = newCh
	n.mtx.Unlock()

	close(toClose)
}

func (s *Store) linearizableReadNotify(ctx context.Context) error {
	s.readMtx.RLock()
	nc := s.readNotifier
	s.readMtx.RUnlock()

	// signal linearizable loop for current notify if it hasn't been already
	select {
	case s.readWaitCh <- struct{}{}:
	default:
	}

	// wait for read state notification
	select {
	case <-nc.ch:
		return nc.err
	case <-ctx.Done():
		return ctx.Err()
	case <-s.done:
		return ErrStopped
	}
}

// ReadyNotify returns a channel that will be closed when the
// server is ready to serve client requests.
func (s *Store) readyNotify() <-chan struct{} {
	return s.readyCh
}

func (s *Store) stop() {
	select {
	case s.stopped <- struct{}{}:
	// Not already stopped, so trigger it
	case <-s.done:
		// has already been stopped - no need to do anything
		return
	}

	// Block until the stop has been acknowledged by start()
	<-s.done
}

func (s *Store) onStop() {
	s.raftNode.Stop()
	s.ticker.Stop()
	s.transport.Stop()
	db := s.db.Swap(nil)
	_ = db.Close()

	if err := s.raftStorage.Close(); err != nil {
		s.logger.Panic("failed to close raft storage", zap.Error(err))
	}

	close(s.done)
}

func (s *Store) adjustTicks() {
	members := len(s.transport.Peers())

	// single-node fresh start, or single-node recovers from snapshot
	if members == 1 {
		ticks := electionTick - 1
		s.logger.Info("started as single-node, fast-forwarding election ticks", zap.Int("forward-ticks", ticks))

		s.advanceTicks(ticks)
		return
	}

	// retry up to "transport.connReadTimeout", which is 5sec
	// until peer connection reports; otherwise:
	//
	waitTime := 5 * time.Second
	itv := 50 * time.Millisecond
	for i := int64(0); i < int64(waitTime/itv); i++ {
		select {
		case <-time.After(itv):
		case <-s.stopped:
			return
		}

		// todo: this value should be the active one
		actives := len(s.transport.Peers())
		if actives > 1 {
			// multi-node received peer connection reports
			// adjust ticks, in case slow leader message receive
			ticks := electionTick - 2
			s.logger.Info("initialized peer connections, fast-forwarding election ticks",
				zap.Int("forward-ticks", ticks))

			s.advanceTicks(ticks)
			return
		}
	}
}
