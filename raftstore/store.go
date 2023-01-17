package raftstore

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/raftstore/membership"
	"github.com/f1shl3gs/manta/raftstore/pb"
	"github.com/f1shl3gs/manta/raftstore/transport"
	"github.com/f1shl3gs/manta/raftstore/wal"
	"github.com/prometheus/client_golang/prometheus"
	bolt "go.etcd.io/bbolt"
	"go.etcd.io/raft/v3"
	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// initialMmapSize is the initial size of the mmapped region. Setting this
	// larger than the potential max db size can prevent writer from blocking reader.
	//
	// This only works for linux
	initialMmapSize = 128 * 1024 * 1024 // 128MiB

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
	metaBucket = []byte("__meta")

	consistentIndexKey = []byte("consistentIndex")
)

type Store struct {
	logger  *zap.Logger
	readyCh chan struct{}
	db      atomic.Pointer[bolt.DB]
	cluster *membership.Cluster
	wait    *wait
	idGen   *idGenerator

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

	ds, err := wal.Init(cf.DataDir, logger)
	if err != nil {
		return nil, err
	}

	store := &Store{
		logger:        logger,
		wait:          newWait(),
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

		leaderChanges: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "manta",
			Subsystem: "raftstore",
			Name:      "leader_changes_seen_total",
			Help:      "The number of leader changes seen.",
		}),
		hasLeader: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "manta",
			Subsystem: "raftstore",
			Name:      "has_leader",
			Help:      "Whether or not a leader exists. 1 is existence, 0 is not.",
		}),
		isLeader: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "manta",
			Subsystem: "raftstore",
			Name:      "is_leader",
			Help:      "Whether or not this member is a leader. 1 if is, 0 otherwise.",
		}),
		slowReadInex: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "manta",
			Subsystem: "raftstore",
			Name:      "slow_read_indexes_total",
			Help:      "The total number of pending read indexes not in sync with leader's or timed out read index requests.",
		}),
		readIndexFailed: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "manta",
			Subsystem: "raftstore",
			Name:      "read_indexes_failed_total",
			Help:      "The total number of failed read indexes seen.",
		}),
	}

	if db, err := openDB(cf); err != nil {
		return nil, err
	} else {
		store.db.Store(db)
	}

	// TODO: what if there is multiple term in ds
	term, index, err := store.consistentIndex()
	if err != nil {
		panic("read consistent index failed")
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
		Applied: index,

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
	if term == 0 && index == 0 {
		logger.Info("start a brand new raft cluster")

		rcf.ID = membership.GenerateID(cf.Listen)
		peers := []raft.Peer{
			{
				ID:      rcf.ID,
				Context: unsafeStringToBytes(cf.Listen),
			},
		}

		trans := transport.New(logger)
		for _, peer := range peers {
			err = trans.AddPeer(peer.ID, unsafeBytesToString(peer.Context))
			if err != nil {
				return nil, err
			}
		}

		store.idGen = newGenerator(uint16(rcf.ID), time.Now())
		store.transport = trans
		store.raftNode = raft.StartNode(rcf, peers)
		return store, nil
	}

	// restart raft node

	// replaying unapplied entry to backend and restore cluster from it
	lastIndex, err := ds.LastIndex()
	if err != nil {
		return nil, err
	}

	entries, err := ds.Entries(index, lastIndex, lastIndex-index)
	if err != nil {
		return nil, err
	}

	for _, ent := range entries {
		if ent.Index <= index {
			panic("replying entry smaller than current backend")
		}
	}

	return store, nil
}

func (s *Store) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		s.leaderChanges,
		s.hasLeader,
	}
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

	db, err := bolt.Open(path, 0600, opt)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(metaBucket)
		return err
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (s *Store) consistentIndex() (term, index uint64, err error) {
	err = s.db.Load().View(func(tx *bolt.Tx) error {
		b := tx.Bucket(metaBucket)
		value := b.Get(consistentIndexKey)

		// first setup
		if len(value) == 0 {
			return nil
		}

		if len(value) != 16 {
			term = binary.BigEndian.Uint64(value)
			index = binary.BigEndian.Uint64(value[8:])
		} else {
			panic(fmt.Sprintf("consistent value is not 16 bytes"))
		}

		return nil
	})

	return
}

func (s *Store) Run(ctx context.Context) {
	var batched uint64

	go s.startRaftLoop()
	go s.syncBolt(ctx)
	go s.linearizableReadLoop(ctx)

	// apply loop
	for {
		select {
		case <-ctx.Done():
			return

		case apply := <-s.applyCh:
			db := s.db.Load()

			err := db.Update(func(tx *bolt.Tx) error {
				for _, ent := range apply.entries {
					switch ent.Type {
					case raftpb.EntryNormal:
						s.applyNormal(tx, &ent)
						s.appliedIndex.Store(ent.Index)

					case raftpb.EntryConfChange:
						// TODO

						s.appliedIndex.Store(ent.Index)
					default:
						s.logger.Error("trying to apply unknown entry type",
							zap.Int32("type", int32(ent.Type)))
					}
				}

				return nil
			})
			if err != nil {
				s.logger.Fatal("error while applying to boltdb",
					zap.Error(err))
			}

			// wait for the raft routine to finish the disk writes before triggering
			// a snapshot. or applied index might be greater than the last index in
			// raft storage. since the raft routine might be slower than toApply
			// routine.
			<-apply.notifyCh

			batched += uint64(len(apply.entries))
			if batched > batchLimit {
				err := s.db.Load().Sync()
				if err != nil {
					s.logger.Fatal("sync boltdb failed",
						zap.Error(err))
				}

				batched = 0
			}
		}
	}
}

func (s *Store) applyNormal(tx *bolt.Tx, ent *raftpb.Entry) {
	var req pb.InternalRequest

	err := req.Unmarshal(ent.Data)
	if err != nil {
		s.logger.Error("unmarshal internal request failed",
			zap.Error(err))
		return
	}

	err = func() error {
		if txn := req.GetTxn(); txn != nil {
			return applyTxn(tx, txn)
		}

		if cb := req.GetCreateBucket(); cb != nil {
			_, err = tx.CreateBucket(cb.Name)
			return err
		}

		if d := req.GetDeleteBucket(); d != nil {
			return tx.DeleteBucket(d.Name)
		}

		return errors.New("empty internal request")
	}()

	s.wait.Trigger(req.ID, err)
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

func (s *Store) syncBolt(ctx context.Context) {
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return
	case <-ticker.C:
		err := s.db.Load().Sync()
		if err != nil {
			s.logger.Fatal("sync boltdb failed",
				zap.Error(err))
		}

		ticker.Reset(batchInterval)
	}
}
