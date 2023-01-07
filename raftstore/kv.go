package raftstore

import (
	"context"
	"io"

	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/raftstore/raftpb"
)

type Store struct {
	raft *raftNode
}

// CreateBucket creates a bucket on the underlying store if it does not exist
func (s *Store) CreateBucket(ctx context.Context, bucket []byte) error {
	return s.raft.Propose(ctx, &raftpb.CreateBucket{
		Name: bucket,
	})
}

// DeleteBucket deletes a bucket on the underlying store if it exists
func (s *Store) DeleteBucket(ctx context.Context, bucket []byte) error {
	return s.raft.Propose(ctx, &raftpb.DeleteBucket{
		Name: bucket,
	})
}

// View opens up a transaction that will not write to any data. Implementing interfaces
// should take care to ensure that all view transactions do not mutate any data.
func (s *Store) View(ctx context.Context, fn func(kv.Tx) error) error {
	err := s.raft.linearizableReadNotify(ctx)
	if err != nil {
		return err
	}

	return fn(&readTx{ctx})
}

// Update opens up a transaction that will mutate data.
func (s *Store) Update(ctx context.Context, fn func(kv.Tx) error) error {
	return fn(&writeTx{ctx})
}

// Backup copies all K:Vs to a writer, file format determined by implementation.
func (s *Store) Backup(ctx context.Context, w io.Writer) error {
	panic("not implement")
}
