package raftstore

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/f1shl3gs/manta/pkg/fsutil"
	"github.com/f1shl3gs/manta/raftstore/internal"
	"github.com/f1shl3gs/manta/raftstore/membership"
	"github.com/f1shl3gs/manta/raftstore/rawkv"
	"github.com/f1shl3gs/manta/raftstore/transport"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.etcd.io/etcd/pkg/v3/contention"
	"go.etcd.io/etcd/pkg/v3/idutil"
	"go.etcd.io/etcd/pkg/v3/pbutil"
	"go.etcd.io/etcd/pkg/v3/wait"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.etcd.io/etcd/server/v3/wal"
	"go.uber.org/zap"
)

const (
	recommendedMaxRequestBytes = 10 * 1024 * 1024

	readIndexRetryTime = 500 * time.Millisecond

	DefaultSnapshotCount = 100000

	// DefaultSnapshotCatchUpEntries is the number of entries for a slow follower
	// to catch-up after compacting the raft storage entries.
	// We expect the follower has a millisecond level latency with the leader.
	// The max throughput is around 10K. Keep a 5K entries is enough for helping
	// follower to catch up.
	DefaultSnapshotCatchUpEntries uint64 = 5000

	// In the health case, there might be a small gap (10s of entries) between
	// the applied index and committed index.
	// However, if the committed entries are very heavy to apply, the gap might grow.
	// We should stop accepting new proposals if the gap growing to a certain point.
	maxGapBetweenApplyAndCommitIndex = 5000
)

const (
	// The max throughput of etcd will not exceed 100MB/s (100K * 1KB value).
	// Assuming the RTT is around 10ms, 1MB max size is large enough.
	maxSizePerMsg = 1 * 1024 * 1024
	// Never overflow the rafthttp buffer, which is 4096.
	// TODO: a better const?
	maxInflightMsgs = 4096 / 8
)

var (
	ErrUnknownMethod                 = errors.New("unknown method")
	ErrStopped                       = errors.New("server stopped")
	ErrCanceled                      = errors.New("request cancelled")
	ErrTimeout                       = errors.New("request timed out")
	ErrTimeoutDueToLeaderFail        = errors.New("request timed out, possibly due to previous leader failure")
	ErrTimeoutDueToConnectionLost    = errors.New("request timed out, possibly due to connection lost")
	ErrTimeoutLeaderTransfer         = errors.New("request timed out, leader transfer took too long")
	ErrLeaderChanged                 = errors.New("leader changed")
	ErrNotEnoughStartedMembers       = errors.New("re-configuration failed due to not enough started members")
	ErrLearnerNotReady               = errors.New("can only promote a learner member which is in sync with leader")
	ErrNoLeader                      = errors.New("no leader")
	ErrNotLeader                     = errors.New("not leader")
	ErrRequestTooLarge               = errors.New("request is too large")
	ErrNoSpace                       = errors.New("no space")
	ErrTooManyRequests               = errors.New("too many requests")
	ErrUnhealthy                     = errors.New("unhealthy cluster")
	ErrKeyNotFound                   = errors.New("key not found")
	ErrCorrupt                       = errors.New("corrupt cluster")
	ErrBadLeaderTransferee           = errors.New("bad leader transferee")
	ErrClusterVersionUnavailable     = errors.New("cluster version not found during downgrade")
	ErrWrongDowngradeVersionFormat   = errors.New("wrong downgrade target version format")
	ErrInvalidDowngradeTargetVersion = errors.New("invalid downgrade target version")
	ErrDowngradeInProcess            = errors.New("cluster has a downgrade job in progress")
	ErrNoInflightDowngrade           = errors.New("no inflight downgrade job")
	ErrRejectFromRemovedMember       = errors.New("cannot process message from removed member")
)

type Store struct {
	id uint64

	reqIDGen *idutil.Generator
	stopCh   chan struct{}

	// read routine notifies server that it waits for reading by
	// sending an empty struct to readWaitCh
	readWaitCh chan struct{}
	readMtx    sync.RWMutex
	// readNotifier is used to notify the read routine that it can
	// process the request when there is no error
	readNotifier *notifier
	wal          *wal.WAL

	config *Config
	logger *zap.Logger

	// raft stuff
	applyCh     chan apply
	errCh       chan error
	msgSnapCh   chan raftpb.Message
	raftNode    raft.Node
	raftStorage *raft.MemoryStorage
	// todo: temporary interface
	storage              Storage
	readStateC           chan raft.ReadState
	firstCommitInTermMtx sync.RWMutex
	firstCommitInTermCh  chan struct{}
	// leaderChangedCh is used to notify the linearizable read loop to drop the old read requests.
	leaderChangedCh  chan struct{}
	leaderChangedMtx sync.RWMutex
	tickMtx          sync.Mutex
	ticker           *time.Ticker
	// contention detectors for raft heartbeat message
	heartbeat    time.Duration // for logging?
	td           *contention.TimeoutDetector
	readyHandler ReadyHandler
	cluster      membership.Cluster
	transporter  transport.Transporter
	wg           sync.WaitGroup
	// atomic states
	appliedIndex      uint64
	committedIndex    uint64
	term              uint64
	inflightSnapshots int64
	lead              uint64

	// server
	wait            wait.Wait
	applyWait       wait.WaitTime
	leadTimeMtx     sync.RWMutex
	leadElectedTime time.Time

	// metrics
	slowReadIndex           prometheus.Counter
	readIndexFailed         prometheus.Counter
	leaderChanges           prometheus.Counter
	heartbeatSendFailures   prometheus.Counter
	proposalsFailed         prometheus.Counter
	proposalsPending        prometheus.Gauge
	hasLeader               prometheus.Gauge
	leadership              prometheus.Gauge
	proposalsCommitted      prometheus.Gauge
	proposalsApplied        prometheus.Gauge
	applySnapshotInProgress prometheus.Gauge
	learnerStatus           prometheus.Gauge

	// KV
	engine *pebble.DB
}

func New(cf *Config, logger *zap.Logger) (*Store, error) {
	var (
		w  *wal.WAL
		n  raft.Node
		cl membership.Cluster
		s  *raft.MemoryStorage
		id uint64
	)

	db, err := pebble.Open(cf.DataDir, defaultPebbleOptions())
	if err != nil {
		return nil, err
	}

	if cf.MaxRequestBytes > recommendedMaxRequestBytes {
		logger.Warn("exceeded recommended request bytes limit",
			zap.Uint64("recommended", recommendedMaxRequestBytes),
			zap.Uint64("configured", cf.MaxRequestBytes))
	}

	if err := fsutil.TouchDirAll(cf.DataDir); err != nil {
		return nil, errors.Wrapf(err, "cannot access data directory %q", cf.DataDir)
	}

	haveWAL := wal.Exist(cf.WALDir)

	if err := fsutil.TouchDirAll(cf.SnapDir); err != nil {
		return nil, errors.Wrapf(err, "cannot create snapshot directory %q", cf.SnapDir)
	}

	// todo: clean up temp files, eg: snapshots

	ss := snap.New(logger, cf.SnapDir)

	switch {
	case !haveWAL && !cf.NewCluster:
		// todo: verify join existing cluster config
		cl = membership.NewClusterFromAddresses(cf.InitialPeers)
		id, n, s, w = startNode(cf, logger, cl)

	case !haveWAL && cf.NewCluster:
		// new cluster
		cl := membership.NewClusterFromAddresses(cf.InitialPeers)

		id, n, s, w = startNode(cf, logger, cl)

	case haveWAL:
		if err := fsutil.IsDirWriteable(cf.MemberDir); err != nil {
			return nil, errors.Wrapf(err, "cannot write to member directory %q", cf.MemberDir)
		}

		if err := fsutil.IsDirWriteable(cf.WALDir); err != nil {
			return nil, errors.Wrapf(err, "cannot write to WAL directory %q", cf.WALDir)
		}

		// Find a snapshot to start or restart a raft node
		walSnaps, err := wal.ValidSnapshotEntries(logger, cf.WALDir)
		if err != nil {
			return nil, err
		}

		// snapshot files can be orphaned if manta crashes after writing them
		// but before writing the corresponding wal log entries
		snapshot, err := ss.LoadNewestAvailable(walSnaps)
		if err != nil && err != snap.ErrNoSnapshot {
			return nil, err
		}

		if snapshot != nil {
			// todo: recovery

		} else {
			logger.Info("no snapshot found. Recovering WAL from scratch")
		}

		if !cf.ForceNewCluster {
			id, cl, n, s, w = restartNode(cf, logger, snapshot)
		} else {
			id, cl, n, s, w = restartAsStandaloneNode(cf, logger, snapshot)
		}

	default:
		return nil, errors.New("unsupported bootstrap config")
	}

	heartbeat := time.Duration(cf.TickMs) * time.Millisecond
	store := &Store{
		id:                  id,
		config:              cf,
		logger:              logger,
		raftNode:            n,
		raftStorage:         s,
		cluster:             cl,
		wal:                 w,
		reqIDGen:            idutil.NewGenerator(uint16(id), time.Now()),
		applyCh:             make(chan apply),
		storage:             NewStorage(w, ss),
		td:                  contention.NewTimeoutDetector(2 * heartbeat),
		firstCommitInTermCh: make(chan struct{}),
		engine:              db,
		readStateC:          make(chan raft.ReadState, 1),

		slowReadIndex: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "slow_read_indexes_total",
		}),
		readIndexFailed: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "read_indexes_failed_total",
		}),
		leaderChanges: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "leader_changed_total",
		}),
		heartbeatSendFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "heartbeat_send_failed_total",
		}),
		proposalsFailed: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "proposal_failed_total",
		}),
		proposalsPending: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "proposal_pending",
		}),
		hasLeader: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "has_leader",
		}),
		leadership: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "is_leader",
		}),
		proposalsCommitted: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "proposals_committed_total",
		}),
		proposalsApplied: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "proposal_applied_total",
		}),
		applySnapshotInProgress: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "snapshot_apply_in_progress",
		}),
		learnerStatus: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "manta",
			Subsystem: "raft",
			Name:      "is_learner",
		}),
	}

	store.ticker = time.NewTicker(heartbeat)

	// TODO: add peers to transport
	store.transporter = transport.New(store, logger)
	for _, m := range cl.Members() {
		err := store.transporter.AddPeer(m.ID, m.Address)
		if err != nil {
			return nil, err
		}
	}

	return store, nil
}

func (s *Store) Run(ctx context.Context) error {
	s.start(ctx)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.adjustTicks(ctx)
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.linearizableReadLoop(ctx)
	}()

	s.wg.Wait()
	return nil
}

// start prepares and starts server in a new goroutine.
// It is no longer safe to modify a server's fields after
// it has been sent to Start. This function is just used for testing
func (s *Store) start(ctx context.Context) {
	if s.config.SnapshotCount == 0 {
		s.logger.Info("updating snapshot count to default",
			zap.Uint64("default", DefaultSnapshotCount))
		s.config.SnapshotCount = DefaultSnapshotCount
	}

	if s.config.SnapshotCatchUpEntries == 0 {
		s.logger.Info("updating snapshot catch-up entries to default",
			zap.Uint64("default", DefaultSnapshotCatchUpEntries))
		s.config.SnapshotCatchUpEntries = DefaultSnapshotCatchUpEntries
	}

	s.wait = wait.New()
	s.applyWait = wait.NewTimeList()
	s.readWaitCh = make(chan struct{}, 1)
	s.readNotifier = newNotifier()
	s.leaderChangedCh = make(chan struct{})

	s.logger.Info("staring raft store",
		zap.String("id", strconv.FormatUint(s.id, 16)))

	// TODO: if this is an empty log, writes all peer infos
	// into the first entry
	go s.run(ctx)
}

func (s *Store) run(ctx context.Context) {
	snapshot, err := s.raftStorage.Snapshot()
	if err != nil {
		s.logger.Fatal("failed to get snapshot from Raft Storage",
			zap.Error(err))
	}

	// asynchronously accept apply packets, dispatch progress in-order
	// TODO: fifo queue

	s.readyHandler = ReadyHandler{
		getLead: func() uint64 {
			return s.getLead()
		},
		updateLead: func(lead uint64) {
			s.setLead(lead)
		},
		updateLeadership: func(newLeader bool) {
			if newLeader {
				s.leaderChangedMtx.Lock()
				lc := s.leaderChangedCh
				s.leaderChangedCh = make(chan struct{})
				close(lc)
				s.leaderChangedMtx.Unlock()
			}
		},
		updateCommittedIndex: func(ci uint64) {
			curr := s.getCommittedIndex()
			if ci > curr {
				s.setCommittedIndex(ci)
			}
		},
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.raftLoop(ctx)
	}()

	sp := &progress{
		confState:     snapshot.Metadata.ConfState,
		snapshotIndex: snapshot.Metadata.Index,
		appliedTerm:   snapshot.Metadata.Term,
		appliedIndex:  snapshot.Metadata.Index,
	}

	for {
		select {
		case <-ctx.Done():
			return
		case ap := <-s.applyCh:
			s.applyAll(sp, &ap)
		case err := <-s.errCh:
			s.logger.Warn("server error, data-dir used by this member must be removed",
				zap.Error(err))
			return
		}
	}
}

type progress struct {
	confState     raftpb.ConfState
	snapshotIndex uint64
	appliedTerm   uint64
	appliedIndex  uint64
}

func (s *Store) applyAll(p *progress, apply *apply) {
	start := time.Now()
	defer func(before uint64) {
		s.logger.Info("apply all done",
			zap.Uint64("applied", p.appliedIndex-before),
			zap.Uint64("applied-index", p.appliedIndex),
			zap.Duration("elapsed", time.Since(start)))
	}(p.appliedIndex)

	s.applySnapshot(p, apply)
	s.applyEntries(p, apply)

	s.proposalsApplied.Set(float64(p.appliedIndex))
	s.applyWait.Trigger(p.appliedIndex)

	// wait for the raft routine to finish the disk writes before
	// triggering a snapshot. or applied index might be greater than
	// the last index in raft storage, since the raft routine might
	// be slower than apply routine.
	<-apply.notifyc

	s.triggerSnapshot(p)

	select {
	case m := <-s.msgSnapCh:
		merged := s.createMergedSnapshotMessage(m, p.appliedTerm, p.appliedIndex, p.confState)
		s.sendMergedSnap(merged)
	default:
	}
}

func (s *Store) raftLoop(ctx context.Context) {
	isLead := false
	internalTimeout := time.Second

	for {
		select {
		case <-ctx.Done():
			return

		case <-s.ticker.C:
			s.tickMtx.Lock()
			s.raftNode.Tick()
			s.tickMtx.Unlock()

		case rd := <-s.raftNode.Ready():
			if rd.SoftState != nil {
				newLeader := rd.SoftState.Lead != raft.None && s.readyHandler.getLead() != rd.SoftState.Lead
				if newLeader {
					s.leaderChanges.Inc()
				}

				if rd.SoftState.Lead == raft.None {
					s.hasLeader.Set(0)
				} else {
					s.hasLeader.Set(1)
				}

				s.readyHandler.updateLead(rd.SoftState.Lead)
				isLead = rd.RaftState == raft.StateLeader
				if isLead {
					s.leadership.Set(1)
				} else {
					s.leadership.Set(0)
				}

				s.readyHandler.updateLeadership(newLeader)
				s.td.Reset()
			}

			if len(rd.ReadStates) != 0 {
				select {
				case s.readStateC <- rd.ReadStates[len(rd.ReadStates)-1]:
				case <-time.After(internalTimeout):
					s.logger.Warn("timed out sending read state",
						zap.Duration("timeout", internalTimeout))
				case <-ctx.Done():
					return
				}
			}

			notifyCh := make(chan struct{}, 1)
			ap := apply{
				entries:  rd.CommittedEntries,
				snapshot: rd.Snapshot,
				notifyc:  notifyCh,
			}

			updateCommittedIndex(&ap, s.readyHandler)

			select {
			case s.applyCh <- ap:
			case <-ctx.Done():
				return
			}

			// the leader can write to its disk in parallel with replicating
			// to the followers and them writing to their disks. For more
			// detail, check raft thesis 10.2.1
			if isLead {
				s.transporter.Send(s.processMessages(rd.Messages))
			}

			// Must save the snapshot file and WAL snapshot entry before saving any
			// other entries or hardstate to ensure that reovery after a snapshot
			// restore is possible
			if !raft.IsEmptySnap(rd.Snapshot) {
				if err := s.storage.SaveSnap(rd.Snapshot); err != nil {
					s.logger.Fatal("save Raft snapshot failed",
						zap.Error(err))
				}
			}

			if err := s.storage.Save(rd.HardState, rd.Entries); err != nil {
				s.logger.Fatal("save Raft hard state and entries failed",
					zap.Error(err))
			}

			if !raft.IsEmptyHardState(rd.HardState) {
				s.proposalsCommitted.Set(float64(rd.HardState.Commit))
			}

			if !raft.IsEmptySnap(rd.Snapshot) {
				// Force WAL to fsync its hard state before Release() release
				// old data from the WAL. Otherwise could get an error like:
				// panic: tocommit(107) is out of range [lastIndex(84)]. Was
				// the raft log corrupted, truncated, or lost?
				// See https://github.com/etcd-io/etcd/issues/10219 for more details.
				if err := s.storage.Sync(); err != nil {
					s.logger.Fatal("sync raft snapshot failed",
						zap.Error(err))
				}

				// server now claim the snapshot has been persisted onto the disk
				notifyCh <- struct{}{}

				s.raftStorage.ApplySnapshot(rd.Snapshot)
				s.logger.Info("applied incoming Raft snapshot",
					zap.Uint64("term", rd.Snapshot.Metadata.Term),
					zap.Uint64("index", rd.Snapshot.Metadata.Index))

				if err := s.storage.Release(rd.Snapshot); err != nil {
					s.logger.Fatal("release Raft wal failed",
						zap.Error(err))
				}
			}

			s.raftStorage.Append(rd.Entries)

			if !isLead {
				// finish processing incoming messages before we signal raftdone chan
				msgs := s.processMessages(rd.Messages)

				// now unblocks 'applyAll' that waits on Raft log disk writes
				// before triggering snapshot
				notifyCh <- struct{}{}

				// Candidate or follower needs to wait for all pending configuration
				// changes to be applied before sending messages. Otherwise we might
				// incorrectly count votes (e.g. votes from removed members). Also
				// slow machine's follower raft-layer could proceed to become the
				// leader on its own single-node cluster, before apply-layer applies
				// the config change. We simply wait for ALL pending entries to be
				// applied for now. We might improve this later on if it causes unnecessary
				// long blocking issues.
				waitApply := false
				for _, ent := range rd.CommittedEntries {
					if ent.Type == raftpb.EntryConfChange {
						waitApply = true
						break
					}
				}

				if waitApply {
					// block until 'applyApp' calls 'applyWait.Trigger'
					// to be in sync with scheduled config-change job
					// (assume notifyCh has cap of 1)
					select {
					case <-ctx.Done():
						return
					case notifyCh <- struct{}{}:
					}
				}

				s.transporter.Send(msgs)
			} else {
				// leader already processed 'MsgSnap' and signaled
				notifyCh <- struct{}{}
			}

			s.raftNode.Advance()
		}
	}
}

func (s *Store) processMessages(messages []raftpb.Message) []raftpb.Message {
	sentAppResp := false

	for i := len(messages) - 1; i >= 0; i-- {
		if s.cluster.IDRemoved(messages[i].To) {
			messages[i].To = 0
		}

		if messages[i].Type == raftpb.MsgAppResp {
			if sentAppResp {
				messages[i].To = 0
			} else {
				sentAppResp = true
			}
		}

		if messages[i].Type == raftpb.MsgSnap {
			// The msgSnap only contains the most recent snapshot of store without KV.
			// So we need to redirect the msgSnap to server main loop for merging in
			// the current store snapshot and KV snapshot.
			select {
			case s.msgSnapCh <- messages[i]:
			default:
				// drop msgSnap if the inflight chan is full
			}

			messages[i].To = 0
		}

		if messages[i].Type == raftpb.MsgHeartbeat {
			ok, exceed := s.td.Observe(messages[i].To)
			if !ok {
				// TODO: limit request rate
				s.logger.Warn("leader failed to send out heartbeat on time;"+
					" took too long, leader is overloaded likely from slow disk",
					zap.String("to", fmt.Sprintf("%x", messages[i].To)),
					zap.Duration("heartbeat-interval", s.heartbeat),
					zap.Duration("expected-duration", 2*s.heartbeat),
					zap.Duration("exceeded-duration", exceed))

				s.heartbeatSendFailures.Inc()
			}
		}
	}

	return messages
}

func updateCommittedIndex(a *apply, rh ReadyHandler) {
	var ci uint64
	if len(a.entries) != 0 {
		ci = a.entries[len(a.entries)-1].Index
	}

	if a.snapshot.Metadata.Index > ci {
		ci = a.snapshot.Metadata.Index
	}

	if ci != 0 {
		rh.updateCommittedIndex(ci)
	}
}

func (s *Store) adjustTicks(ctx context.Context) {
	size := len(s.cluster.Members())

	// single-node fresh start, or single-node recovers from snapsho
	if size == 1 {
		ticks := s.config.ElectionTicks - 1
		s.logger.Info(
			"started as single-node; fast-forwarding election ticks",
			zap.String("member-id", strconv.FormatUint(s.id, 16)),
			zap.Int("forward-ticks", ticks),
			zap.String("forward-duration", tickToDur(ticks, s.config.TickMs)),
			zap.Int("election-ticks", s.config.ElectionTicks),
			zap.String("election-timeout", tickToDur(s.config.ElectionTicks, s.config.TickMs)),
		)

		s.advanceTicks(ticks)

		return
	}

	if !s.config.InitialElectionTickAdvance {
		s.logger.Info(
			"skipping initial election tick advance",
			zap.Int("election-ticks", s.config.ElectionTicks),
		)

		return
	}

	s.logger.Info("starting initial election tick advance",
		zap.Int("election-ticks", s.config.ElectionTicks))

	// retry up to "rafthttp.ConnReadTimeout", which is 5-sec
	// until peer connection reports; otherwise:
	// 1. all connection s fails
	// 2. no active peers
	// 3. restarted single-node with no snapshot
	// then, do nothing, because advancing ticks would have no effect
	waitTime := rafthttp.ConnReadTimeout
	itv := 50 * time.Millisecond
	for i := int64(0); i < int64(waitTime/itv); i++ {
		select {
		case <-time.After(itv):
		case <-ctx.Done():
			return
		}

		if peerN := s.ActivePeers(); peerN > 1 {
			// multi-node received peer connection reports
			// adjust ticks, in case slow leader message receive
			ticks := s.config.ElectionTicks - 2

			s.logger.Info(
				"initialized peer connections; fast-forwarding election ticks",
				zap.String("member-id", strconv.FormatUint(s.id, 16)),
				zap.Int("forward-ticks", ticks),
				zap.String("forward-duration", tickToDur(ticks, s.config.TickMs)),
				zap.Int("election-ticks", s.config.ElectionTicks),
				zap.String("election-timeout", tickToDur(s.config.ElectionTicks, s.config.TickMs)),
				zap.Int("active-remote-members", peerN),
			)

			s.advanceTicks(ticks)
			return
		}
	}
}

func (s *Store) ActivePeers() int {
	return s.transporter.ActivePeers()
}

// advanceTicks advances ticks of Raft node.
// This can be used for fast-forwarding election
// ticks in multi data-center deployments, thus
// speeding up election process.
func (s *Store) advanceTicks(ticks int) {
	for i := 0; i < ticks; i++ {
		s.tickMtx.Lock()
		s.raftNode.Tick()
		s.tickMtx.Unlock()
	}
}

func tickToDur(ticks int, tickMs uint) string {
	return fmt.Sprintf("%v", time.Duration(ticks)*time.Duration(tickMs)*time.Millisecond)
}

func (s *Store) linearizableReadLoop(ctx context.Context) {
	for {
		requestID := s.reqIDGen.Next()
		leaderChangeNotifier := s.leaderChangedNotify()

		select {
		case <-s.readWaitCh:
		case <-ctx.Done():
			return
		case <-leaderChangeNotifier:
			continue
		}

		nextNr := newNotifier()
		s.readMtx.Lock()
		nr := s.readNotifier
		s.readNotifier = nextNr
		s.readMtx.Unlock()

		confirmedIndex, err := s.requestCurrentIndex(ctx, leaderChangeNotifier, requestID)
		if err == raft.ErrStopped || err == ErrStopped || ctx.Err() != nil {
			return
		}

		if err != nil {
			nr.notify(err)
			continue
		}

		appliedIndex := s.getAppliedIndex()
		if appliedIndex < confirmedIndex {
			s.logger.Info("wait",
				zap.Uint64("applied", appliedIndex),
				zap.Uint64("confirmed", confirmedIndex))

			select {
			case <-s.applyWait.Wait(confirmedIndex):
			case <-ctx.Done():
				return
			}
		}

		// unblock all l-reads requested at indices before confirmedIndex
		nr.notify(nil)
	}
}

type notifier struct {
	c   chan struct{}
	err error
}

func newNotifier() *notifier {
	return &notifier{
		c: make(chan struct{}),
	}
}

func (nc *notifier) notify(err error) {
	nc.err = err
	close(nc.c)
}

func (s *Store) requestCurrentIndex(ctx context.Context, leaderChangeNotifier <-chan struct{}, reqID uint64) (uint64, error) {
	err := s.sendReadIndex(ctx, reqID)
	if err != nil {
		return 0, err
	}

	errTimer := time.NewTimer(s.config.RequestTimeout())
	defer errTimer.Stop()
	retryTimer := time.NewTimer(readIndexRetryTime)
	defer retryTimer.Stop()

	firstCommitInTermNotifier := s.FirstCommitInTermNotify()

	for {
		select {
		case rs := <-s.readStateC:
			requestIdBytes := uint64ToBigEndianBytes(reqID)
			gotOwnResp := bytes.Equal(rs.RequestCtx, requestIdBytes)
			if !gotOwnResp {
				// a previous request might time out. now we should
				// ignore the resp of it and continue waiting for
				// the response of the current requests
				respID := uint64(0)
				if len(rs.RequestCtx) == 8 {
					respID = binary.BigEndian.Uint64(rs.RequestCtx)
				}

				s.logger.Warn("ignore out-of-date read index response; local node read indexes queueing up and waiting to be in sync with leader",
					zap.Uint64("req", reqID),
					zap.Uint64("resp", respID))

				s.slowReadIndex.Inc()
				continue
			}

			return rs.Index, nil

		case <-leaderChangeNotifier:
			s.readIndexFailed.Inc()
			// return a retryable error
			return 0, ErrLeaderChanged

		case <-firstCommitInTermNotifier:
			firstCommitInTermNotifier = s.FirstCommitInTermNotify()
			s.logger.Info("first commit in current term: resending ReadIndex request")
			err := s.sendReadIndex(ctx, reqID)
			if err != nil {
				return 0, err
			}

			retryTimer.Reset(readIndexRetryTime)
			continue

		case <-retryTimer.C:
			s.logger.Warn(
				"waiting for ReadIndex response took too long, retrying",
				zap.Uint64("req", reqID),
				zap.Duration("retry-timeout", readIndexRetryTime),
			)

			err := s.sendReadIndex(ctx, reqID)
			if err != nil {
				return 0, err
			}

			retryTimer.Reset(readIndexRetryTime)
			continue

		case <-errTimer.C:
			s.logger.Warn(
				"timed out waiting for read index response, local node might have slow network",
				zap.Duration("timeout", s.config.RequestTimeout()),
			)

			s.slowReadIndex.Inc()
			return 0, ErrTimeout

		case <-ctx.Done():
			return 0, ErrStopped
		}
	}
}

func uint64ToBigEndianBytes(number uint64) []byte {
	byteResult := make([]byte, 8)
	binary.BigEndian.PutUint64(byteResult, number)
	return byteResult
}

func (s *Store) sendReadIndex(ctx context.Context, reqIndex uint64) error {
	data := uint64ToBigEndianBytes(reqIndex)

	ctx, cancel := context.WithTimeout(ctx, s.config.RequestTimeout())
	err := s.raftNode.ReadIndex(ctx, data)
	cancel()

	if err == raft.ErrStopped {
		return err
	}

	if err != nil {
		s.logger.Warn("failed to get read index from Raft",
			zap.Error(err))
		s.readIndexFailed.Inc()
		return err
	}

	return nil
}

// FirstCommitInTermNotify returns channel that will be unlocked
// on first entry committed in new term, which is necessary for new
// leader to answer read-only request (leadser is not able to respond
// any read-only requests as long as linearizable semantic is required)
func (s *Store) FirstCommitInTermNotify() <-chan struct{} {
	s.firstCommitInTermMtx.RLock()
	ch := s.firstCommitInTermCh
	s.firstCommitInTermMtx.RUnlock()

	return ch
}

func (s *Store) LeaderChangedNotify() <-chan struct{} {
	s.leaderChangedMtx.RLock()
	ch := s.leaderChangedCh
	s.leaderChangedMtx.RUnlock()

	return ch
}

func (s *Store) leaderChangedNotify() <-chan struct{} {
	s.leaderChangedMtx.RLock()
	ch := s.leaderChangedCh
	s.leaderChangedMtx.RUnlock()

	return ch
}

func (s *Store) getAppliedIndex() uint64 {
	return atomic.LoadUint64(&s.appliedIndex)
}

func (s *Store) getCommittedIndex() uint64 {
	return atomic.LoadUint64(&s.committedIndex)
}

func (s *Store) setCommittedIndex(ci uint64) {
	atomic.StoreUint64(&s.committedIndex, ci)
}

func (s *Store) triggerSnapshot(p *progress) {
	if p.appliedIndex-p.snapshotIndex <= s.config.SnapshotCount {
		return
	}

	s.logger.Info("triggering snapshot",
		zap.Uint64("applied", p.appliedIndex),
		zap.Uint64("snapshot", p.snapshotIndex),
		zap.Uint64("snapshot-count", s.config.SnapshotCount))

	s.snapshot(p.appliedIndex, p.confState)
	p.snapshotIndex = p.appliedIndex
}

// TODO: non-blocking snapshot
func (s *Store) snapshot(index uint64, confState raftpb.ConfState) {
	s.wg.Add(1)
	defer s.wg.Done()

	// TODO: take snapshot of the KV engine
	data := []byte{}
	snapshot, err := s.raftStorage.CreateSnapshot(index, &confState, data)
	if err != nil {
		// the snapshot was done asynchronously with the progress of raft.
		// raft might have already got a newer snapshot.
		if err == raft.ErrSnapOutOfDate {
			s.logger.Fatal("create snapshot failed",
				zap.Error(err))
		}
		return
	}

	// SaveSnap saves the snapshot to file and appends the corresponding
	// WAL entry.
	if err = s.storage.SaveSnap(snapshot); err != nil {
		s.logger.Fatal("save snapshot failed",
			zap.Error(err))
	}

	if err = s.storage.Release(snapshot); err != nil {
		s.logger.Fatal("release snapshot failed",
			zap.Error(err))
	}

	s.logger.Info("snapshot saved",
		zap.Uint64("snapshot-index", snapshot.Metadata.Index))

	// When sending a snapshot, server will pause compaction. After
	// receives a snapshot, the slow follower needs to get all the entries
	// right after snapshot sent to catch up. If we do not pause
	// compaction, the log entries right after the snapshot sent might already be
	// compacted. It happens when the snapshot takes long time to send
	// and save. Pausing compaction avoids triggering a snapshot sending cycle.
	if atomic.LoadInt64(&s.inflightSnapshots) != 0 {
		s.logger.Info("skip compaction since there is an inflight snapshot")
		return
	}

	// keep some in memory log entries for slow followers
	compactIndex := uint64(1)
	if index > s.config.SnapshotCatchUpEntries {
		compactIndex = index - s.config.SnapshotCatchUpEntries
	}

	err = s.raftStorage.Compact(compactIndex)
	if err != nil {
		// the compaction was done asynchronously with the progress of raft.
		// raft log might already been compact.
		if err == raft.ErrCompacted {
			s.logger.Fatal("compact failed",
				zap.Error(err))
		}
	}

	s.logger.Info("raft logs compacted",
		zap.Uint64("compact-index", compactIndex))
}

func (s *Store) applySnapshot(p *progress, apply *apply) {
	if raft.IsEmptySnap(apply.snapshot) {
		return
	}

	s.applySnapshotInProgress.Inc()

	s.logger.Info("applying snapshot",
		zap.Uint64("current-snapshot-index", p.snapshotIndex),
		zap.Uint64("current-applied-index", p.appliedIndex),
		zap.Uint64("incoming-leader-snapshot-index", apply.snapshot.Metadata.Index),
		zap.Uint64("incoming-leader-snapshot-term", apply.snapshot.Metadata.Term))

	defer func() {
		s.logger.Info("snapshot applied successfully",
			zap.Uint64("current-snapshot-index", p.snapshotIndex),
			zap.Uint64("current-spplied-index", p.appliedIndex),
			zap.Uint64("incoming-leader-snapshot-index", apply.snapshot.Metadata.Index),
			zap.Uint64("incoming-leader-snapshot-term", apply.snapshot.Metadata.Term))

		s.applySnapshotInProgress.Dec()
	}()

	if apply.snapshot.Metadata.Index <= p.appliedIndex {
		s.logger.Fatal("unexpected leader snapshot from outdated index",
			zap.Uint64("current-snapshot-index", p.snapshotIndex),
			zap.Uint64("current-applied-index", p.appliedIndex),
			zap.Uint64("incoming-leader-snapshot-index", apply.snapshot.Metadata.Index),
			zap.Uint64("incoming-leader-snapshot-term", apply.snapshot.Metadata.Term))
	}

	// wait for raftNode to persist snapshot onto the disk
	<-apply.notifyc

	// TODO: create new KV engine for the snapshot

	s.logger.Info("restoring kv store")

	s.transporter.RemoveAllPeers()

	s.logger.Info("remove old peers from network")
	s.logger.Info("adding peers from new cluster configuration")

	s.logger.Info("added peers from new cluster configuration")

	p.appliedTerm = apply.snapshot.Metadata.Term
	p.appliedIndex = apply.snapshot.Metadata.Index
	p.snapshotIndex = p.appliedIndex
	p.confState = apply.snapshot.Metadata.ConfState
}

func (s *Store) applyEntries(p *progress, apply *apply) {
	if len(apply.entries) == 0 {
		return
	}

	firstIndex := apply.entries[0].Index
	if firstIndex > p.appliedIndex+1 {
		s.logger.Fatal("unexpected committed entry index",
			zap.Uint64("current-applied-index", p.appliedIndex),
			zap.Uint64("first-committed-entry-index", firstIndex))
	}

	var entries []raftpb.Entry
	if p.appliedIndex+1-firstIndex < uint64(len(apply.entries)) {
		entries = apply.entries[p.appliedIndex+1-firstIndex:]
	}

	if len(entries) == 0 {
		return
	}

	var shouldStop bool
	if p.appliedTerm, p.appliedIndex, shouldStop = s.applyBatch(entries, &p.confState); shouldStop {
		go s.stopWithDelay(10*100*time.Millisecond,
			fmt.Errorf("the member has been permanently removed from the cluster"))
	}
}

func (s *Store) stopWithDelay(duration time.Duration, err error) {
	// TODO: pass the root context
	ctx := context.Background()

	select {
	case <-time.After(duration):
	case <-ctx.Done():
		return
	}

	select {
	case s.errCh <- err:
	default:
	}
}

type confChangeResponse struct {
	membs []*membership.Member
	err   error
}

// apply takes entries received from Raft (after it has been committed) and
// applies them to the current state of the Store. The given entries should
// not be empty
func (s *Store) apply(entries []raftpb.Entry, confState *raftpb.ConfState) (uint64, uint64, bool) {
	var (
		appliedIndex, appliedTerm uint64
		shouldStop                bool
	)

	s.logger.Debug("applying entries",
		zap.Int("count", len(entries)))

	for i := range entries {
		entry := entries[i]
		s.logger.Debug("applying entry",
			zap.Uint64("index", entry.Index),
			zap.Uint64("term", entry.Term),
			zap.Stringer("type", entry.Type))

		switch entry.Type {
		case raftpb.EntryNormal:
			s.applyEntryNormal(&entry)
			s.setAppliedIndex(entry.Index)
			s.setTerm(entry.Term)

		case raftpb.EntryConfChange:
			var cc raftpb.ConfChange
			pbutil.MustUnmarshal(&cc, entry.Data)
			removedSelf, err := s.applyConfChange(cc, confState)
			s.setAppliedIndex(entry.Index)
			s.setTerm(entry.Term)

			shouldStop = shouldStop || removedSelf
			s.wait.Trigger(cc.ID, &confChangeResponse{s.cluster.Members(), err})

		default:
			s.logger.Fatal("unknown entry type, it must be EntryNormal or EntryConfChange",
				zap.Stringer("type", entry.Type))
		}

		appliedIndex, appliedTerm = entry.Index, entry.Term
	}

	return appliedTerm, appliedIndex, shouldStop
}

func (s *Store) applyBatch(entries []raftpb.Entry, confState *raftpb.ConfState) (uint64, uint64, bool) {
	var (
		appliedIndex, appliedTerm uint64
		shouldStop                bool
		start, i                  int
	)

	for i = start; i < len(entries); i++ {
		entry := entries[i]
		if entry.Type == raftpb.EntryNormal {
			continue
		}

		if entry.Type == raftpb.EntryConfChange {
			// apply pending EntryConfChange
			s.applyNormalEntries(entries[start:i])
			s.setAppliedIndex(entry.Index)
			s.setTerm(entry.Index)

			appliedIndex, appliedTerm = entry.Index, entry.Term

			// apply EntryConfChange
			var cc raftpb.ConfChange
			pbutil.MustUnmarshal(&cc, entry.Data)
			removedSelf, err := s.applyConfChange(cc, confState)
			s.setAppliedIndex(entry.Index)
			s.setTerm(entry.Term)

			shouldStop = shouldStop || removedSelf
			s.wait.Trigger(cc.ID, &confChangeResponse{s.cluster.Members(), err})

			appliedIndex, appliedTerm = entry.Index, entry.Term

			start = i

			continue
		}

		s.logger.Fatal("unknown entry type, it must be EntryNormal or EntryConfChange",
			zap.Stringer("type", entries[i].Type))
	}

	// apply remaining entries
	if start != len(entries) {
		entry := entries[len(entries)-1]

		s.applyNormalEntries(entries[start:i])
		s.setAppliedIndex(entry.Index)
		s.setTerm(entry.Index)

		appliedIndex, appliedTerm = entry.Index, entry.Term
	}

	return appliedTerm, appliedIndex, shouldStop
}

func (s *Store) applyNormalEntries(entries []raftpb.Entry) {
	wo := &pebble.WriteOptions{Sync: true}
	batch := s.engine.NewBatch()
	defer batch.Commit(wo)

	for _, entry := range entries {
		req := &rawkv.PutRequest{}
		err := req.Unmarshal(entry.Data)
		if err != nil {
			s.wait.Trigger(req.Id, err)
			continue
		}

		err = batch.Set(req.Key, req.Value, wo)
		if err != nil {
			panic(err)
		} else {
			s.wait.Trigger(req.Id, nil)
		}
	}
}

// applyEntryNormal applies an EntryNormal type raftpb request to the Store
func (s *Store) applyEntryNormal(entry *raftpb.Entry) {
	s.logger.Debug("apply entry normal",
		zap.Uint64("index", entry.Index),
		zap.Uint64("term", entry.Term))

	// raft state machine may generate noop entry when leader confirmation
	// skip it in advance to avoid some potential bug in the future
	if len(entry.Data) == 0 {
		s.notifyAboutFirstCommitInTerm()
		return
	}

	// TODO: apply to KV engine
	req := &rawkv.PutRequest{}
	err := req.Unmarshal(entry.Data)
	if err != nil {
		s.logger.Warn("decode req failed",
			zap.Error(err))
		return
	}

	err = s.engine.Set(req.Key, req.Value, nil)
	if err != nil {
		s.logger.Warn("write kv failed",
			zap.Error(err))
	}

	s.wait.Trigger(req.Id, nil)
}

func (s *Store) notifyAboutFirstCommitInTerm() {
	notifyCh := make(chan struct{})
	s.firstCommitInTermMtx.Lock()
	notifierToClose := s.firstCommitInTermCh
	s.firstCommitInTermCh = notifyCh
	s.firstCommitInTermMtx.Unlock()
	close(notifierToClose)
}

// applyConfChange applies a ConfChange to the server. It is only
// invoked with a ConfChange that has already passed through Raft
func (s *Store) applyConfChange(cc raftpb.ConfChange, confState *raftpb.ConfState) (bool, error) {
	if err := s.cluster.ValidateConfigurationChange(cc); err != nil {
		cc.NodeID = raft.None
		s.raftNode.ApplyConfChange(cc)
		return false, err
	}

	*confState = *s.raftNode.ApplyConfChange(cc)
	switch cc.Type {
	case raftpb.ConfChangeAddNode, raftpb.ConfChangeAddLearnerNode:
		confChangeContext := &membership.ConfigChangeContext{}
		if err := json.Unmarshal(cc.Context, confChangeContext); err != nil {
			s.logger.Fatal("failed to unmarshal member",
				zap.Error(err))
		}

		if cc.NodeID != uint64(confChangeContext.Member.ID) {
			s.logger.Fatal("got different member ID",
				zap.String("member-id-from-config-change-entry", strconv.FormatUint(cc.NodeID, 16)),
				zap.String("member-id-from-message", strconv.FormatUint(confChangeContext.Member.ID, 16)))
		}

		if confChangeContext.IsPromote {
			s.cluster.PromoteMember(confChangeContext.Member.ID)
		} else {
			s.cluster.AddMember(&confChangeContext.Member)

			if confChangeContext.Member.ID != s.id {
				s.transporter.AddPeer(confChangeContext.Member.ID, confChangeContext.Member.Address)
			}
		}

		// update the isLearner metric when this server id is equal to the id in raft member confChange
		if confChangeContext.Member.ID == s.id {
			if cc.Type == raftpb.ConfChangeAddLearnerNode {
				s.learnerStatus.Set(1)
			} else {
				s.learnerStatus.Set(0)
			}
		}

	case raftpb.ConfChangeRemoveNode:
		id := cc.NodeID
		s.cluster.RemoveMember(id)
		if id == s.id {
			return true, nil
		}

		s.transporter.RemovePeer(id)

	case raftpb.ConfChangeUpdateNode:
		m := &membership.Member{}
		if err := json.Unmarshal(cc.Context, m); err != nil {
			panic("failed to unmarshal member, err: " + err.Error())
		}

		if cc.NodeID != m.ID {
			s.logger.Panic("got different member ID",
				zap.String("member-id-from-config-change-entry", strconv.FormatUint(cc.NodeID, 16)),
				zap.String("member-id-from-message", strconv.FormatUint(m.ID, 16)))
		}

		s.cluster.UpdateMember(m.ID, m.Address)
		if m.ID != s.id {
			s.transporter.UpdatePeer(m.ID, m.Address)
		}
	}

	return false, nil
}

func (s *Store) setAppliedIndex(index uint64) {
	atomic.StoreUint64(&s.appliedIndex, index)
}

func (s *Store) setTerm(term uint64) {
	atomic.StoreUint64(&s.term, term)
}

// Process implement transport.Raft, it takes a raft message and applies it
// to the server's raft state machine, respecting any timeout of given context.
func (s *Store) Process(ctx context.Context, m raftpb.Message) error {
	if s.cluster.IDRemoved(m.From) {
		s.logger.Warn("rejected Raft message from removed member",
			zap.String("local", internal.IDToString(s.id)),
			zap.String("from", internal.IDToString(m.From)))

		return ErrRejectFromRemovedMember
	}

	return s.raftNode.Step(ctx, m)
}

// IsIDRemoved implement transport.Raft
func (s *Store) IsIDRemoved(id uint64) bool {
	return s.cluster.IDRemoved(id)
}

// ReportUnreachable implement transport.Raft
func (s *Store) ReportUnreachable(id uint64) {
	s.raftNode.ReportUnreachable(id)
}

// ReportSnapshot implement transport.Raft, it reports snapshot sent status
// to the raft state machin, and clears the used snapshot from the snapshot
// store.
func (s *Store) ReportSnapshot(id uint64, status raft.SnapshotStatus) {
	s.raftNode.ReportSnapshot(id, status)
}
