package raft

import (
	"github.com/f1shl3gs/manta/raftstore/membership"
	"github.com/f1shl3gs/manta/raftstore/wal"
	"go.etcd.io/raft/v3"
	"go.uber.org/zap"
)

func defaultRaftConfig() *raft.Config {
	return &raft.Config{
		ElectionTick:    20, // 2s
		HeartbeatTick:   1,
		MaxInflightMsgs: 256,

		// 512 KB should be enough for most txn.
		MaxSizePerMsg:  512 * 1024,
		ReadOnlyOption: raft.ReadOnlySafe,

		// When a disconnected node joins back, it forces a leader change,
		// as it starts with a higher term, as described in Raft thesis
		// (not the paper) in section 9.6. This setting can avoid that by
		// only increasing the term, if the node has a good chance of
		// becoming the leader.
		PreVote: true,
	}
}

func newNode(cf *Config, logger *zap.Logger) (*raftNode, error) {
	store, err := wal.Init(cf.DataDir, logger)
	if err != nil {
		return nil, err
	}

	nid := store.Uint(wal.RaftId)
	if nid == 0 {
		nid = membership.GenerateID(cf.Listen)
	}

	snap, err := store.Snapshot()
	if err != nil {
		return nil, err
	}

	cfg := &raft.Config{
		ID:              nid,
		ElectionTick:    1000 / 100,
		HeartbeatTick:   1,
		MaxInflightMsgs: 256,

		// Setting applied to the first index in the raft log, so it does not
		// derive it separately, thus avoiding a crash when the Applied is set
		// to below snapshot index by Raft.
		//
		// In case this is a new Raft log, first would be 1, and therefore
		// Applied would be zero, hence meeting the condition by the library
		// that Applied should only be set during a restart.
		Applied: snap.Metadata.Index,

		// Storage is the storage for raft. it is used to store wal and snapshots.
		Storage: store,

		// MaxSizePerMsg specifies the maximum aggregate byte size of Raft
		// log entries that a leader will send to followers in a single MsgApp.
		MaxSizePerMsg: 64 << 10, // 64KB should allow more txn

		// MaxCommittedSizePerReady specifies the maximum aggregate
		// byte size of the committed log entries which a node will receive in a
		// single Ready.
		MaxCommittedSizePerReady: 64 << 20, // 64MB

		PreVote: true,
		Logger:  newRaftLoggerZap(logger.Named("raft")),
	}

	return &raftNode{}, nil
}

func (r *raftNode) initAndStart() error {
	restart := r.storage.NumEntries() > 1

	if restart {
		r.logger.Info("restarting node from wal and snapshot")

	} else {
		r.logger.Info("start a brand new raft node")
	}
}
