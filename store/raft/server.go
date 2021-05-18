package raft

import (
	"context"
	"os"
	"time"

	"github.com/cespare/xxhash"
	"github.com/f1shl3gs/manta"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/pkg/v3/types"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"
	"go.etcd.io/etcd/server/v3/wal"
	"go.etcd.io/etcd/server/v3/wal/walpb"
	"go.uber.org/zap"
)

type Config struct {
	Listen  string `json:"listen"`
	DataDir string `json:"data_dir"`
	WALDir  string `json:"wal_dir"`
}

type commit struct {
	data       [][]byte
	applyDoneC chan<- struct{}
}

type Node struct {
	proposeCh    <-chan []byte
	confChangeCh <-chan raftpb.ConfChange
	commitCh     chan<- *commit
	errCh        chan<- error

	id    manta.ID
	peers []string

	conf Config

	// State
	confState     raftpb.ConfState
	snapshotIndex uint64
	appliedIndex  uint64

	// raft backing for the commit/error channel
	raftNode    raft.Node
	raftStorage *raft.MemoryStorage
	wal         *wal.WAL

	snapshotter      *snap.Snapshotter
	snapshotterReady chan *snap.Snapshotter // signals when snapshotter is ready

	snapCount uint64
	transport *rafthttp.Transport
	stopCh    chan struct{} // signals proposal channel closed

	logger *zap.Logger

	getSnapshot func() ([]byte, error)
}

func New(id manta.ID, conf Config, logger *zap.Logger) (*Node, error) {
	node := &Node{
		proposeCh:    make(chan []byte),
		confChangeCh: make(chan raftpb.ConfChange),
		commitCh:     make(chan *commit),
		errCh:        make(chan error),
		id:           id,

		conf: conf,

		stopCh: make(chan struct{}),
		logger: logger,
	}

	return node, nil
}

func (node *Node) startRaft(join bool) error {
	oldWal := wal.Exist(node.conf.WALDir)
	w, err := node.replayWAL()
	if err != nil {
		return err
	}

	node.wal = w

	// signal replay has finished
	node.snapshotterReady <- node.snapshotter

	peers := make([]raft.Peer, len(node.peers))
	for i, addr := range node.peers {
		peers[i] = raft.Peer{ID: xxhash.Sum64String(addr)}
	}

	c := &raft.Config{
		ID:                        uint64(node.id),
		ElectionTick:              10,
		HeartbeatTick:             1,
		Storage:                   node.raftStorage,
		MaxSizePerMsg:             1024 * 1024,
		MaxInflightMsgs:           256,
		MaxUncommittedEntriesSize: 1 << 30,
	}

	if oldWal || join {
		node.raftNode = raft.RestartNode(c)
	} else {
		node.raftNode = raft.StartNode(c, peers)
	}

	node.transport = &rafthttp.Transport{
		Logger:      node.logger,
		ID:          types.ID(node.id),
		ClusterID:   0x1000,
		Raft:        node,
		ServerStats: stats.NewServerStats("", ""),
		LeaderStats: stats.NewLeaderStats(node.logger, node.id.String()),
		ErrorC:      make(chan error),
	}

	if err = node.transport.Start(); err != nil {
		return err
	}

	for i := range node.peers {
		if node.peers[i] == node.conf.Listen {
			continue
		}
		addr := node.peers[i]
		node.transport.AddPeer(types.ID(xxhash.Sum64String(addr)), []string{addr})
	}

	go node.serveRaft()
	go node.serveChannels()

	return nil
}

func (node *Node) serveRaft() {
	snapshot, err := node.raftStorage.Snapshot()
	if err != nil {
		panic(err)
	}

	node.confState = snapshot.Metadata.ConfState
	node.snapshotIndex = snapshot.Metadata.Index
	node.appliedIndex = snapshot.Metadata.Index

	defer node.wal.Close()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// send proposals over raft
	go func() {
		confChangeCount := uint64(0)

		for node.proposeCh != nil && node.confChangeCh != nil {
			select {
			case data, ok := <-node.proposeCh:
				if !ok {
					node.proposeCh = nil
				} else {
					// blocks until accepted by raft state machine
					err = node.raftNode.Propose(context.Background(), data)
					if err != nil {
						node.logger.Warn("proposal failed",
							zap.Error(err))
					}
				}

			case cc, ok := <-node.confChangeCh:
				if !ok {
					node.confChangeCh = nil
				} else {
					confChangeCount += 1
					cc.ID = confChangeCount
					err = node.raftNode.ProposeConfChange(context.Background(), cc)
					if err != nil {
						node.logger.Warn("propose config change failed",
							zap.Error(err))
					}
				}
			}
		}

		// client closed channel; shutdown raft if not already
	}()

	// event loop on raft state machine updates
	for {
		select {
		case <-ticker.C:
			node.raftNode.Tick()

		// store raft entries to wal, then publish over commit channel
		case rd := <-node.raftNode.Ready():
			err = node.wal.Save(rd.HardState, rd.Entries)
			if err != nil {
				panic(err)
			}

			if !raft.IsEmptySnap(rd.Snapshot) {
				err = node.saveSnapshot(rd.Snapshot)
				if err != nil {
					panic(err)
				}

				err = node.raftStorage.ApplySnapshot(rd.Snapshot)
				if err != nil {
					panic(err)
				}

				node.publishSnapshot(rd.Snapshot)
			}

			err = node.raftStorage.Append(rd.Entries)
			if err != nil {
				panic(err)
			}

			node.transport.Send(rd.Messages)
			applyDoneCh, ok := node.publishEntries(node.entriesToApply(rd.CommittedEntries))
			if !ok {
				node.stop()
				return
			}

			node.maybeTriggerSnapshot(applyDoneCh)
			node.raftNode.Advance()

		case err := <-node.transport.ErrorC:
			panic(err)
			return

		case <-node.stopCh:
			node.stop()
			return
		}
	}
}

func (node *Node) entriesToApply(entries []raftpb.Entry) []raftpb.Entry {
	if len(entries) == 0 {
		return entries
	}

	firstIndex := entries[0].Index
	if firstIndex > node.appliedIndex+1 {
		node.logger.Fatal("first index of committed entry should less equal to appliedIndex",
			zap.Uint64("first-index", firstIndex),
			zap.Uint64("applied-index", node.appliedIndex+1))
	}

	if node.appliedIndex-firstIndex+1 < uint64(len(entries)) {
		return entries[node.appliedIndex-firstIndex+1:]
	}

	return nil
}

var snapshotCatchUpEntriesN uint64 = 100000

func (node *Node) maybeTriggerSnapshot(applyDoneC <-chan struct{}) {
	if node.appliedIndex-node.snapshotIndex <= node.snapCount {
		return
	}

	// wait until all committed entries are applied
	// or server is closed
	if applyDoneC != nil {
		select {
		case <-applyDoneC:
		case <-node.stopCh:
			return
		}
	}

	node.logger.Info("start snapshot",
		zap.Uint64("applied-index", node.appliedIndex),
		zap.Uint64("last-snapshot-index", node.snapshotIndex))

	data, err := node.getSnapshot()
	if err != nil {
		panic(err)
	}

	snapshot, err := node.raftStorage.CreateSnapshot(node.appliedIndex, &node.confState, data)
	if err != nil {
		panic(err)
	}

	if err := node.saveSnapshot(snapshot); err != nil {
		panic(err)
	}

	compactIndex := uint64(1)
	if node.appliedIndex > snapshotCatchUpEntriesN {
		compactIndex = node.appliedIndex - snapshotCatchUpEntriesN
	}

	if err := node.raftStorage.Compact(compactIndex); err != nil {
		panic(err)
	}

	node.logger.Info("compacted log",
		zap.Uint64("index", compactIndex))

	node.snapshotIndex = node.appliedIndex
}

func (node *Node) saveSnapshot(snapshot raftpb.Snapshot) error {
	walSnap := walpb.Snapshot{
		Index: snapshot.Metadata.Index,
		Term:  snapshot.Metadata.Term,
		// TODO: update to new version of etcd
		// ConfState: &snapshot.Metadata.ConfState,
	}

	// save the snapshot file before writing the snapshot to the wal.
	// This makes it possible fore the snapshot file to become
	// orphaned, but prevents a WAL snapshot entry from having no
	// corresponding snapshot file.
	if err := node.snapshotter.SaveSnap(snapshot); err != nil {
		return err
	}

	if err := node.wal.SaveSnapshot(walSnap); err != nil {
		return err
	}

	return node.wal.ReleaseLockTo(snapshot.Metadata.Index)
}

func (node *Node) serveChannels() {

}

func (node *Node) stop() {

}

// publishEntries writes committed log entries to commit channel and
// returns whether all entries could be published.
func (node *Node) publishEntries(entries []raftpb.Entry) (<-chan struct{}, bool) {
	if len(entries) == 0 {
		return nil, true
	}

	data := make([][]byte, 0, len(entries))
	for i := range entries {
		switch entries[i].Type {
		case raftpb.EntryNormal:
			if len(entries[i].Data) == 0 {
				// ignore empty messages
				break
			}

			data = append(data, entries[i].Data)

		case raftpb.EntryConfChange:
			var cc raftpb.ConfChange
			if err := cc.Unmarshal(entries[i].Data); err != nil {
				panic(err)
			}

			node.confState = *node.raftNode.ApplyConfChange(cc)

			switch cc.Type {
			case raftpb.ConfChangeAddNode:
				if len(cc.Context) > 0 {
					node.transport.AddPeer(types.ID(cc.NodeID), []string{string(cc.Context)})
				}

			case raftpb.ConfChangeRemoveNode:
				if cc.NodeID == uint64(node.id) {

					return nil, false
				}

				node.transport.RemovePeer(types.ID(cc.NodeID))
			}
		}
	}

	var applyDoneCh chan struct{}

	if len(data) > 0 {
		applyDoneCh = make(chan struct{}, 1)

		select {
		case node.commitCh <- &commit{data, applyDoneCh}:
		case <-node.stopCh:
			return nil, false
		}
	}

	// after commit, update appliedIndex
	node.appliedIndex = entries[len(entries)-1].Index

	return applyDoneCh, true
}

func (node *Node) publishSnapshot(snapshot raftpb.Snapshot) {
	if raft.IsEmptySnap(snapshot) {
		return
	}

	node.logger.Info("publishing snapshot",
		zap.Uint64("index", node.snapshotIndex))

	defer func(start time.Time) {
		node.logger.Info("publish snapshot finished",
			zap.Uint64("index", node.snapshotIndex),
			zap.Duration("elapsed", time.Since(start)))
	}(time.Now())

	if snapshot.Metadata.Index <= node.appliedIndex {
		node.logger.Fatal("snapshot index should great than progress applied index",
			zap.Uint64("snapshot-index", snapshot.Metadata.Index),
			zap.Uint64("applied-index", node.appliedIndex))
	}

	// trigger kvstore to load snapshot
	node.commitCh <- nil

	node.confState = snapshot.Metadata.ConfState
	node.snapshotIndex = snapshot.Metadata.Index
	node.appliedIndex = snapshot.Metadata.Index
}

func (node *Node) replayWAL() (*wal.WAL, error) {
	snapshot, err := node.loadSnapshot()
	if err != nil {
		return nil, err
	}

	w, err := node.openWAL(snapshot)
	if err != nil {
		return nil, err
	}

	_, st, ents, err := w.ReadAll()
	if err != nil {
		return nil, err
	}

	node.raftStorage = raft.NewMemoryStorage()
	if snapshot != nil {
		err = node.raftStorage.ApplySnapshot(*snapshot)
		if err != nil {
			return nil, err
		}
	}

	err = node.raftStorage.SetHardState(st)
	if err != nil {
		return nil, err
	}

	// append to storage so raft starts at the right place in log
	err = node.raftStorage.Append(ents)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (node *Node) loadSnapshot() (*raftpb.Snapshot, error) {
	if !wal.Exist(node.conf.WALDir) {
		return &raftpb.Snapshot{}, nil
	}

	walSnap, err := wal.ValidSnapshotEntries(node.logger, node.conf.WALDir)
	if err != nil {
		return nil, errors.Wrap(err, "valid snapshot entries failed")
	}

	snapshot, err := node.snapshotter.LoadNewestAvailable(walSnap)
	if err != nil && err != snap.ErrNoSnapshot {
		return nil, err
	}

	return snapshot, nil
}

func (node *Node) openWAL(snapshot *raftpb.Snapshot) (*wal.WAL, error) {
	if !wal.Exist(node.conf.WALDir) {
		if err := os.Mkdir(node.conf.WALDir, 0750); err != nil {
			return nil, err
		}

		w, err := wal.Create(node.logger, node.conf.WALDir, nil)
		if err != nil {
			return nil, err
		}

		if err := w.Close(); err != nil {
			return nil, err
		}
	}

	walSnap := walpb.Snapshot{}
	if snapshot != nil {
		walSnap.Index, walSnap.Term = snapshot.Metadata.Index, snapshot.Metadata.Term
	}

	node.logger.Info("loading WAL",
		zap.Uint64("term", walSnap.Term),
		zap.Uint64("index", walSnap.Index))

	w, err := wal.Open(node.logger, node.conf.WALDir, walSnap)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Process implement rafthttp.Raft
func (node *Node) Process(ctx context.Context, m raftpb.Message) error {
	return node.raftNode.Step(ctx, m)
}

func (node *Node) IsIDRemoved(id uint64) bool {
	return false
}

func (node *Node) ReportUnreachable(id uint64) {
	node.raftNode.ReportUnreachable(id)
}

func (node *Node) ReportSnapshot(id uint64, status raft.SnapshotStatus) {
	node.raftNode.ReportSnapshot(id, status)
}
