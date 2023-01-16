package raftstore

import (
	"context"

	"github.com/f1shl3gs/manta/raftstore/membership"
)

type RaftMaintanceService interface {
	Members() []membership.Member

	// Join add a node to Raft cluster
	Join(ctx context.Context, addr string) error

	// Leave removes member of Raft cluster
	Leave(ctx context.Context, id uint64) error
}
