package raftstore

import (
	"fmt"
	"github.com/f1shl3gs/manta/raftstore/wal"
	"sync"
	"time"

	"go.etcd.io/raft/v3"
	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
)

const (
	// max number of in-flight snapshot messages server allows to have
	// This number is more than enough for most clusters with 5 machines.
	maxInflightMsgSnap = 16
)

// ReadyHandler contains a set of operations to be called by Node,
// and helps decouple state machine logic from Raft algorithms.
// TODO: add a state machine interface to ToApply the commit entries and do snapshot/recover
type ReadyHandler interface {
	GetLead() uint64
	UpdateLead(uint64)
	UpdateLeadership(bool)
	UpdateCommittedIndex(uint64)
}

type raftNodeConfig struct {
	// to check if msg receiver is removed from cluster
	isIDRemoved func(id uint64) bool
	raft.Node
	storage   *wal.DiskStorage
	heartbeat time.Duration // for logging
	// transport specifies the transport to send and receive msgs to members.
	// Sending messages MUST NOT block. It is okay to drop messages, since
	// clients should timeout and reissue their messages.
	// If transport is nil, server will panic.
	transport Transporter
}

// ToApply contains entries, snapshot to be applied. Once
// an ToApply is consumed, the entries will be persisted to
// to raft storage concurrently; the application must read
// raftDone before assuming the raft messages are stable.
type ToApply struct {
	entries  []raftpb.Entry
	snapshot raftpb.Snapshot
	// notifyCh synchronizes etcd server applies with the raft node
	notifyCh chan struct{}
}

type Node struct {
	logger *zap.Logger

	tickMu *sync.Mutex
	raftNodeConfig

	// a chan to send/receive snapshot
	msgSnapCh chan raftpb.Message

	// a chan to send out apply
	applyCh chan ToApply

	// a chan to send out readState
	readStateCh chan raft.ReadState

	// utility
	ticker *time.Ticker
	// contention detectors for raft heartbeat message
	td *TimeoutDetector

	stopped chan struct{}
	done    chan struct{}
}

func newRaftNode(cfg raftNodeConfig, logger *zap.Logger) *Node {
	raft.SetLogger(newRaftLoggerZap(logger))

	r := &Node{
		logger:         logger.With(zap.String("service", "raft")),
		tickMu:         new(sync.Mutex),
		raftNodeConfig: cfg,
		// set up contention detectors for raft heatbeat message.
		// expect to send a heartbeat within 2 heartbeat intervals.
		td:          NewTimeoutDetector(2 * cfg.heartbeat),
		readStateCh: make(chan raft.ReadState, 1),
		msgSnapCh:   make(chan raftpb.Message, maxInflightMsgSnap),
		applyCh:     make(chan ToApply),
		stopped:     make(chan struct{}),
		done:        make(chan struct{}),
	}

	if r.heartbeat == 0 {
		r.ticker = &time.Ticker{}
	} else {
		r.ticker = time.NewTicker(r.heartbeat)
	}

	return r
}

// raft.Node does not have locks in Raft package
func (n *Node) tick() {
	n.tickMu.Lock()
	n.Tick()
	n.tickMu.Unlock()
}

// start prepares and starts Node in a new goroutine. It is no longer safe
// to modify the fields after it has been started.
func (n *Node) start(rh ReadyHandler) {
	internalTimeout := time.Second

	go func() {
		defer n.onStop()
		isLead := false

		for {
			select {
			case <-n.ticker.C:
				n.tick()

			case rd := <-n.Ready():
				if rd.SoftState != nil {
					newLeader := rd.SoftState.Lead != raft.None && rh.GetLead() != rd.SoftState.Lead
					if newLeader {
						// todo: leader change metric
					}

					/*
											   // TODO: leader state metric
						                    if rd.SoftState.Lead == raft.None {
												hasLeader.Set(0)
											} else {
												hasLeader.Set(1)
						                    }
					*/

					rh.UpdateLead(rd.SoftState.Lead)
					isLead = rd.RaftState == raft.StateLeader
					// TODO: is leader state metric

					rh.UpdateLeadership(newLeader)
					n.td.Reset()
				}

				if len(rd.ReadStates) != 0 {
					select {
					case n.readStateCh <- rd.ReadStates[len(rd.ReadStates)-1]:
					case <-time.After(internalTimeout):
						n.logger.Warn("timed out sending read state",
							zap.Duration("timeout", internalTimeout))
					case <-n.stopped:
						return
					}
				}

				notifyCh := make(chan struct{}, 1)
				ap := ToApply{
					entries:  rd.CommittedEntries,
					snapshot: rd.Snapshot,
					notifyCh: notifyCh,
				}

				updateCommittedIndex(&ap, rh)

				select {
				case n.applyCh <- ap:
				case <-n.stopped:
					return
				}

				// the leader can write to its disk in parallel with replicating to the
				// followers and them writing to their disks.
				// For more details, check raft thesis 10.2.1
				if isLead {
					n.transport.Send(n.processMessages(rd.Messages))
				}

				// Must save the snapshot file and WAL snapshot entry before saving any other
				// entries or hardstate to ensure that recovery after a snapshot restore is
				// possible.
				if !raft.IsEmptySnap(rd.Snapshot) {
                    // we may not need this
                    n.logger.Fatal("save snapshot is not implement")
				}

				if err := n.storage.Save(&rd.HardState, rd.Entries, &rd.Snapshot); err != nil {
					n.logger.Fatal("failed to save raft hard state and entries", zap.Error(err))
				}
				if !raft.IsEmptyHardState(rd.HardState) {
					// TODO: proposalsCommitted metrics
				}

				if !raft.IsEmptySnap(rd.Snapshot) {
					// Force WAL to fsync its hard state before Release() releases
					// old data from the WAL. Otherwise could get an error like:
					// panic: tocommit(107) is out of range [lastIndex(84)]. Was the raft log corrupted, truncated, or lost?
					// See https://github.com/etcd-io/etcd/issues/10219 for more details.
					if err := n.storage.Sync(); err != nil {
						n.logger.Fatal("failed to sync raft snapshot", zap.Error(err))
					}

					// now claim the snapshot has been persisted onto the disk
					notifyCh <- struct{}{}

					n.raftStorage.ApplySnapshot(rd.Snapshot)
					n.logger.Info("applied incoming raft snapshot",
						zap.Uint64("snapshot-index", rd.Snapshot.Metadata.Index))

					if err := n.storage.Release(rd.Snapshot); err != nil {
						n.logger.Fatal("failed to release raft wal", zap.Error(err))
					}
				}

				n.storage.Append(rd.Entries)

				if !isLead {
					// finish processing incoming messages before we signal raftdone chan
					msgs := n.processMessages(rd.Messages)

					// now unblocks 'applyAll' that waits on Raft log disk writes before
					// triggering snapshots
					notifyCh <- struct{}{}

					// Candidate or follower needs to wait for all pending configuration
					// changes to be applied before sending messages.
					// Otherwise we might incorrectly count votes (e.g. votes from removed members).
					// Also slow machine's follower raft-layer could proceed to become the leader
					// on its own single-node cluster, before ToApply-layer applies the config change.
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
						case <-n.stopped:
							return
						}
					}

					n.transport.Send(msgs)
				} else {
					// leader already processed 'MsgSnap' and signaled
					notifyCh <- struct{}{}
				}

				n.Advance()

			case <-n.stopped:
				return
			}
		}
	}()
}

func updateCommittedIndex(ap *ToApply, rh ReadyHandler) {
	var ci uint64

	if len(ap.entries) != 0 {
		ci = ap.entries[len(ap.entries)-1].Index
	}
	if ap.snapshot.Metadata.Index > ci {
		ci = ap.snapshot.Metadata.Index
	}

	if ci != 0 {
		rh.UpdateCommittedIndex(ci)
	}
}

func (n *Node) processMessages(msgs []raftpb.Message) []raftpb.Message {
	sentAppResp := false

	for i := len(msgs) - 1; i >= 0; i-- {
		if n.isIDRemoved(msgs[i].To) {
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
			case n.msgSnapCh <- msgs[i]:
			default:
				// drop msgSnap if the inflight chan is full.
			}

			msgs[i].To = 0
		}

		if msgs[i].Type == raftpb.MsgHeartbeat {
			ok, exceed := n.td.Observe(msgs[i].To)
			if !ok {
				// TODO: limit request rate!?
				n.logger.Warn(
					"leader failed to send out heartbeat on time; took too long, leader is overloaded likely from slow disk",
					zap.String("to", fmt.Sprintf("%x", msgs[i].To)),
					zap.Duration("heartbeat-interval", n.heartbeat),
					zap.Duration("expected-duration", 2*n.heartbeat),
					zap.Duration("exceeded-duration", exceed),
				)
			}
		}
	}

	return msgs
}

func (n *Node) Apply() chan ToApply {
	return n.applyCh
}

func (n *Node) stop() {
	select {
	case n.stopped <- struct{}{}:
	// Not already stopped, so trigger it
	case <-n.done:
		// has already been stopped - no need to do anything
		return
	}

	// Block until the stop has been acknowledged by start()
	<-n.done
}

func (n *Node) onStop() {
	n.Stop()
	n.ticker.Stop()

	if err := n.transport.Stop(); err != nil {
		n.logger.Error("failed to close raft transport", zap.Error(err))
	}

	if err := n.storage.Close(); err != nil {
		n.logger.Panic("failed to close raft storage", zap.Error(err))
	}

	close(n.done)
}

// advanceTicks advances ticks of raft node.
// This can be used for fast-forwarding election
// ticks in multi data-center deployments, thus
// speeding up election process.
func (n *Node) advanceTicks(ticks int) {
	for i := 0; i < ticks; i++ {
		n.tick()
	}
}
