package raftstore

import (
	"context"

	"github.com/f1shl3gs/manta/raftstore/internal"
	"github.com/f1shl3gs/manta/raftstore/transport"

	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.uber.org/zap"
)

// Send implement transport.RaftServer
func (s *Store) Send(ctx context.Context, m *raftpb.Message) (*transport.Done, error) {
	if s.cluster.IDRemoved(m.From) {
		s.logger.Warn("rejected Raft message from removed member",
			zap.String("local-member-id", internal.IDToString(s.id)),
			zap.String("from", internal.IDToString(m.From)))
		return nil, ErrRejectFromRemovedMember
	}

	return nil, s.raftNode.Step(ctx, *m)
}

// SendSnapshot implement transport.RaftServer
func (s *Store) SendSnapshot(server transport.Raft_SendSnapshotServer) error {
	// TODO: there are something need pay attention
	// 1. do we need to save it to disk, and trigger a GC?
	// 2. restore the kv engine right here?
	panic("implement me")
}
