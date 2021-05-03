package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.etcd.io/etcd/pkg/v3/contention"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	"go.uber.org/zap"
)

// apply contains entries, snapshot to be applied. Once
// an apply is consumed, the entries will be persisted to
// to raft storage concurrently; the application must read
// raftDone before assuming the raft messages are stable.
type apply struct {
	entries  []raftpb.Entry
	snapshot raftpb.Snapshot
	// notifyc synchronizes etcd server applies with the raft node
	notifyc chan struct{}
}

type raftNode struct {
	raft.Node

	transport rafthttp.Transport
	// isIDRemoved to check if msg receiver is removed from cluster
	isIDRemoved func(id uint64) bool

	// misc
	heartbeat   time.Duration
	logger      *zap.Logger
	stopped     chan struct{}
	applyC      chan apply
	raftStorage *raft.MemoryStorage
	storage     Storage
	// a chan to send/receive snapshot
	msgSnapC chan raftpb.Message

	tickMtx sync.Mutex
	ticker  *time.Ticker
	// contention detectors for raft heartbeat message
	td *contention.TimeoutDetector

	// a chan to send out readState
	readStateC chan raft.ReadState
	// metrics
	leaderChanges         prometheus.Counter
	hasLeader             prometheus.Gauge
	isLeader              prometheus.Gauge
	heartbeatSendFailures prometheus.Counter
	proposalsCommitted    prometheus.Gauge
}

// raftReadyHandler contains a set of EtcdServer operations to be called by raftNode,
// and helps decouple state machine logic from Raft algorithms.
// TODO: add a state machine interface to apply the commit entries and do snapshot/recover
type raftReadyHandler struct {
	getLead              func() (lead uint64)
	updateLead           func(lead uint64)
	updateLeadership     func(newLeader bool)
	updateCommittedIndex func(uint64)
}

func (r *raftNode) start(rh *raftReadyHandler) {
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
					newLeader := rd.SoftState.Lead != raft.None && rh.getLead() != rd.SoftState.Lead
					if newLeader {
						r.leaderChanges.Inc()
					}

					if rd.SoftState.Lead == raft.None {
						r.hasLeader.Set(0)
					} else {
						r.hasLeader.Set(1)
					}

					rh.updateLead(rd.SoftState.Lead)
					isLead = rd.RaftState == raft.StateLeader
					if isLead {
						r.isLeader.Set(1)
					} else {
						r.isLeader.Set(0)
					}

					rh.updateLeadership(newLeader)
					r.td.Reset()
				}

				if len(rd.ReadStates) != 0 {
					select {
					case r.readStateC <- rd.ReadStates[len(rd.ReadStates)-1]:
					case <-time.After(internalTimeout):
						r.logger.Warn("Timed out sending read state",
							zap.Duration("timeout", internalTimeout))
					case <-r.stopped:
						return
					}
				}

				notifyC := make(chan struct{}, 1)
				ap := apply{
					entries:  rd.CommittedEntries,
					snapshot: rd.Snapshot,
					notifyc:  notifyC,
				}

				updateCommittedIndex(&ap, rh)

				select {
				case r.applyC <- ap:
				case <-r.stopped:
					return
				}

				// the leader can write to its disk in parallel with replication to the
				// followers and them writing to their disks.
				// For more details, check rat thesis 10.2.1
				if isLead {
					r.transport.Send(r.processMessages(rd.Messages))
				}

				// Must save the snapshot file and WAL snapshot entry before saving any
				// other entries or hardstate to ensure that recovery after a snapshot
				// restore is possible
				if !raft.IsEmptySnap(rd.Snapshot) {
					if err := r.storage.SaveSnap(rd.Snapshot); err != nil {
						r.logger.Fatal("failed to save raft snapshot",
							zap.Error(err))
					}
				}

				if err := r.storage.Save(rd.HardState, rd.Entries); err != nil {
					r.logger.Fatal("failed to save Raft hard state and entries",
						zap.Error(err))
				}

				if !raft.IsEmptyHardState(rd.HardState) {
					r.proposalsCommitted.Set(float64(rd.HardState.Commit))
				}

				if !raft.IsEmptySnap(rd.Snapshot) {
					// Force WAL to fsync its hard state before Release() releases
					// old data from the WAL. Otherwise could get an error like:
					// panic: tocommit(107) is out of range [lastIndex(84)]. Was
					// the raft log corrupted, truncated, or lost?
					// See https://github.com/etcd-io/etcd/issues/10219 fore more details.
					if err := r.storage.Sync(); err != nil {
						r.logger.Fatal("failed to sync raft snapshot",
							zap.Error(err))
					}

					// etcdserver now claim the snapshot has been persisted onto the disk
					notifyC <- struct{}{}

					_ = r.raftStorage.ApplySnapshot(rd.Snapshot)
					r.logger.Info("applied incoming raft snapshot",
						zap.Uint64("index", rd.Snapshot.Metadata.Index))

					if err := r.storage.Release(rd.Snapshot); err != nil {
						r.logger.Fatal("failed to release raft wal",
							zap.Error(err))
					}

				}

				r.raftStorage.Append(rd.Entries)

				if !isLead {
					// finish processing incoming messages before we signal raftdone chan
					msgs := r.processMessages(rd.Messages)

					// now unblocks 'applyAll' that waits on Raft log disk writes
					// before triggering snapshots
					notifyC <- struct{}{}

					// Candidate or follower needs to wait for all pending configuration
					// changes to be applied before sending messages.
					// Otherwise we might incorrectly count votes(eg: votes from removed members).
					// Also slow machin's follower raft-layer could proceed to become the leader
					// on its own single-node cluster, before apply-layer applies the config change.
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
						// block until 'applyAll' calls 'applyWait.Trigger'
						// to be in sync with scheduled config-change job
						// (assume notifyC has cap of 1)
						select {
						case notifyC <- struct{}{}:
						case <-r.stopped:
							return
						}
					}

					r.transport.Send(msgs)
				} else {
					// leader already processed 'MsgSnap' and signaled
					notifyC <- struct{}{}
				}

				r.Advance()

			case <-r.stopped:
				return
			}
		}
	}()
}

func updateCommittedIndex(ap *apply, rh *raftReadyHandler) {
	var ci uint64
	if len(ap.entries) != 0 {
		ci = ap.entries[len(ap.entries)-1].Index
	}
	if ap.snapshot.Metadata.Index > ci {
		ci = ap.snapshot.Metadata.Index
	}
	if ci != 0 {
		rh.updateCommittedIndex(ci)
	}
}

func (r *raftNode) processMessages(ms []raftpb.Message) []raftpb.Message {
	sendAppResp := false
	for i := len(ms) - 1; i >= 0; i-- {
		if r.isIDRemoved(ms[i].To) {
			ms[i].To = 0
		}

		if ms[i].Type == raftpb.MsgAppResp {
			if sendAppResp {
				ms[i].To = 0
			} else {
				sendAppResp = true
			}
		}

		if ms[i].Type == raftpb.MsgSnap {
			// There are two separate data stores: the store for v2, and the
			// kv for v3. The msgSnap only contains the most recent snapshot
			// of store without kv. So we need to redirect the msgSnap to
			// etcd server main loop for merging in the current store snapshot
			// and kv snapshot
			select {
			case r.msgSnapC <- ms[i]:
			default:
				// drop msgSnap if the inflight chan is full
			}

			ms[i].To = 0
		}

		if ms[i].Type == raftpb.MsgHeartbeat {
			ok, exceed := r.td.Observe(ms[i].To)
			if !ok {
				// TODO: limit request rate
				r.logger.Warn("leader failed to send out heartbeat on time; took too long, leader is overloaded likely from slow disk",
					zap.String("to", fmt.Sprintf("%x", ms[i].To)),
					zap.Duration("heartbeat-interval", r.heartbeat),
					zap.Duration("expected-duration", 2*r.heartbeat),
					zap.Duration("exceeded-duration", exceed))
			}

			r.heartbeatSendFailures.Inc()
		}
	}

	return ms
}

func (r *raftNode) onStop() {

}

func (r *raftNode) tick() {
	r.tickMtx.Lock()
	r.Tick()
	r.tickMtx.Unlock()
}

func (r *raftNode) apply() chan apply {
	return r.applyC
}
