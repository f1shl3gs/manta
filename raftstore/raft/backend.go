package raft

import "go.etcd.io/raft/v3/raftpb"

type Backend interface {
	ApplyEntries(entries []raftpb.Entry)

	ApplySnapshot(snapshot raftpb.Snapshot)

	ConsistentIndex() uint64

	SetConsistentIndex(index, term uint64)
}
