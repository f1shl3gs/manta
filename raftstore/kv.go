package raftstore

import (
	"context"
	"io"

	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/raftstore/kvpb"
	"github.com/f1shl3gs/manta/raftstore/mvcc"
)

type writeTx struct {
	ctx context.Context
}

func (tx *writeTx) Bucket(b []byte) (kv.Bucket, error) {
	// TODO: check if bucket exist
	return &bucket{
		name:     b,
		writable: true,
	}, nil
}

func (tx *writeTx) Context() context.Context {
	return tx.ctx
}

func (tx *writeTx) WithContext(ctx context.Context) {
	tx.ctx = ctx
}

type Backend interface {
	// CreateBucket creates a bucket on the underlying store if it does not exist
	CreateBucket(ctx context.Context, bucket []byte) error
	// DeleteBucket deletes a bucket on the underlying store if it exists
	DeleteBucket(ctx context.Context, bucket []byte) error
}

type Store struct {
	backend Backend
	raft    *raftNode
	kv      mvcc.KV
}

// CreateBucket creates a bucket on the underlying store if it does not exist
func (s *Store) CreateBucket(ctx context.Context, bucket []byte) error {
	req := kvpb.Request{
		Id: s.raft.reqIDGen.Next(),
		Payload: &kvpb.Request_CreateBucket{
			CreateBucket: &kvpb.CreateBucket{
				Name: bucket,
			},
		},
	}

	data, err := req.Marshal()
	if err != nil {
		return err
	}

	return s.raft.Propose(ctx, data)
}

// DeleteBucket deletes a bucket on the underlying store if it exists
func (s *Store) DeleteBucket(ctx context.Context, bucket []byte) error {
	req := kvpb.Request{
		Id: s.raft.reqIDGen.Next(),
		Payload: &kvpb.Request_DeleteBucket{
			DeleteBucket: &kvpb.DeleteBucket{
				Name: bucket,
			},
		},
	}

	data, err := req.Marshal()
	if err != nil {
		return err
	}

	return s.raft.Propose(ctx, data)
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

type bucket struct {
	name []byte

	writable bool
}

// Get returns a key within this bucket. Errors if key does not exist.
func (b *bucket) Get(key []byte) ([]byte, error) {
	panic("not implement")
}

// GetBatch returns a corresponding set of values for the provided
// set of keys. If a value cannot be found for any provided key its
// value will be nil at the same index for the provided key.
func (b *bucket) GetBatch(keys ...[]byte) ([][]byte, error) {
	panic("not implement")
}

// Cursor returns a cursor at the beginning of this bucket optionally
// using the provided hints to improve performance.
func (b *bucket) Cursor(hints ...kv.CursorHint) (kv.Cursor, error) {
	panic("not implement")
}

// Put should error if the transaction it was called in is not writable.
func (b *bucket) Put(key, value []byte) error {
	if !b.writable {
		return kv.ErrTxNotWritable
	}

	panic("not implement")
}

// Delete should error if the transaction it was called in is not writable.
func (b *bucket) Delete(key []byte) error {
	if !b.writable {
		return kv.ErrTxNotWritable
	}

	panic("not implement")
}

// ForwardCursor returns a forward cursor from the seek position provided.
// Other options can be supplied to provide direction and hints.
func (b *bucket) ForwardCursor(seek []byte, opts ...kv.CursorOption) (kv.ForwardCursor, error) {
	panic("not implement")
}

type forwardCursor struct {
}

// Next moves the cursor to the next key in the bucket.
func (c *forwardCursor) Next() (k, v []byte) {
	panic("not implement")
}

// Err returns non-nil if an error occurred during cursor iteration.
// This should always be checked after Next returns a nil key/value.
func (c *forwardCursor) Err() error {
	panic("not implement")
}

// Close is responsible for freeing any resources created by the cursor.
func (c *forwardCursor) Close() error {
	panic("not implement")
}
