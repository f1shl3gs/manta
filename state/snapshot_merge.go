package state

import (
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
)

// createMergedSnapshotMessage creates a snapshot message that contains:
// raft status (term, conf), a snapshot of v2 store inside raft.Snapshot as
// []byte, a snapshot of v3 KV in the top level message as ReadCloser.
func (s *Server) createMergedSnapshotMessage(m raftpb.Message, snapTerm, snapIndex uint64, confState raftpb.ConfState) snap.Message {
	// storage snapshot

	snapshot := raftpb.Snapshot{
		Metadata: raftpb.SnapshotMetadata{
			Index:     snapIndex,
			Term:      snapTerm,
			ConfState: confState,
		},
		Data: []byte{},
	}

	m.Snapshot = snapshot

	return *snap.NewMessage(m, nil, 0)
}
