package state

import (
	"context"
	"hash/fnv"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/etcdserver/stats"
	"github.com/f1shl3gs/manta/pkg/snowflake"
	"github.com/f1shl3gs/manta/state/cindex"
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
	maxInFlightMsgSnap = 16

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

	ID     types.ID
	logger *zap.Logger
	node   raft.Node

	idGen *snowflake.Generator

	// raft
	raftNode        *raftNode
	raftStorage     raft.Storage
	consistentIndex cindex.ConsistentIndexer

	// example implement
	transport     *rafthttp.Transport
	wait          wait.Wait
	applyWait     wait.WaitTime
	doneCh        chan struct{}
	stopCh        chan struct{}
	stopping      chan struct{}
	cancel        context.CancelFunc
	ctx           context.Context
	readwaitC     chan struct{}
	readNotifier  *notifier
	leaderChanged chan struct{}
	lead          uint64

	// stats
	stats       *stats.ServerStats
	leaderStats *stats.LeaderStats

	committedIndex uint64

	// metrics
	applySnapshotInProgress prometheus.Gauge

	// unknown field
	forceVersionC chan struct{}
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

func (s *Server) run() {
	sn, err := s.raftStorage.Snapshot()
	if err != nil {
		s.logger.Panic("failed to get snapshot from raft storage",
			zap.Error(err))
	}

	// asynchronously accept apply packets, dispatch progress in-order
	sched := schedule.NewFIFOScheduler()

	var (
		mtx   sync.RWMutex
		syncC <-chan time.Time
	)

	setSyncC := func(ch <-chan time.Time) {
		mtx.Lock()
		syncC = ch
		mtx.Unlock()
	}

	getSyncC := func() <-chan time.Time {
		mtx.RLock()
		ch := syncC
		mtx.RUnlock()

		return ch
	}

	rh := &raftReadyHandler{
		getLead: func() (lead uint64) {
			return s.getLead()
		},
		updateLead: func(lead uint64) {
			s.setLead(lead)
		},
		updateLeadership: func(newLeader bool) {
			if !s.isLeader() {
				setSyncC(nil)

				// todo: handle compactor properly
			} else {
				if newLeader {
					t := time.Now()
					s.leadTimeMu.Lock()
					s.leadElectedTime = t
					s.leadTimeMu.Unlock()
				}

				setSyncC(s.SyncTicker.C)
				// todo: handle compactor
			}

			if newLeader {
				s.leaderChangeMu.Lock()
				lc := s.leaderChanged
				s.leaderChanged = make(chan struct{})
				close(lc)
				s.leaderChangeMu.Unlock()
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

	// etcProcess ?

	for {
		select {
		case ap := <-s.raftNode.apply():
			f := func(ctx context.Context) {
				s.applyAll(&ep, &ap)
			}
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

	proposalsApplied.Set(float64(ep.appliedi))
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
		}
	}
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

	if raftReq.V2 != nil {

	}
}
