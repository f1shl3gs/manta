package raft

import (
	"context"

	"go.etcd.io/raft/v3/raftpb"
)

type Transporter interface {
	Start(ctx context.Context) error
	Stop() error

	// Send sends out the given messages to the remote peers.
	Send(msgs []raftpb.Message)

	AddRemote()
}
