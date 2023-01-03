package raftstore

import (
	"context"
	"go.etcd.io/bbolt"

	"github.com/f1shl3gs/manta/kv"
)

type readTx struct {
	ctx   context.Context
	inner *bbolt.Tx
}

func (tx *readTx) Bucket(name []byte) (kv.Bucket, error) {
	b := tx.inner.Bucket(name)
	if b == nil {
		return nil, kv.ErrBucketNotFound
	}

	return &readBucket{
		Bucket: b,
	}, nil
}

func (tx *readTx) Context() context.Context {
	return tx.ctx
}

func (tx *readTx) WithContext(ctx context.Context) {
	tx.ctx = ctx
}

type readBucket struct {
	*bbolt.Bucket
}

// Get returns a key within this bucket. Errors if key does not exist.
func (b *readBucket) Get(key []byte) ([]byte, error) {
	value := b.Bucket.Get(key)
	return value, nil
}

// GetBatch returns a corresponding set of values for the provided
// set of keys. If a value cannot be found for any provided key its
// value will be nil at the same index for the provided key.
func (b *readBucket) GetBatch(keys ...[]byte) ([][]byte, error) {
	values := make([][]byte, len(keys))
	for _, key := range keys {
		values = append(values, b.Bucket.Get(key))
	}

	return values, nil
}

// Cursor returns a cursor at the beginning of this bucket optionally
// using the provided hints to improve performance.
func (b *readBucket) Cursor(hints ...kv.CursorHint) (kv.Cursor, error) {
	panic("not implement")
}

// Put should error if the transaction it was called in is not writable.
func (b *readBucket) Put(key, value []byte) error {
	return kv.ErrTxNotWritable
}

// Delete should error if the transaction it was called in is not writable.
func (b *readBucket) Delete(key []byte) error {
	return kv.ErrTxNotWritable
}

// ForwardCursor returns a forward cursor from the seek position provided.
// Other options can be supplied to provide direction and hints.
func (b *readBucket) ForwardCursor(seek []byte, opts ...kv.CursorOption) (kv.ForwardCursor, error) {
	panic("not implement")
}
