package raftstore

import (
	"context"
	"github.com/f1shl3gs/manta/raftstore/pb"
	"go.etcd.io/raft/v3/raftpb"

	"github.com/f1shl3gs/manta/raftstore/membership"
)

type RaftMaintanceService interface {
	Members() []membership.Member

	// Join add a node to Raft cluster
	Join(ctx context.Context, addr string) error

	// Leave removes member of Raft cluster
	Leave(ctx context.Context, id uint64) error
}

var (
	done = &pb.Done{}
)

// Send implement RaftServer
func (s *Store) Send(ctx context.Context, msg *raftpb.Message) (*pb.Done, error) {
	err := s.raftNode.Step(ctx, *msg)
	if err != nil {
		return nil, err
	}

	return done, nil
}
