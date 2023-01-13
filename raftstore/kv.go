package raftstore

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"

	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/raftstore/kvpb"

	bolt "go.etcd.io/bbolt"
	"go.etcd.io/raft/v3"
)

type KV struct {
	logger *zap.Logger
	db     *bolt.DB
	raft   raft.Node
	idGen  *idGenerator
	wait   Wait

	linearizableReadNotify func(ctx context.Context) error
}

// CreateBucket creates a bucket on the underlying store if it does not exist
func (s *KV) CreateBucket(ctx context.Context, bucket []byte) error {
	id := s.idGen.Next()
	waitCh := s.wait.Register(id)

	req := kvpb.InternalRequest{
		ID: id,
		Request: &kvpb.InternalRequest_CreateBucket{
			CreateBucket: &kvpb.CreateBucket{
				Name: bucket,
			},
		},
	}

	data, err := req.Marshal()
	if err != nil {
		return err
	}

	err = s.raft.Propose(ctx, data)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-waitCh:
		return nil
	}
}

// DeleteBucket deletes a bucket on the underlying store if it exists
func (s *KV) DeleteBucket(ctx context.Context, bucket []byte) error {
	id := s.idGen.Next()
	waitCh := s.wait.Register(id)

	req := kvpb.InternalRequest{
		ID: id,
		Request: &kvpb.InternalRequest_DeleteBucket{
			DeleteBucket: &kvpb.DeleteBucket{
				Name: bucket,
			},
		},
	}

	data, err := req.Marshal()
	if err != nil {
		return err
	}

	if err = s.raft.Propose(ctx, data); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-waitCh:
		return nil
	}
}

// View opens up a transaction that will not write to any value. Implementing interfaces
// should take care to ensure that all view transactions do not mutate any value.
func (s *KV) View(ctx context.Context, fn func(kv.Tx) error) error {
	err := s.linearizableReadNotify(ctx)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin(false)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			panic(fmt.Sprintf("read tx fallback failed, %s", err))
		}
	}()

	return fn(&readTx{
		ctx: ctx,
		Tx:  tx,
	})
}

// Update opens up a transaction that will mutate value.
func (s *KV) Update(ctx context.Context, fn func(kv.Tx) error) error {
	// TODO: some transication do not need to read anything, so we can move
	// this when read happends
	err := s.linearizableReadNotify(ctx)
	if err != nil {
		return err
	}

	// write operation is cached and it will be propose throught raft,
	// the it will be applied, so we don't need this tx to be writable.
	tx, err := s.db.Begin(false)
	if err != nil {
		return err
	}

	wtx := &writeTx{
		ctx:  ctx,
		tx:   tx,
		rset: make(map[string]readSet),
		wset: make(map[string]writeSet),
	}
	err = fn(wtx)
	if rErr := tx.Rollback(); rErr != nil {
		panic(fmt.Sprintf("read tx fallback failed, %s", rErr))
	}
	if err != nil {
		return err
	}

	txn := wtx.txn()
	if len(txn.Successes) == 0 && len(txn.Failures) == 0 {
		// This txn is readonly, so is should not be proposal to raft
		return nil
	}

	id := s.idGen.Next()
	waitCh := s.wait.Register(id)
	req := &kvpb.InternalRequest{
		ID: id,
		Request: &kvpb.InternalRequest_Txn{
			Txn: txn,
		},
	}

	data, err := req.Marshal()
	if err != nil {
		return err
	}

	if err = s.raft.Propose(ctx, data); err != nil {
		return err
	}

	// TODO: retry !?

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-waitCh:
		return nil
	}

	return nil
}

// Backup copies all K:Vs to a writer, file format determined by implementation.
func (s *KV) Backup(ctx context.Context, w io.Writer) error {
	panic("not implement")
}
