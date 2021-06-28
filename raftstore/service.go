package raftstore

import (
	"context"

	"github.com/f1shl3gs/manta/raftstore/internal"
	"github.com/f1shl3gs/manta/raftstore/rawkv"
	"github.com/f1shl3gs/manta/raftstore/transport"
	bolt "go.etcd.io/bbolt"

	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.uber.org/zap"
)

// Send implement transport.RaftServer
func (s *Store) Send(ctx context.Context, m *raftpb.Message) (*transport.Done, error) {
	if s.cluster.IDRemoved(m.From) {
		s.logger.Warn("rejected Raft message from removed member",
			zap.String("local-member-id", internal.IDToString(s.id)),
			zap.String("from", internal.IDToString(m.From)))
		return nil, ErrRejectFromRemovedMember
	}

	return nil, s.raftNode.Step(ctx, *m)
}

// SendSnapshot implement transport.RaftServer
func (s *Store) SendSnapshot(server transport.Raft_SendSnapshotServer) error {
	// TODO: there are something need pay attention
	// 1. do we need to save it to disk, and trigger a GC?
	// 2. restore the kv engine right here?
	panic("implement me")
}

// Put implement RawKV
func (s *Store) Put(ctx context.Context, req *rawkv.PutRequest) (*rawkv.Empty, error) {
	req.Id = s.reqIDGen.Next()
	waitCh := s.wait.Register(req.Id)

	data, _ := req.Marshal()

	err := s.raftNode.Propose(ctx, data)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-waitCh:
	}

	return &rawkv.Empty{}, nil
}

// Get implement RawKV
func (s *Store) Get(ctx context.Context, req *rawkv.GetRequest) (*rawkv.GetResponse, error) {
	err := s.linearizableReadNotify(ctx)
	if err != nil {
		return nil, err
	}

	var val []byte

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("kv"))
		val = b.Get(req.Key)
		return nil
	})

	return &rawkv.GetResponse{Value: val}, nil
}

func (s *Store) linearizableReadNotify(ctx context.Context) error {
	s.readMtx.RLock()
	nc := s.readNotifier
	s.readMtx.RUnlock()

	// signal linearizable loop for current notify if it hasn't been already
	select {
	case s.readWaitCh <- struct{}{}:
	default:
	}

	// wait for read state notification
	select {
	case <-nc.c:
		return nc.err
	case <-ctx.Done():
		return ctx.Err()
	case <-s.stopCh:
		return ErrStopped
	}
}
