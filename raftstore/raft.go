package raftstore

import (
	"context"
	"github.com/f1shl3gs/manta/raftstore/id"
	"time"
)

const (
	recommendedMaxRequestBytes = 4 * 1024 * 1024

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

type raftNode struct {
	id uint64

	reqIDGen *id.Generator
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
	db *bolt.DB
}

func (r *raftNode) Propose(ctx context.Context, data []byte) error {
	panic("not implement")
}

func (r *raftNode) linearizableReadNotify(ctx context.Context) error {
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
	case <-nc.c:
		return nc.err
	case <-ctx.Done():
		return ctx.Err()
	case <-s.stopCh:
		return ErrStopped
	}
}
