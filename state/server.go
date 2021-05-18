package state

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/etcdserver/stats"
	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/snowflake"
	"github.com/f1shl3gs/manta/state/cindex"
	"github.com/f1shl3gs/manta/state/membership"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/pkg/v3/fileutil"
	"go.etcd.io/etcd/pkg/v3/pbutil"
	"go.etcd.io/etcd/pkg/v3/schedule"
	"go.etcd.io/etcd/pkg/v3/types"
	"go.etcd.io/etcd/pkg/v3/wait"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.etcd.io/etcd/server/v3/wal"
	"go.uber.org/zap"
)

const (
	DefaultSnapshotCount = 100000

	// DefaultSnapshotCatchUpEntries is the number of entries for a slow
	// follower to catch-up after compacting the raft storage entries.
	// We expect the follower has a millisecond level latency with the
	// leader. The max throughput is around 10K. Keep a 5K entries is enough
	// for helping follower to catch up
	DefaultSnapshotCatchUpEntries uint64 = 5000

	// HealthInterval is the minimum time the cluster should be healthy
	// before accepting add member requests.
	HealthInterval = 5 * time.Second

	purgeFileInterval = 30 * time.Second

	// Max number of in-flight snapshot messages the server allows
	// to have. This number is more than enough for most clusters with 5
	// machines.
	maxInFlightMsgSnap        = 16
	releaseDelayAfterSnapshot = 30 * time.Second

	// maxPendingRevokes is the maximum number of outstanding expired
	// lease revocations
	maxPendingRevokes = 16

	recommendedMaxRequestBytes = 10 * 1024 * 1024

	defaultMaxSizePerMsg   = 1024 * 1024
	defaultMaxInflightMsgs = 256
)

type Config struct {
	Listen  string
	DataDir string
	Peers   []string

	SnapshotCount          uint64
	SnapshotCatchUpEntries uint64
}

func (cf *Config) WALDir() string {
	return filepath.Join(cf.DataDir, "wal")
}

func (cf *Config) SnapDir() string {
	return filepath.Join(cf.DataDir, "snap")
}

type Server struct {
	Config

	ID                manta.ID
	logger            *zap.Logger
	node              raft.Node
	errCh             chan error
	wg                sync.WaitGroup
	inflightSnapshots int64
	cluster           membership.Cluster

	idGen *snowflake.Generator

	// raft
	raftNode        *raftNode
	raftStorage     raft.Storage
	consistentIndex cindex.ConsistentIndexer

	// example implement
	transport    *rafthttp.Transport
	wait         wait.Wait
	applyWait    wait.WaitTime
	doneCh       chan struct{}
	stopCh       chan struct{}
	stopping     chan struct{}
	cancel       context.CancelFunc
	ctx          context.Context
	readwaitC    chan struct{}
	readNotifier *notifier

	// leaderChanged is used to notify the linearizable read loop to drop the old read requests.
	leaderChanged   chan struct{}
	leaderChangedMu sync.RWMutex
	lead            uint64

	// stats
	stats       *stats.ServerStats
	leaderStats *stats.LeaderStats

	committedIndex uint64

	// metrics
	applySnapshotInProgress prometheus.Gauge
	proposalsApplied        prometheus.Gauge
	isLearner               prometheus.Gauge

	// unknown field
	forceVersionC chan struct{}

	leadTimeMu      sync.RWMutex
	leadElectedTime time.Time
	appliedIndex    uint64
	term            uint64
}

// Process implement rafthttp.Raft
func (s *Server) Process(ctx context.Context, m raftpb.Message) error {
	return s.node.Step(ctx, m)
}

// IsIDRemoved implement rafthttp.Raft
// todo: find out the etcd implement, and figure out what it is
func (s *Server) IsIDRemoved(id uint64) bool {
	return false
}

// ReportUnreachable implement rafthttp.Raft
func (s *Server) ReportUnreachable(id uint64) {
	s.node.ReportUnreachable(id)
}

// ReportSnapshot implement rafthttp.Raft
func (s *Server) ReportSnapshot(id uint64, status raft.SnapshotStatus) {
	s.node.ReportSnapshot(id, status)
}

func New(cf *Config, logger *zap.Logger) (*Server, error) {
	var (
		w    *wal.WAL
		node raft.Node
		ms   *raft.MemoryStorage
	)

	svr := &Server{
		idGen: snowflake.New(os.Getpid() % 1024),
	}

	if err := fileutil.TouchDirAll(cf.DataDir); err != nil {
		return nil, errors.Errorf("cannot access data dir %v, err: %s",
			cf.DataDir, err)
	}

	haveWAL := wal.Exist(cf.WALDir())

	if err := fileutil.TouchDirAll(cf.SnapDir()); err != nil {
		return nil, errors.Wrapf(err, "failed to create snapshot directory")
	}

	// todo: remove tmp files under snap dir

	ss := snap.New(logger, cf.SnapDir())

	rc := &raft.Config{
		ID:                        svr.idGen.Next(),
		ElectionTick:              10,
		HeartbeatTick:             1,
		Storage:                   ms,
		MaxSizePerMsg:             defaultMaxSizePerMsg,
		MaxInflightMsgs:           defaultMaxInflightMsgs,
		MaxUncommittedEntriesSize: 1 << 30,
	}

	if haveWAL {
		walSnap, err := wal.ValidSnapshotEntries(logger, cf.WALDir())
		if err != nil {
			return nil, err
		}

		// snapshot files can be orphaned if etcd crashes
		// after writing them but before writing the corresponding
		// wal log entries
		snapshot, err := ss.LoadNewestAvailable(walSnap)
		if err != nil && err != snap.ErrNoSnapshot {
			return nil, err
		}

		if snapshot != nil {

		}

		node = raft.RestartNode(rc)
	} else {
		rps := make([]raft.Peer, len(cf.Peers))
		for i := range rps {
			rps[i] = raft.Peer{ID: generateRaftID(cf.Peers[i])}
		}
		node = raft.StartNode(rc, rps)
	}

	// find a snapshot to start/restart a raft node
	walSnaps, err := wal.ValidSnapshotEntries(logger, cf.WALDir())
	if err != nil {
		return nil, err
	}

	// snapshot files can be orphaned if etcd crashes after writing
	// them but before writing the corresponding wal log entries
	snapshot, err := ss.LoadNewestAvailable(walSnaps)
	if err != nil {
		return nil, err
	}

	st := &Store{}
	if snapshot != nil {
		if err := st.Recovery(snapshot.Data); err != nil {
			return nil, err
		}

		logger.Info("Recover store from snapshot",
			zap.Uint64("index", snapshot.Metadata.Index),
			zap.Int("size", snapshot.Size()))
	}

	svr.transport = &rafthttp.Transport{
		Logger:    logger,
		ID:        types.ID(generateRaftID(cf.Listen)),
		ClusterID: 0x1000,
		Raft:      svr,
	}

	svr.transport.Start()

	for i := range cf.Peers {
		if cf.Peers[i] == cf.Listen {
			continue
		}

		svr.transport.AddPeer(types.ID(generateRaftID(cf.Peers[i])), []string{cf.Peers[i]})
	}

	return svr, nil
}

func generateRaftID(addr string) uint64 {
	h := fnv.New64()
	h.Sum([]byte(addr))
	return h.Sum64()
}

// Start performs any initialization of the Server necessary for it to
// begin serving requests. It must be called before Do or Process.
// Start must be non-blocking; any long-running server functionality
// should be implemented in goroutines
func (s *Server) Start() {

}

// start prepares and starts server in a new goroutine. It is no longer
// safe to modify a server's fields after it has been sent to Start.
// This function is just used for testing.
func (s *Server) start() {
	if s.Config.SnapshotCount == 0 {
		s.logger.Info("updating snapshot count to default",
			zap.Uint64("count", DefaultSnapshotCount))
		s.Config.SnapshotCount = DefaultSnapshotCount
	}

	if s.Config.SnapshotCatchUpEntries == 0 {
		s.logger.Info("updating snapshot catch-up entries to default",
			zap.Uint64("count", DefaultSnapshotCatchUpEntries))
		s.Config.SnapshotCatchUpEntries = DefaultSnapshotCatchUpEntries
	}

	s.wait = wait.New()
	s.applyWait = wait.NewTimeList()
	s.doneCh = make(chan struct{})
	s.stopCh = make(chan struct{})
	s.stopping = make(chan struct{}, 1)
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.readwaitC = make(chan struct{}, 1)
	s.readNotifier = newNotifier()
	s.leaderChanged = make(chan struct{})

	/*
		if s.ClusterVersion() != nil {

		}
	*/

	// TODO: if this is an empty log, writes all peer infos
	// inot the first entry
	go s.run()
}

func (s *Server) run(ctx context.Context) {
	sn, err := s.raftStorage.Snapshot()
	if err != nil {
		s.logger.Panic("failed to get snapshot from raft storage",
			zap.Error(err))
	}

	// asynchronously accept apply packets, dispatch progress in-order
	sched := schedule.NewFIFOScheduler()

	rh := &raftReadyHandler{
		getLead: func() (lead uint64) {
			return s.getLead()
		},
		updateLead: func(lead uint64) {
			s.setLead(lead)
		},
		updateLeadership: func(newLeader bool) {
			if !s.isLeader() {

				// todo: handle compactor properly
			} else {
				if newLeader {
					t := time.Now()
					s.leadTimeMu.Lock()
					s.leadElectedTime = t
					s.leadTimeMu.Unlock()
				}

				// todo: handle compactor
			}

			if newLeader {
				s.leaderChangedMu.Lock()
				lc := s.leaderChanged
				s.leaderChanged = make(chan struct{})
				close(lc)
				s.leaderChangedMu.Unlock()
			}

			// TODO: remove the nil checking
			// current test utility does not provide the stats
			if s.stats != nil {
				s.stats.BecomeLeader()
			}
		},
		updateCommittedIndex: func(ci uint64) {
			cci := s.getCommittedIndex()
			if ci > cci {
				s.setCommittedIndex(ci)
			}
		},
	}

	s.raftNode.start(rh)

	ep := etcdProgress{
		confState: sn.Metadata.ConfState,
		snapi:     sn.Metadata.Index,
		appliedt:  sn.Metadata.Term,
		appliedi:  sn.Metadata.Index,
	}

	defer func() {
		close(s.stopping)

	}()

	for {
		select {
		case ap := <-s.raftNode.apply():
			f := func(ctx context.Context) {
				s.applyAll(&ep, &ap)
			}
			sched.Schedule(f)

		case err := <-s.errCh:
			s.logger.Warn("server error",
				zap.Error(err))
			return

		case <-s.stopCh:
			return
		}
	}
}

type etcdProgress struct {
	confState raftpb.ConfState
	snapi     uint64
	appliedt  uint64
	appliedi  uint64
}

func (s *Server) applyAll(ep *etcdProgress, apply *apply) {
	s.applySnapshot(ep, apply)
	s.applyEntries(ep, apply)

	s.proposalsApplied.Set(float64(ep.appliedi))
	s.applyWait.Trigger(ep.appliedi)

	// wait for the raft routine to finish the disk writes before
	// triggering a snapshot, or applied index might be greater than
	// the last index in raft storage, since the raft routine might
	// be slower than apply routine.
	<-apply.notifyc

	s.triggerSnapshot(ep)

	select {
	case m := <-s.raftNode.msgSnapC:
		merged := s.createMergedSnapshotMessage(m, ep.appliedt, ep.appliedi, ep.confState)
		s.sendMergedSnap(merged)
	default:
	}
}

func (s *Server) sendMergedSnap(merged snap.Message) {
	atomic.AddInt64(&s.inflightSnapshots, 1)

	fields := []zap.Field{
		zap.String("from", s.ID.String()),
		zap.String("to", types.ID(merged.To).String()),
		zap.Int64("bytes", merged.TotalSize),
	}

	now := time.Now()
	s.raftNode.transport.SendSnapshot(merged)

	s.logger.Info("sending merged snapshot", fields...)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		select {
		case ok := <-merged.CloseNotify():
			// delay releasing inflight snapshot for another 30 seconds
			// to block log compaction. If the follower still fails to catch
			// up, it is probably just too slow to catch up. We cannot avoid
			// the snapshot cycle any way
			if ok {
				select {
				case <-time.After(releaseDelayAfterSnapshot):
				case <-s.stopping:
				}
			}

			atomic.AddInt64(&s.inflightSnapshots, -1)
			s.logger.Info("sent merged snapshot",
				append(fields, zap.Duration("took", time.Since(now)))...)

		case <-s.stopping:
			s.logger.Warn("canceled sending merged snapshot; server stopping",
				fields...)
		}
	}()
}

func (s *Server) triggerSnapshot(ep *etcdProgress) {
	if ep.appliedi-ep.snapi <= s.Config.SnapshotCount {
		return
	}

	s.logger.Info("triggering snapshot",
		zap.Uint64("applied", ep.appliedi),
		zap.Uint64("snapshot", ep.snapi),
		zap.Uint64("snapshot-count", s.Config.SnapshotCount))

	s.snapshot(ep.appliedi, ep.confState)
	ep.snapi = ep.appliedi
}

func (s *Server) snapshot(snapIndex uint64, confState raftpb.ConfState) {
	// flush !?

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// data
		snap, err := s.raftNode.raftStorage.CreateSnapshot(snapIndex, &confState, []byte("a"))
		if err != nil {
			// the snapshot was done asynchronously with the progress of raft
			// raft might have already got a newer snapshot
			if err == raft.ErrSnapOutOfDate {
				return
			}

			s.logger.Panic("failed to create snapshot",
				zap.Error(err))
		}

		// SaveSnap saves the snapshot to file and appends the corresponding WAL entry
		if err = s.raftNode.storage.SaveSnap(snap); err != nil {
			s.logger.Panic("failed to save snapshot",
				zap.Error(err))
		}

		if err = s.raftNode.storage.Release(snap); err != nil {
			s.logger.Panic("failed to release wal",
				zap.Error(err))
		}

		s.logger.Info("save snapshot success",
			zap.Uint64("index", snap.Metadata.Index))

		// When sending a snapshot, etcd will pause compaction
		// After receives a snapshot, the slow follower needs to get all the
		// entries right after the snapshot sent to catch up. If we
		// do not pause compaction, the log entries right after the snapshot
		// send might already be compacted. It happens when the snapshot
		// takes long time to send and save. Pausing compaction avoids
		// triggering a snapshot sending cyley
		if atomic.LoadInt64(&s.inflightSnapshots) != 0 {
			s.logger.Info("skip compaction since there is an inflight snapshot")
			return
		}

		// keep some in memory log entries for slow followers
		compactIndex := uint64(1)
		if snapIndex > s.Config.SnapshotCatchUpEntries {
			compactIndex = snapIndex - s.Config.SnapshotCatchUpEntries
		}

		err = s.raftNode.raftStorage.Compact(compactIndex)
		if err != nil {
			// the compaction was done asynchronously with the progress of
			// raft. Raft log might already been compact
			if err == raft.ErrCompacted {
				return
			}

			s.logger.Panic("failed to compact",
				zap.Error(err))
		}

		s.logger.Info("compacted raft logs",
			zap.Uint64("comapct-index", compactIndex))
	}()
}

func (s *Server) applySnapshot(ep *etcdProgress, apply *apply) {
	if raft.IsEmptySnap(apply.snapshot) {
		return
	}

	s.applySnapshotInProgress.Inc()

	s.logger.Info("applying snapshot",
		zap.Uint64("snapshot-index", ep.snapi),
		zap.Uint64("applied-index", ep.appliedi),
		zap.Uint64("incoming-leader-snapshot-index", apply.snapshot.Metadata.Index),
		zap.Uint64("incoming-leader-snapshot-term", apply.snapshot.Metadata.Term))

	defer func(start time.Time) {
		s.logger.Info("applied snapshot",
			zap.Uint64("snapshot-index", ep.snapi),
			zap.Uint64("applied-index", ep.appliedi),
			zap.Uint64("incoming-leader-snapshot-index", apply.snapshot.Metadata.Index),
			zap.Uint64("incoming-leader-snapshot-term", apply.snapshot.Metadata.Term),
			zap.Duration("elapsed", time.Since(start)))
	}(time.Now())

	if apply.snapshot.Metadata.Index <= ep.appliedi {
		s.logger.Panic("unexpected leader snapshot from outdated index",
			zap.Uint64("snapshot-index", ep.snapi),
			zap.Uint64("applied-index", ep.appliedi),
			zap.Uint64("incoming-leader-snapshot-index", apply.snapshot.Metadata.Index),
			zap.Uint64("incoming-leader-snapshot-term", apply.snapshot.Metadata.Term))
	}

	// wait for raftNode to persist snapshot onto the disk
	<-apply.notifyc

	// TODO: handle new storage

	// TODO: handle cluster members

	ep.appliedt = apply.snapshot.Metadata.Term
	ep.appliedi = apply.snapshot.Metadata.Index
	ep.snapi = ep.appliedi
	ep.confState = apply.snapshot.Metadata.ConfState
}

func (s *Server) applyEntries(ep *etcdProgress, apply *apply) {
	if len(apply.entries) == 0 {
		return
	}

	firstIndex := apply.entries[0].Index
	if firstIndex > ep.appliedi+1 {
		s.logger.Panic("unexpected committed entry index",
			zap.Uint64("current-applied-index", ep.appliedi),
			zap.Uint64("first-committed-entry-index", firstIndex))
	}

	var ents []raftpb.Entry
	if ep.appliedi+1-firstIndex < uint64(len(apply.entries)) {
		ents = apply.entries[ep.appliedi+1-firstIndex:]
	}

	if len(ents) == 0 {
		return
	}

	var shouldStop bool
	if ep.appliedt, ep.appliedi, shouldStop = s.apply(ents, &ep.confState); shouldStop {
		go s.stopWithDelay(10*100*time.Millisecond,
			errors.Errorf("the member has been permanetly removed from the cluster"))
	}
}

func (s *Server) stopWithDelay(d time.Duration, err error) {
	select {
	case <-time.After(d):
	case <-s.doneCh:
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

// apply takes entries received from Raft (after it has been committed)
// and applies them to the current state of the EtcdServer.
// The given entries should not be empty
func (s *Server) apply(
	ents []raftpb.Entry,
	confState *raftpb.ConfState,
) (appliedTerm, appliedIndex uint64, shouldStop bool) {
	for i := range ents {
		e := ents[i]
		switch e.Type {
		case raftpb.EntryNormal:
			s.applyEntryNormal(&e)
			s.setAppliedIndex(e.Index)
			s.setTerm(e.Term)

		case raftpb.EntryConfChange:
			// set the consistent index of current executing entry
			if e.Index > s.consistentIndex.ConsistentIndex() {
				s.consistentIndex.SetConsistentIndex(e.Index)
			}

			var cc raftpb.ConfChange
			pbutil.MustUnmarshal(&cc, e.Data)
			removedSelf, err := s.applyConfChange(cc, confState)
			s.setAppliedIndex(e.Index)
			s.setTerm(e.Term)
			shouldStop = shouldStop || removedSelf
			s.wait.Trigger(cc.ID, &confChangeResponse{s.cluster.Members(), err})

		default:
			// It must be either EntryNormal or EntryConfChange
			s.logger.Panic("unknown entry type",
				zap.String("type", e.Type.String()))
		}
	}

	return appliedTerm, appliedIndex, shouldStop
}

// applyConfChange applies a ConfChange to the server.
// It is only invoked with a confChange that has already
// passed through Raft
func (s *Server) applyConfChange(cc raftpb.ConfChange, confState *raftpb.ConfState) (bool, error) {
	if err := s.cluster.ValidateConfigurationChange(cc); err != nil {
		cc.NodeID = raft.None
		s.raftNode.ApplyConfChange(cc)
		return false, err
	}

	*confState = *s.raftNode.ApplyConfChange(cc)
	switch cc.Type {
	case raftpb.ConfChangeAddNode, raftpb.ConfChangeAddLearnerNode:
		confChangeContext := new(membership.ConfigChangeContext)
		if err := json.Unmarshal(cc.Context, confChangeContext); err != nil {
			s.logger.Panic("failed to unmarshal member",
				zap.Error(err))
		}

		if cc.NodeID != uint64(confChangeContext.Member.ID) {
			s.logger.Panic("got different member ID",
				zap.String("member-id-from-config-change-entry", types.ID(cc.NodeID).String()),
				zap.String("member-id-from-message", confChangeContext.Member.ID.String()))
		}

		if confChangeContext.IsPromote {
			s.cluster.PromoteMember(uint64(confChangeContext.Member.ID))
		} else {
			s.cluster.AddMember(&confChangeContext.Member)

			if confChangeContext.Member.ID != s.ID {
				s.raftNode.transport.AddPeer(types.ID(confChangeContext.Member.ID), confChangeContext.Addresses)
			}
		}

		// update the isLearner metric when this server id is equal to the id in raft member confChange
		if confChangeContext.Member.ID == s.ID {
			if cc.Type == raftpb.ConfChangeAddLearnerNode {
				s.isLearner.Set(1)
			} else {
				s.isLearner.Set(0)
			}
		}

	case raftpb.ConfChangeRemoveNode:
		id := manta.ID(cc.NodeID)
		s.cluster.RemoveMember(uint64(id))
		if id == s.ID {
			return true, nil
		}

		s.raftNode.transport.RemovePeer(types.ID(id))

	case raftpb.ConfChangeUpdateNode:
		m := &membership.Member{}
		if err := json.Unmarshal(cc.Context, m); err != nil {
			s.logger.Panic("failed to unmarshal member",
				zap.Error(err))
		}

		if cc.NodeID != uint64(m.ID) {
			s.logger.Panic("got different member ID",
				zap.String("member-id-from-config-change-entry", types.ID(cc.NodeID).String()),
				zap.String("member-id-from-message", m.ID.String()))
		}

		s.cluster.UpdateAttributes(m.ID, m.Attributes)
		if m.ID != s.ID {
			s.raftNode.transport.UpdatePeer(types.ID(m.ID), m.Addresses)
		}
	}

	return false, nil
}

func (s *Server) setAppliedIndex(index uint64) {
	atomic.StoreUint64(&s.appliedIndex, index)
}

func (s *Server) setTerm(term uint64) {
	atomic.StoreUint64(&s.term, term)
}

func (s *Server) getCommittedIndex() uint64 {
	return atomic.LoadUint64(&s.committedIndex)
}

func (s *Server) setCommittedIndex(i uint64) {
	atomic.StoreUint64(&s.committedIndex, i)
}

func (s *Server) isLeader() bool {
	return uint64(s.ID) == s.getLead()
}

func (s *Server) setLead(v uint64) {
	atomic.StoreUint64(&s.lead, v)
}

func (s *Server) getLead() uint64 {
	return atomic.LoadUint64(&s.lead)
}

// applyEntryNormal apples an EntryNormal type raftpb request to the EtcdServer
func (s *Server) applyEntryNormal(entry *raftpb.Entry) {
	shouldApplyV3 := false
	index := s.consistentIndex.ConsistentIndex()
	if entry.Index > index {
		// set the consistent index of current executing entry
		s.consistentIndex.SetConsistentIndex(entry.Index)
		shouldApplyV3 = true
	}

	s.logger.Debug("apply entry normal",
		zap.Uint64("consistent-index", index),
		zap.Uint64("entry-index", entry.Index),
		zap.Bool("should-applyV3", shouldApplyV3))

	// raft state machine may generate noop entry when leader confirmation.
	// skip it in advance to avoid some potential bug in the future
	if len(entry.Data) == 0 {
		select {
		case s.forceVersionC <- struct{}{}:
		default:
		}

		return
	}

	var raftReq pb.InternalRaftRequest
	if !pbutil.MaybeUnmarshal(&raftReq, entry.Data) {
		var r pb.Request
		rp := &r
		pbutil.MustUnmarshal(rp, entry.Data)
		// s.w.Trigger(r.ID, s.applyV2Request((*RequestV2)(rp)))
		return
	}
}
