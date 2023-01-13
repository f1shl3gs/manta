package raftstore

import (
	"context"

	"github.com/f1shl3gs/manta/raftstore/membership"
)

type RaftMaintanceService interface {
	Members() []membership.Member

	Join(ctx context.Context, addr string) error
}
