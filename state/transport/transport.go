package transport

import (
	"time"

	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
)

type Resolver interface {
	Resolve(id uint64) (string, error)
}

type Transport interface {
	// Send sends out the given messages to the remote peers.
	// Each message has a To field, which is an id that maps
	// to an existing peer in the transport. If the id cannot
	// be found in the transport, the message will be ignored.
	Send(m []raftpb.Message)

	// SendSnapshot sends out the given snapshot message to a
	// remote peer. The behavior of SendSnapshot is similar to
	// Send.
	SendSnapshot(m snap.Message)

	// ActiveSince returns the time that the connection with
	// the peer of the given id becomes active. If the connection
	// is active since peer was added, it returns the adding
	// time. If the connection is currently inactive, it returns
	// zero time.
	ActiveSince(id uint64) time.Time

	// ActivePeers returns the number of active peers.
	ActivePeers() int
}
