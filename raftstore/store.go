package raftstore

import (
	"context"
	"encoding/binary"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/pkg/fsutil"
	"github.com/f1shl3gs/manta/raftstore/pb"
	"github.com/f1shl3gs/manta/raftstore/transport"
	"github.com/f1shl3gs/manta/raftstore/wal"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	bolt "go.etcd.io/bbolt"
	"go.etcd.io/raft/v3"
	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
)

const (
	// initialMmapSize is the initial size of the mmapped region. Setting this
	// larger than the potential max db size can prevent writer from blocking reader.
	//
	// This only works for linux
	initialMmapSize = 2 * 1024 * 1024 * 1024 // 2 GiB

	// max number of in-flight snapshot messages server allows to have
	// This number is more than enough for most clusters with 5 machines.
	maxInflightMsgSnap = 16

	electionMs   = 1000
	tickMs       = 100
	electionTick = electionMs / tickMs

	batchLimit    = 10000
	batchInterval = 100 * time.Millisecond
)

var (
	membershipBucket = []byte("__membership")
)

type Store struct {
	self      pb.Member
	logger    *zap.Logger
	readyCh   chan struct{}
	db        atomic.Pointer[bolt.DB]
	wait      *wait[error]
	idGen     *idGenerator
	confState atomic.Pointer[raftpb.ConfState]

	// read routine notifies server that it waits for reading by
	// sending an emtpy struct to readWaitCh
	readMtx    sync.RWMutex
	readWaitCh chan struct{}
	// readNotifier is used to notify the read reoutine that it can
	// process the request when there is no error
	readNotifier *errNotifier

	// leaderChanged is used to notify the linearizable read loop to drop the old read requests.
	leaderChanged *notifier

	firstCommitInTerm *notifier
	applyWait         *waitTime

	// raft staff
	raftNode    raft.Node
	raftStorage *wal.DiskStorage
	// a chan to send/receive snapshot
	msgSnapCh chan raftpb.Message
	// a chan to send out apply
	applyCh chan toApply
	// a chan to send out readState
	readStateCh chan raft.ReadState
	// utility
	tickMu *sync.Mutex
	ticker *time.Ticker
	// contention detectors for raft heartbeat message
	td        *TimeoutDetector
	stopped   chan struct{}
	done      chan struct{}
	transport *transport.Transporter

	// raft stats
	lead           atomic.Uint64
	committedIndex atomic.Uint64
	appliedIndex   atomic.Uint64

	// metrics
	leaderChanges   prometheus.Counter
	hasLeader       prometheus.Gauge
	isLeader        prometheus.Gauge
	slowReadInex    prometheus.Counter
	readIndexFailed prometheus.Counter
}

func New(cf *Config, logger *zap.Logger) (*Store, error) {
	logger = logger.Named("raftstore")

	err := fsutil.TouchDirAll(cf.DataDir)
	if err != nil {
		return nil, err
	}

	ds, err := wal.Init(cf.DataDir, logger)
	if err != nil {
		return nil, err
	}

	store := &Store{
		logger:        logger,
		wait:          newWait[error](),
		readWaitCh:    make(chan struct{}, 1),
		leaderChanged: newNotifier(),
		tickMu:        new(sync.Mutex),
		ticker:        time.NewTicker(heartbeat),
		readyCh:       make(chan struct{}),
		// set up contention detectors for raft heatbeat message.
		// expect to send a heartbeat within 2 heartbeat intervals.
		td:                NewTimeoutDetector(2 * heartbeat),
		readStateCh:       make(chan raft.ReadState, 1),
		msgSnapCh:         make(chan raftpb.Message, maxInflightMsgSnap),
		applyCh:           make(chan toApply),
		stopped:           make(chan struct{}),
		done:              make(chan struct{}),
		raftStorage:       ds,
		applyWait:         newWaitTime(),
		firstCommitInTerm: newNotifier(),
		readNotifier:      newErrNotifier(),
		transport:         transport.New(logger),

		leaderChanges: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "leader_changes_seen_total",
			Help:      "The number of leader changes seen.",
		}),
		hasLeader: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "has_leader",
			Help:      "Whether or not a leader exists. 1 is existence, 0 is not.",
		}),
		isLeader: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "is_leader",
			Help:      "Whether or not this member is a leader. 1 if is, 0 otherwise.",
		}),
		slowReadInex: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "slow_read_indexes_total",
			Help:      "The total number of pending read indexes not in sync with leader's or timed out read index requests.",
		}),
		readIndexFailed: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "read_indexes_failed_total",
			Help:      "The total number of failed read indexes seen.",
		}),
	}

	db, err := openDB(cf)
	if err != nil {
		return nil, err
	}

	store.db.Store(db)

	appliedIndex, err := ds.Checkpoint()
	if err != nil {
		logger.Fatal("read checkpoint failed",
			zap.Error(err))
	}

	rcf := &raft.Config{
		ElectionTick:    electionMs / tickMs,
		HeartbeatTick:   1,
		MaxInflightMsgs: 256,
		// Setting applied to the first index in the raft log, so it does not
		// derive it separately, thus avoiding a crash when the Applied is set
		// to below snapshot index by Raft.
		//
		// In case this is a new Raft log, first would be 1, and therefore
		// Applied would be zero, hence meeting the condition by the library
		// that Applied should only be set during a restart.
		Applied: appliedIndex,

		// Storage is the storage for raft. it is used to store wal and snapshots.
		Storage: ds,

		// MaxSizePerMsg specifies the maximum aggregate byte size of Raft
		// log entries that a leader will send to followers in a single MsgApp.
		MaxSizePerMsg: 64 << 10, // 64KB should allow more txn

		// MaxCommittedSizePerReady specifies the maximum aggregate
		// byte size of the committed log entries which a node will receive in a
		// single Ready.
		MaxCommittedSizePerReady: 64 << 20, // 64MB

		// When a disconnected node joins back, it forces a leader change,
		// as it starts with a higher term, as described in Raft thesis
		// (not the paper) in section 9.6. This setting can avoid that by
		// only increasing the term, if the node has a good chance of
		// becoming the leader.
		PreVote: true,
		Logger:  newRaftLoggerZap(logger),
	}

	// if node never start before, we don't need to replay
	if appliedIndex == 0 {
		logger.Info("start a brand new raft cluster")

		rcf.ID = generateID(cf.Listen)
		ds.SetNodeID(rcf.ID)
		peers := []raft.Peer{
			{
				ID:      rcf.ID,
				Context: unsafeStringToBytes(cf.Listen),
			},
		}

		for _, peer := range peers {
			err = store.transport.AddPeer(peer.ID, unsafeBytesToString(peer.Context))
			if err != nil {
				return nil, err
			}
		}

		store.self.ID = rcf.ID
		store.self.Addr = cf.Listen
		store.idGen = newGenerator(uint16(rcf.ID), time.Now())
		store.raftNode = raft.StartNode(rcf, peers)
		return store, nil
	}

	// restart raft node
	rcf.ID = ds.NodeID()
	err = store.db.Load().View(func(tx *bolt.Tx) error {
		b := tx.Bucket(membershipBucket)

		return b.ForEach(func(k, v []byte) error {
			id := binary.BigEndian.Uint64(k)
			addr := unsafeBytesToString(v)

			return store.transport.AddPeer(id, addr)
		})
	})
	if err != nil {
		return nil, errors.Errorf("load membership from storage failed, %s", err)
	}

	logger.Info("restart raft node")

	// if !raft.IsEmptySnap(sp) {
	// It is important that we pick up the conf state here.
	// Otherwise, we'll lose the store conf state, and it
	// would get overwritten with an empty state when a new
	// snapshot is taken. This causes a node to just hang
	// on restart, because it finds a zero-member Raft group.

	// TODO: set confState
	// }

	store.appliedIndex.Store(appliedIndex)
	store.self.ID = rcf.ID
	store.self.Addr = cf.Listen
	store.idGen = newGenerator(uint16(rcf.ID), time.Now())
	store.raftNode = raft.RestartNode(rcf)

	/*
		Note: without this restart will hanging forever, and i notice logs from raftexample when restart it

			raft2023/01/19 04:24:57 INFO: 1 switched to configuration voters=()
			raft2023/01/19 04:24:57 INFO: 1 became follower at term 3
			raft2023/01/19 04:24:57 INFO: newRaft 1 [peers: [], term: 3, commit: 1, applied: 0, lastindex: 3, lastterm: 3]
			raft2023/01/19 04:24:57 INFO: 1 switched to configuration voters=(1)
			raft2023/01/19 04:24:58 INFO: 1 is starting a new election at term 3
			raft2023/01/19 04:24:58 INFO: 1 became candidate at term 4
			raft2023/01/19 04:24:58 INFO: 1 received MsgVoteResp from 1 at term 4
			raft2023/01/19 04:24:58 INFO: 1 became leader at term 4
			raft2023/01/19 04:24:58 INFO: raft.node: 1 elected leader 1 at term 4

		raft.StartNode has no voters at first,

			raft2023/01/19 04:24:57 INFO: 1 switched to configuration voters=()

		after it call raft.NewRawNode, rawNode.Bootstrap will be called which
		call applyConfChange, then we can see that there we have one voter popout
		and finally cluster can elect a leader.

			raft2023/01/19 04:24:57 INFO: 1 switched to configuration voters=(1)

		But we use disk-based state machine, so we don't need to replay WAL like raftexample does,
		therefore the voters part will be empty and cluster can never elect a leader.
	*/
	for id, peer := range store.transport.Peers() {
		cc := raftpb.ConfChange{Type: raftpb.ConfChangeAddNode, NodeID: id, Context: unsafeStringToBytes(peer)}
		cs := store.raftNode.ApplyConfChange(cc)
		store.confState.Store(cs)
	}

	return store, nil
}

// openDB open a boltdb with default options, and setup(if none) meta buckets
// to store membership and consistent index.
func openDB(cf *Config) (*bolt.DB, error) {
	path := filepath.Join(cf.DataDir, "state.bolt")

	opt := bolt.DefaultOptions
	opt.Timeout = 3 * time.Second
	opt.InitialMmapSize = initialMmapSize
	opt.FreelistType = bolt.FreelistMapType

	// sync will be done periodly by another goroutine
	opt.NoSync = true
	opt.NoGrowSync = true

	db, err := bolt.Open(path, 0600, opt)
	if err != nil {
		return nil, err
	}

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (s *Store) propose(ctx context.Context, req pb.InternalRequest) error {
	req.ID = s.idGen.Next()
	data, err := req.Marshal()
	if err != nil {
		return err
	}

	waitCh := s.wait.Register(req.ID)
	if err = s.raftNode.Propose(ctx, data); err != nil {
		return err
	}

	// TODO: retry !?

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-waitCh:
		return nil
	}
}

func (s *Store) Run(ctx context.Context) {
	go s.startRaftLoop()
	go s.syncLoop(ctx)
	go s.linearizableReadLoop(ctx)
	go s.publish(ctx)
	go s.adjustTicks()
	go s.checkSnapshot(ctx)

	defer func() {
		s.stop()

		s.raftStorage.SetUint(wal.CheckpointIndex, s.appliedIndex.Load())
		_ = s.raftStorage.Sync()
	}()

	// apply loop
	for {
		select {
		case <-ctx.Done():
			return

		// we need this for now, may be remove this later!?
		case <-s.done:
			return

		case apply := <-s.applyCh:
			if len(apply.entries) == 0 {
				continue
			}

			var entries []raftpb.Entry

			first := apply.entries[0].Index
			applied := s.appliedIndex.Load()
			if applied+1-first < uint64(len(apply.entries)) {
				entries = apply.entries[applied+1-first:]
			}
			if len(entries) == 0 {
				continue
			}

			_, appliedIndex := s.apply(entries)

			s.appliedIndex.Store(appliedIndex)
			s.applyWait.Trigger(appliedIndex)
			// s.logger.Info("trigger apply wait", zap.Uint64("index", appliedIndex))

			// wait for the raft routine to finish the disk writes before triggering
			// a snapshot. or applied index might be greater than the last index in
			// raft storage. since the raft routine might be slower than toApply
			// routine.
			<-apply.notifyCh
		}
	}
}

// publish registers server information into the cluster.
// The function keeps attempting to register until it succeeds,
// or its server is stopped.
func (s *Store) publish(ctx context.Context) {
	req := pb.InternalRequest{
		Txn: &pb.Txn{
			Successes: []*pb.Operation{
				{
					Type:   pb.Put,
					Bucket: membershipBucket,
					Key:    uint64ToBigEndianBytes(s.self.ID),
					Value:  unsafeStringToBytes(s.self.Addr),
				},
			},
		},
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		sCtx, cancel := context.WithTimeout(ctx, s.reqTimeout())
		err := s.propose(sCtx, req)
		cancel()
		if err == nil {
			close(s.readyCh)
			s.logger.Info("published local member to cluster through raft")

			return
		}

		s.logger.Warn("failed to publish local member to cluster through raft",
			zap.Error(err))
	}
}

func (s *Store) apply(
	entries []raftpb.Entry,
) (appliedt, appliedi uint64) {
	for i := range entries {
		ent := entries[i]

		s.logger.Debug("Applying entry",
			zap.Uint64("term", ent.Term),
			zap.Uint64("index", ent.Index),
			zap.Stringer("type", ent.Type))

		switch ent.Type {
		case raftpb.EntryNormal:
			s.applyNormal(&ent)
			s.appliedIndex.Store(ent.Index)

		case raftpb.EntryConfChange:
			var cc raftpb.ConfChange
			err := cc.Unmarshal(ent.Data)
			if err != nil {
				s.logger.Panic("unmarshal config change failed",
					zap.Error(err))
			}

			s.applyConfChange(&cc)

			s.appliedIndex.Store(ent.Index)

		default:
			s.logger.Panic("trying to apply unknown entry type",
				zap.Int32("type", int32(ent.Type)))
		}

		appliedt, appliedi = ent.Term, ent.Index
	}

	return appliedt, appliedi
}

func (s *Store) applyNormal(ent *raftpb.Entry) {
	// raft state machine may generate noop entry when leader confirmation.
	// skip it in advance to avoid some potential bug in the future
	if len(ent.Data) == 0 {
		s.firstCommitInTerm.notify()
		return
	}

	var req pb.InternalRequest

	err := req.Unmarshal(ent.Data)
	if err != nil {
		s.logger.Error("unmarshal internal request failed",
			zap.Error(err))
		return
	}

	if snap := req.Snapshot; snap != nil {
		// do snapshot and clean wals
		err = s.raftStorage.CreateSnapshot(snap.Index, s.confState.Load(), nil)
		if err != nil {
			s.logger.Fatal("create snapshot failed",
				zap.Uint64("checkpoint", snap.Index),
				zap.Error(err))
		} else {
			s.logger.Info("create snapshot success",
				zap.Uint64("index", snap.Index))
		}
	} else {
		db := s.db.Load()
		if db == nil {
			s.wait.Trigger(req.ID, ErrStopped)
			return
		}

		err = s.db.Load().Update(func(tx *bolt.Tx) error {
			if txn := req.Txn; txn != nil {
				return applyTxn(tx, txn)
			}

			if cb := req.CreateBucket; cb != nil {
				_, err = tx.CreateBucket(cb.Name)
				return err
			}

			if d := req.DeleteBucket; d != nil {
				return tx.DeleteBucket(d.Name)
			}

			return errors.New("empty internal request")
		})
	}

	s.wait.Trigger(req.ID, err)
}

func (s *Store) applyConfChange(cc *raftpb.ConfChange) {
	err := s.db.Load().Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(membershipBucket)
		if err != nil {
			return err
		}

		key := uint64ToBigEndianBytes(cc.NodeID)
		switch cc.Type {
		case raftpb.ConfChangeAddNode, raftpb.ConfChangeAddLearnerNode, raftpb.ConfChangeUpdateNode:
			s.logger.Info("add/update node",
				zap.String("id", strconv.FormatUint(cc.NodeID, 16)),
				zap.ByteString("addr", cc.Context))
			return b.Put(key, cc.Context)
		case raftpb.ConfChangeRemoveNode:
			s.logger.Info("remove node",
				zap.String("id", strconv.FormatUint(cc.NodeID, 16)))
			return b.Delete(key)

		default:
			return errors.New("unsupported config change type")
		}
	})
	if err != nil {
		s.logger.Fatal("apply conf change failed",
			zap.Error(err))
	}

	cs := s.raftNode.ApplyConfChange(cc)
	s.confState.Store(cs)
}

func applyTxn(tx *bolt.Tx, txn *pb.Txn) error {
	var err error

	for _, op := range txn.Successes {
		b := tx.Bucket(op.Bucket)
		if b == nil {
			return kv.ErrBucketNotFound
		}

		if op.Type == pb.Put {
			err = b.Put(op.Key, op.Value)
		} else {
			err = b.Delete(op.Key)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) syncLoop(ctx context.Context) {
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.sync()
			ticker.Reset(batchInterval)
		}
	}
}

func (s *Store) sync() {
	start := time.Now()
	db := s.db.Load()
	if db == nil {
		return
	}

	err := db.Sync()
	elapsed := time.Since(start)
	if err != nil {
		s.logger.Fatal("sync boltdb failed",
			zap.Error(err))
	}

	s.logger.Debug("sync done", zap.Duration("elapsed", elapsed))
}
