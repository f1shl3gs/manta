package raftstore

import (
	"context"
	"github.com/f1shl3gs/manta/raftstore/transport"
	"go.etcd.io/raft/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	logger *zap.Logger
	raft   raft.Node
}

func (s *server) Send(ctx context.Context, batch *transport.Batch, opts ...grpc.CallOption) (*transport.SendResponse, error) {
	for _, m := range batch.Msgs {
		s.raft.Step(ctx, m)
	}
}
