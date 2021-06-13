package raftstore

import (
	"context"
	"io"
	"sync/atomic"
	"time"

	"github.com/f1shl3gs/manta/raftstore/internal"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.uber.org/zap"
)

const (
	releaseDelayAfterSnapshot = 30 * time.Second
)

func (s *Store) createMergedSnapshotMessage(m raftpb.Message, term uint64, index uint64, confState raftpb.ConfState) snap.Message {
	panic("not implement")
}

func newSnapshotReaderCloser(logger *zap.Logger) io.ReadCloser {
	panic("not implement")
}

func (s *Store) sendMergedSnap(merged snap.Message) {
	// TODO: ctx
	ctx := context.Background()

	atomic.AddInt64(&s.inflightSnapshots, 1)

	fields := []zap.Field{
		zap.String("from", internal.IDToString(s.id)),
		zap.String("to", internal.IDToString(merged.To)),
		zap.Int64("bytes", merged.TotalSize),
	}

	start := time.Now()
	s.transporter.SendSnapshot(merged)
	s.logger.Info("sending merged snapshot", fields...)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		select {
		case ok := <-merged.CloseNotify():
			// delay releasing inflight snapshot for another 30 seconds to
			// block log compaction. If the follower still fails to catch up,
			// it is probably just too slow to catch up. We cannot avoid the
			// snapshot cycle anyway
			if ok {
				select {
				case <-time.After(releaseDelayAfterSnapshot):
				case <-ctx.Done():
				}
			}

			atomic.AddInt64(&s.inflightSnapshots, -1)
			s.logger.Info("sent merged snapshot",
				append(fields, zap.Duration("took", time.Since(start)))...)

		case <-ctx.Done():
			return
		}
	}()
}

// SnapshotStorage is used to allow for flexible implementations
// of snapshot storage and retrieval. For example, a client could
// implement a shared state store such as S3, allowing new nodes
// to restore snapshots without streaming from the leader
type SnapshotStorage interface {
	// Create is used to begin a snapshot at a give term and index
	Create(term, index uint64, rc io.ReadCloser) error
}
