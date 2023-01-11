package raft

import (
	"context"

	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
)

type Transporter interface {
	Start(ctx context.Context) error
	Stop() error

	// Send sends out the given messages to the remote peers.
	Send(msgs []raftpb.Message)

	SendSnapshot(m snap.Message)

	AddRemote()
}
