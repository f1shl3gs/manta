package raft

import (
	"fmt"
	"github.com/f1shl3gs/manta/pkg/fsutil"
	"github.com/f1shl3gs/manta/raftstore/raft/cindex"
	"github.com/f1shl3gs/manta/raftstore/raft/membership"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.etcd.io/etcd/server/v3/wal"
	"go.etcd.io/etcd/server/v3/wal/walpb"
	"go.uber.org/zap"
	"os"
	"strings"
)

const (
	// The max throughput of etcd will not exceed 100MB/s (100K * 1KB value).
	// Assuming the RTT is around 10ms, 1MB max size is large enough.
	maxSizePerMsg = 1 * 1024 * 1024
	// Never overflow the rafthttp buffer, which is 4096.
	// TODO: a better const?
	maxInflightMsgs = 4096 / 8
)

func New(cf *Config, logger *zap.Logger, db *bolt.DB) (raft.Node, error) {
	var (
		w  *wal.WAL
		n  raft.Node
		s  *raft.MemoryStorage
		cl *membership.Cluster
	)

	if cf.MaxRequestBytes > maxRequestBytes {
		logger.Warn(
			"exceeded recommended request limit, use max",
			zap.Uint("max-request-bytes", maxRequestBytes),
		)
	}

	// make sure dirs is created and writable
	if err := fsutil.TouchDirAll(cf.DataDir); err != nil {
		return nil, errors.Wrapf(err, "cannot create/access data directory %s", cf.DataDir)
	}

	if err := fsutil.TouchDirAll(cf.WalDir()); err != nil {
		return nil, errors.Wrapf(err, "cannot access/create wal directory %s", cf.WalDir())
	}

	haveWal := wal.Exist(cf.WalDir())

	if err := fsutil.TouchDirAll(cf.SnapDir()); err != nil {
		return nil, errors.Wrapf(err, "cannot create/access snap directory %s", cf.SnapDir())
	}

	// cleanup temp snaps
	if entries, err := os.ReadDir(cf.SnapDir()); err != nil {
		for _, entry := range entries {
			if !strings.HasPrefix(entry.Name(), "tmp") {
				continue
			}

			if err = os.Remove(entry.Name()); err != nil {
				logger.Error("failed to remove temp file in snapshot directory",
					zap.String("file", entry.Name()),
					zap.Error(err))
			}
		}
	}

	var (
		ss = snap.New(logger, cf.SnapDir())
		ci = cindex.New()
		id uint64
	)

	if haveWal {
		if err := fsutil.TouchDirAll(cf.MemberDir()); err != nil {
			return nil, fmt.Errorf("cannot write to access memeber directory: %v", err)
		}

		// Find a snapshot to start/restart a raft node
		walSnaps, err := wal.ValidSnapshotEntries(logger, cf.WalDir())
		if err != nil {
			return nil, err
		}

		// snapshot files can be orphaned if server crashes after writing them but
		// before writing the corresponding wal log entries
		snapshot, err := ss.LoadNewestAvailable(walSnaps)
		if err != nil {
			return nil, err
		}

		if snapshot != nil {
			if db, err = recoverSnapshot(cf, db, *snapshot); err != nil {
				logger.Panic("failed to recover from snapshot", zap.Error(err))
			}

			logger.Info("recovered from snapshot")
		} else {
			logger.Info("no snapshot found. recovering WAL from scratch!")
		}

		id, cl, n, s, w = restartNode(logger, cf, snapshot)

	} else {

	}
}

func restartNode(logger *zap.Logger, cf *Config, snapshot *raftpb.Snapshot) (uint64, *membership.Cluster, raft.Node, *raft.MemoryStorage, *wal.WAL) {
	var walSnap walpb.Snapshot

	if snapshot != nil {
		walSnap.Index, walSnap.Term = snapshot.Metadata.Index, snapshot.Metadata.Term
	}
	w, id, cid, st, ents := readWAL(logger, cf.WalDir(), walSnap)

	logger.Info(
		"restarting local member",
		zap.String("cluster-id", internal.IDToString(cid)),
		zap.String("local-member", internal.IDToString(id)),
		zap.Uint64("commit-index", st.Commit),
	)

	cl := membership.New()
	cl.SetID(id, cid)
	s := raft.NewMemoryStorage()
	if snapshot != nil {
		s.ApplySnapshot(*snapshot)
	}
	s.SetHardState(st)
	s.Append(ents)

	c := &raft.Config{
		ID:              id,
		ElectionTick:    cf.ElectionTicks,
		HeartbeatTick:   1,
		Storage:         s,
		MaxSizePerMsg:   maxSizePerMsg,
		MaxInflightMsgs: maxInflightMsgs,
		CheckQuorum:     true,
		PreVote:         cf.PreVote,
		Logger:          newRaftLoggerZap(logger),
	}

	n := raft.RestartNode(c)

	return id, cl, n, s, w
}
