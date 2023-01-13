package raft

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

// ReadyHandler contains a set of operations to be called by raftNode,
// and helps decouple state machine logic from Raft algorithms.
// TODO: add a state machine interface to toApply the commit entries and do snapshot/recover
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

type raftNode struct {
	logger *zap.Logger

	tickMu *sync.Mutex
	raftNodeConfig

	// a chan to send/receive snapshot
	msgSnapCh chan raftpb.Message

	// a chan to send out apply
	applyCh chan toApply

	// a chan to send out readState
	readStateCh chan raft.ReadState

	// utility
	ticker *time.Ticker
	// contention detectors for raft heartbeat message
	td *TimeoutDetector

	stopped chan struct{}
	done    chan struct{}
}

func newRaftNode(cfg raftNodeConfig, logger *zap.Logger) *raftNode {
	raft.SetLogger(newRaftLoggerZap(logger))

	r := &raftNode{
		logger:         logger.With(zap.String("service", "raft")),
		tickMu:         new(sync.Mutex),
		raftNodeConfig: cfg,
		// set up contention detectors for raft heatbeat message.
		// expect to send a heartbeat within 2 heartbeat intervals.
		td:          NewTimeoutDetector(2 * cfg.heartbeat),
		readStateCh: make(chan raft.ReadState, 1),
		msgSnapCh:   make(chan raftpb.Message, maxInflightMsgSnap),
		applyCh:     make(chan toApply),
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
func (r *raftNode) tick() {
	r.tickMu.Lock()
	r.Tick()
	r.tickMu.Unlock()
}

// start prepares and starts raftNode in a new goroutine. It is no longer safe
// to modify the fields after it has been started.
func (r *raftNode) start(rh ReadyHandler) {
	internalTimeout := time.Second

	go func() {
		defer r.onStop()
		isLead := false

		for {
			select {
			case <-r.ticker.C:
				r.tick()

			case rd := <-r.Ready():
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
					r.td.Reset()
				}

				if len(rd.ReadStates) != 0 {
					select {
					case r.readStateCh <- rd.ReadStates[len(rd.ReadStates)-1]:
					case <-time.After(internalTimeout):
						r.logger.Warn("timed out sending read state",
							zap.Duration("timeout", internalTimeout))
					case <-r.stopped:
						return
					}
				}

				notifyCh := make(chan struct{}, 1)
				ap := toApply{
					entries:  rd.CommittedEntries,
					snapshot: rd.Snapshot,
					notifyCh: notifyCh,
				}

				updateCommittedIndex(&ap, rh)

				select {
				case r.applyCh <- ap:
				case <-r.stopped:
					return
				}

				// the leader can write to its disk in parallel with replicating to the
				// followers and them writing to their disks.
				// For more details, check raft thesis 10.2.1
				if isLead {
					r.transport.Send(r.processMessages(rd.Messages))
				}

				// Must save the snapshot file and WAL snapshot entry before saving any other
				// entries or hardstate to ensure that recovery after a snapshot restore is
				// possible.
				if !raft.IsEmptySnap(rd.Snapshot) {
					if err := r.storage.SaveSnap(rd.Snapshot); err != nil {
						r.logger.Fatal("failed to save Raft snapshot", zap.Error(err))
					}
				}

				if err := r.storage.Save(rd.HardState, rd.Entries); err != nil {
					r.logger.Fatal("failed to save raft hard state and entries", zap.Error(err))
				}
				if !raft.IsEmptyHardState(rd.HardState) {
					// TODO: proposalsCommitted metrics
				}

				if !raft.IsEmptySnap(rd.Snapshot) {
					// Force WAL to fsync its hard state before Release() releases
					// old data from the WAL. Otherwise could get an error like:
					// panic: tocommit(107) is out of range [lastIndex(84)]. Was the raft log corrupted, truncated, or lost?
					// See https://github.com/etcd-io/etcd/issues/10219 for more details.
					if err := r.storage.Sync(); err != nil {
						r.logger.Fatal("failed to sync raft snapshot", zap.Error(err))
					}

					// now claim the snapshot has been persisted onto the disk
					notifyCh <- struct{}{}

					r.raftStorage.ApplySnapshot(rd.Snapshot)
					r.logger.Info("applied incoming raft snapshot",
						zap.Uint64("snapshot-index", rd.Snapshot.Metadata.Index))

					if err := r.storage.Release(rd.Snapshot); err != nil {
						r.logger.Fatal("failed to release raft wal", zap.Error(err))
					}
				}

				r.raftStorage.Append(rd.Entries)

				if !isLead {
					// finish processing incoming messages before we signal raftdone chan
					msgs := r.processMessages(rd.Messages)

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
						case <-r.stopped:
							return
						}
					}

					r.transport.Send(msgs)
				} else {
					// leader already processed 'MsgSnap' and signaled
					notifyCh <- struct{}{}
				}

				r.Advance()

			case <-r.stopped:
				return
			}
		}
	}()
}

func updateCommittedIndex(ap *toApply, rh ReadyHandler) {
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

func (r *raftNode) processMessages(msgs []raftpb.Message) []raftpb.Message {
	sentAppResp := false

	for i := len(msgs) - 1; i >= 0; i-- {
		if r.isIDRemoved(msgs[i].To) {
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
			case r.msgSnapCh <- msgs[i]:
			default:
				// drop msgSnap if the inflight chan is full.
			}

			msgs[i].To = 0
		}

		if msgs[i].Type == raftpb.MsgHeartbeat {
			ok, exceed := r.td.Observe(msgs[i].To)
			if !ok {
				// TODO: limit request rate!?
				r.logger.Warn(
					"leader failed to send out heartbeat on time; took too long, leader is overloaded likely from slow disk",
					zap.String("to", fmt.Sprintf("%x", msgs[i].To)),
					zap.Duration("heartbeat-interval", r.heartbeat),
					zap.Duration("expected-duration", 2*r.heartbeat),
					zap.Duration("exceeded-duration", exceed),
				)
			}
		}
	}

	return msgs
}

func (r *raftNode) apply() chan toApply {
	return r.applyCh
}

func (r *raftNode) stop() {
	select {
	case r.stopped <- struct{}{}:
	// Not already stopped, so trigger it
	case <-r.done:
		// has already been stopped - no need to do anything
		return
	}

	// Block until the stop has been acknowledged by start()
	<-r.done
}

func (r *raftNode) onStop() {
	r.Stop()
	r.ticker.Stop()

	if err := r.transport.Stop(); err != nil {
		r.logger.Error("failed to close raft transport", zap.Error(err))
	}

	if err := r.storage.Close(); err != nil {
		r.logger.Panic("failed to close raft storage", zap.Error(err))
	}

	close(r.done)
}

// advanceTicks advances ticks of raft node.
// This can be used for fast-forwarding election
// ticks in multi data-center deployments, thus
// speeding up election process.
func (r *raftNode) advanceTicks(ticks int) {
	for i := 0; i < ticks; i++ {
		r.tick()
	}
}
