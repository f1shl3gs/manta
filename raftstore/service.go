package raftstore

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"time"

	"github.com/f1shl3gs/manta/raftstore/pb"

	"go.etcd.io/raft/v3/raftpb"
)

type ClusterService interface {
	Members() []pb.Member

	// Add add a node to Raft cluster
	Add(ctx context.Context, member pb.Member) error

	// Remove removes member of Raft cluster
	Remove(ctx context.Context, id uint64) error
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

func (s *Store) Members() []pb.Member {
	panic("not implement")
}

// Add add a node to Raft cluster
func (s *Store) Add(ctx context.Context, member pb.Member) error {
	cc := raftpb.ConfChange{
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  generateID(member.Addr),
		Context: unsafeStringToBytes(member.Addr),
	}

	return s.raftNode.ProposeConfChange(ctx, cc)
}

// Remove removes member of Raft cluster
func (s *Store) Remove(ctx context.Context, id uint64) error {
	cc := raftpb.ConfChange{
		Type:   raftpb.ConfChangeRemoveNode,
		NodeID: id,
	}

	return s.raftNode.ProposeConfChange(ctx, cc)
}

// generateID generate a new node id with address
func generateID(addr string) uint64 {
	b := []byte(addr)
	b = append(b, []byte(time.Now().String())...)

	hash := sha1.Sum(b)

	return binary.BigEndian.Uint64(hash[:8])
}
