package raftstore

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"sort"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/f1shl3gs/manta/raftstore/membership"
	"github.com/f1shl3gs/manta/raftstore/transport"

	"github.com/gogo/protobuf/proto"
	"go.etcd.io/etcd/pkg/v3/pbutil"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/wal"
	"go.etcd.io/etcd/server/v3/wal/walpb"
	"go.uber.org/zap"
)

func startNode(cf *Config, logger *zap.Logger, cl membership.Cluster) (uint64, raft.Node, *raft.MemoryStorage, *wal.WAL) {
	var (
		w   *wal.WAL
		err error
		n   raft.Node
		id  uint64
	)

	members := cl.Members()
	for _, m := range members {
		if m.Address == cf.BindAddr {
			id = m.ID
		}
	}

	if id == 0 {
		panic("cannot verify local node id")
	}

	metadata := make([]byte, 16)
	binary.BigEndian.PutUint64(metadata, cl.ID())
	binary.BigEndian.PutUint64(metadata[8:], id)

	if w, err = wal.Create(logger, cf.WALDir, metadata); err != nil {
		logger.Panic("failed to create WAL",
			zap.String("path", cf.WALDir),
			zap.Error(err))
	}

	peers := make([]raft.Peer, len(members))
	for i, m := range members {
		data, _ := json.Marshal(m)
		peers[i] = raft.Peer{
			ID:      m.ID,
			Context: data,
		}
	}

	logger.Info("starting local member",
		zap.String("id", strconv.FormatUint(id, 16)),
		zap.String("cluster", strconv.FormatUint(cl.ID(), 16)),
	)

	s := raft.NewMemoryStorage()
	c := &raft.Config{
		ID:              id,
		ElectionTick:    cf.ElectionTicks / int(cf.TickMs),
		HeartbeatTick:   1,
		Storage:         s,
		MaxSizePerMsg:   maxSizePerMsg,
		MaxInflightMsgs: maxInflightMsgs,
		CheckQuorum:     true,
		PreVote:         cf.PreVote,
		Logger:          NewRaftLoggerZap(logger),
	}

	if len(peers) == 0 {
		n = raft.RestartNode(c)
	} else {
		n = raft.StartNode(c, peers)
	}

	return id, n, s, w
}

// getIDs returns an ordered set of IDs included in the given snapshot
// and the entries. The given snapshot/entries can contain three kinds
// of ID-related entry:
// - ConfChangeAddNode, in which case the contained ID will be added into the set
// - ConfChangeRemoveNode, in which case the contained ID will be removed from the set
// - ConfChangeAddLearnerNode, in which case the contained ID will be added into the set
func getIDs(logger *zap.Logger, snap *raftpb.Snapshot, ents []raftpb.Entry) []uint64 {
	ids := make(map[uint64]bool)
	if snap != nil {
		for _, id := range snap.Metadata.ConfState.Voters {
			ids[id] = true
		}
	}

	for _, e := range ents {
		if e.Type != raftpb.EntryConfChange {
			continue
		}

		var cc raftpb.ConfChange
		pbutil.MustUnmarshal(&cc, e.Data)
		switch cc.Type {
		case raftpb.ConfChangeAddLearnerNode, raftpb.ConfChangeAddNode:
			ids[cc.NodeID] = true

		case raftpb.ConfChangeRemoveNode:
			delete(ids, cc.NodeID)

		case raftpb.ConfChangeUpdateNode:
			// do nothing
		default:
			logger.Panic("unknown ConfChange Type",
				zap.String("type", cc.Type.String()))
		}
	}

	sids := make([]uint64, 0, len(ids))
	for id := range ids {
		sids = append(sids, id)
	}

	sort.Slice(sids, func(i, j int) bool {
		return sids[i] > sids[j]
	})

	return sids
}

// createConfigChangeEnts creates a series of Raft entries(ie EntryConfChange)
// to remove the set of given IDs from the cluster. The ID `self` is not removed, even
// if present in the set. If `self` is not inside the given ids, it creates a Raft
// entry to add a default member with the given `self`
func createConfigChangeEnts(logger *zap.Logger, ids []uint64, self uint64, term, index uint64) []raftpb.Entry {
	found := false
	for _, id := range ids {
		if id == self {
			found = true
		}
	}

	var ents []raftpb.Entry
	next := index + 1

	// NB: always add self first, then remove other nodes, Raft will panic
	// if the set of voters ever become empty
	if !found {
		m := &membership.Member{
			ID:      self,
			Address: "localhost:8087",
		}

		data, _ := json.Marshal(m)
		cc := &raftpb.ConfChange{
			Type:   raftpb.ConfChangeAddNode,
			NodeID: self,
			// todo: why we need Context
			Context: data,
		}
		e := raftpb.Entry{
			Type:  raftpb.EntryConfChange,
			Data:  pbutil.MustMarshal(cc),
			Term:  term,
			Index: next,
		}
		ents = append(ents, e)
		next++
	}

	for _, id := range ids {
		if id == self {
			continue
		}

		cc := &raftpb.ConfChange{
			Type:   raftpb.ConfChangeRemoveNode,
			NodeID: id,
		}

		e := raftpb.Entry{
			Type:  raftpb.EntryConfChange,
			Data:  pbutil.MustMarshal(cc),
			Term:  term,
			Index: next,
		}

		ents = append(ents, e)
		next++
	}

	return ents
}

func restartAsStandaloneNode(cf *Config, logger *zap.Logger, snapshot *raftpb.Snapshot) (uint64, membership.Cluster, raft.Node, *raft.MemoryStorage, *wal.WAL) {
	var (
		walsnap walpb.Snapshot
	)

	if snapshot != nil {
		walsnap.Index, walsnap.Term = snapshot.Metadata.Index, snapshot.Metadata.Term
	}

	w, id, cid, st, ents := readWAL(logger, cf.WALDir, walsnap)

	// discard the previously uncommitted entries
	for i, ent := range ents {
		if ent.Index > st.Commit {
			logger.Info("discarding uncommitted WAL entries",
				zap.Uint64("entry-index", ent.Index),
				zap.Uint64("committed", st.Commit),
				zap.Int("number", len(ents)-i))

			ents = ents[:i]
			break
		}
	}

	// force append the configuration change entries
	toAppEnts := createConfigChangeEnts(
		logger,
		getIDs(logger, snapshot, ents),
		id,
		st.Term,
		st.Commit,
	)
	ents = append(ents, toAppEnts...)

	// force commit newly appended entries
	err := w.Save(raftpb.HardState{}, toAppEnts)
	if err != nil {
		logger.Fatal("failed to save hard state and entries",
			zap.Error(err))
	}
	if len(ents) != 0 {
		st.Commit = ents[len(ents)-1].Index
	}

	logger.Info("forcing restart member",
		zap.String("cluster", strconv.FormatUint(cid, 16)),
		zap.String("local-member-id", strconv.FormatUint(id, 16)),
		zap.Uint64("commit", st.Commit))

	cl := membership.NewCluster()
	cl.SetID(cid)
	s := raft.NewMemoryStorage()
	if snapshot != nil {
		s.ApplySnapshot(*snapshot)
	}
	s.SetHardState(st)
	s.Append(ents)
	c := &raft.Config{
		ID:              id,
		ElectionTick:    cf.ElectionTicks,
		HeartbeatTick:   1,
		Storage:         s,
		MaxSizePerMsg:   maxSizePerMsg,
		MaxInflightMsgs: maxInflightMsgs,
		CheckQuorum:     true,
		PreVote:         cf.PreVote,
		Logger:          NewRaftLoggerZap(logger),
	}

	n := raft.RestartNode(c)
	return id, cl, n, s, w
}

func restartNode(cf *Config, logger *zap.Logger, snapshot *raftpb.Snapshot) (uint64, membership.Cluster, raft.Node, *raft.MemoryStorage, *wal.WAL) {
	var walsnap walpb.Snapshot

	if snapshot != nil {
		walsnap.Index, walsnap.Term = snapshot.Metadata.Index, snapshot.Metadata.Term
	}

	w, id, cid, st, ents := readWAL(logger, cf.WALDir, walsnap)

	logger.Info("restarting local member",
		zap.String("cluster", strconv.FormatUint(cid, 16)),
		zap.String("id", strconv.FormatUint(id, 16)),
		zap.Uint64("commit-index", st.Commit),
	)

	cl := membership.NewCluster()
	cl.SetID(cid)
	s := raft.NewMemoryStorage()
	if snapshot != nil {
		s.ApplySnapshot(*snapshot)
	}

	s.SetHardState(st)
	s.Append(ents)

	c := &raft.Config{
		ID:              id,
		ElectionTick:    cf.ElectionTicks,
		HeartbeatTick:   1,
		Storage:         s,
		MaxSizePerMsg:   maxSizePerMsg,
		MaxInflightMsgs: maxInflightMsgs,
		CheckQuorum:     true,
		PreVote:         cf.PreVote,
		Logger:          NewRaftLoggerZap(logger),
	}

	n := raft.RestartNode(c)
	return id, cl, n, s, w
}

// ReadyHandler contains a set of EtcdServer operations to be called by raftNode,
// and helps decouple state machine logic from Raft algorithms.
// TODO: add a state machine interface to apply the commit entries and do snapshot/recover
type ReadyHandler struct {
	getLead              func() (lead uint64)
	updateLead           func(lead uint64)
	updateLeadership     func(newLeader bool)
	updateCommittedIndex func(uint64)
}

/*
type ReadyHandler interface {
	Lead() uint64

	SetLead(id uint64)

	UpdateLeadership(newLeader bool)

	UpdateCommittedIndex(index uint64)
}
*/
// apply contains entries, snapshot to be applied. Once
// an apply is consumed, the entries will be persisted to
// to raft storage concurrently; the application must read
// raftDone before assuming the raft messages are stable.
type apply struct {
	entries  []raftpb.Entry
	snapshot raftpb.Snapshot
	// notifyc synchronizes etcd server applies with the raft node
	notifyc chan struct{}
}

type applyResult struct {
	resp proto.Message
	err  error
	// physCh signals the physical effect of the request has completed in addition
	// to being logically reflected by the node. Currently only used for
	// Compaction requests.
	physCh <-chan struct{}
}

func (s *Store) raftRequest(ctx context.Context, payload []byte) (proto.Message, error) {
	result, err := s.processInternalRaftRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	if result.err != nil {
		return nil, result.err
	}

	return result.resp, nil
}

func (s *Store) processInternalRaftRequest(ctx context.Context, data []byte) (*applyResult, error) {
	appliedIndex := s.getAppliedIndex()
	committedIndex := s.getCommittedIndex()
	if committedIndex > appliedIndex+maxGapBetweenApplyAndCommitIndex {
		return nil, ErrTooManyRequests
	}

	if len(data) > int(s.config.MaxRequestBytes) {
		return nil, ErrRequestTooLarge
	}

	// TODO: set request id to data

	id := s.reqIDGen.Next()

	ch := s.wait.Register(id)

	ctx, cancel := context.WithTimeout(ctx, s.config.RequestTimeout())
	defer cancel()

	start := time.Now()
	err := s.raftNode.Propose(ctx, data)
	if err != nil {
		s.proposalsFailed.Inc()
		s.wait.Trigger(id, nil) // GC wait
		return nil, err
	}

	s.proposalsPending.Inc()
	defer s.proposalsPending.Dec()

	select {
	case x := <-ch:
		return x.(*applyResult), err
	case <-ctx.Done():
		s.proposalsFailed.Inc()
		s.wait.Trigger(id, nil) // GC wait

		return nil, s.parseProposeCtxErr(ctx.Err(), start)
	case <-s.stopCh:
		return nil, ErrStopped
	}
}

func (s *Store) parseProposeCtxErr(err error, start time.Time) error {
	switch err {
	case context.Canceled:
		return ErrCanceled

	case context.DeadlineExceeded:
		s.leadTimeMtx.RLock()
		curLeadElected := s.leadElectedTime
		s.leadTimeMtx.RUnlock()
		prevLeadLost := curLeadElected.Add(-2 * time.Duration(s.config.ElectionTicks) *
			time.Duration(s.config.TickMs) * time.Millisecond)
		if start.After(prevLeadLost) && start.Before(curLeadElected) {
			return ErrTimeoutDueToLeaderFail
		}

		lead := s.getLead()
		switch lead {
		case raft.None:
			// TODO: return error to specify it happens because the cluster does not have leader now
		case s.id:
			if !isConnectedToQuorumSince(s.transporter, start, s.id, s.cluster.Members()) {
				return ErrTimeoutDueToConnectionLost
			}
		default:
			if !isConnectedSince(s.transporter, start, lead) {
				return ErrTimeoutDueToConnectionLost
			}
		}
		return ErrTimeout

	default:
		return err
	}
}

func (s *Store) getLead() uint64 {
	return atomic.LoadUint64(&s.lead)
}

func (s *Store) setLead(lead uint64) {
	atomic.StoreUint64(&s.lead, lead)
}

// isConnectedToQuorumSince checks whether the local member is connected to the
// quorum of the cluster since the given time.
func isConnectedToQuorumSince(transport transport.Transporter, since time.Time, self uint64, members []*membership.Member) bool {
	return numConnectedSince(transport, since, self, members) >= (len(members)/2)+1
}

// isConnectedSince checks whether the local member is connected to the
// remote member since the given time.
func isConnectedSince(transport transport.Transporter, since time.Time, remote uint64) bool {
	t := transport.ActiveSince(remote)
	return !t.IsZero() && t.Before(since)
}

// numConnectedSince counts how many members are connected to the local member
// since the given time.
func numConnectedSince(transport transport.Transporter, since time.Time, self uint64, members []*membership.Member) int {
	connectedNum := 0
	for _, m := range members {
		if m.ID == self || isConnectedSince(transport, since, m.ID) {
			connectedNum++
		}
	}
	return connectedNum
}
