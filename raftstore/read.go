package raftstore

import (
    "context"

    "github.com/f1shl3gs/manta/kv"
)

type readTx struct {
    ctx context.Context
}

// Bucket possibly creates and returns bucket, b.
func (tx *readTx) Bucket(b []byte) (kv.Bucket, error) {
    return &readOnlyBucket{
        readTx: tx,
        name: b,
    }, nil
}

// Context returns the context associated with this Tx.
func (tx *readTx) Context() context.Context {
    return tx.ctx
}

// WithContext associates a context with this Tx.
func (tx *readTx) WithContext(ctx context.Context) {
    tx.ctx = ctx
}

type readOnlyBucket struct {
    *readTx
    name []byte
}

// Get returns a key within this bucket. Errors if key does not exist.
func (b *readOnlyBucket) Get(key []byte) ([]byte, error) {

}

// GetBatch returns a corresponding set of values for the provided
// set of keys. If a value cannot be found for any provided key its
// value will be nil at the same index for the provided key.
func (b *readOnlyBucket) GetBatch(keys ...[]byte) ([][]byte, error) {

}

// Cursor returns a cursor at the beginning of this bucket optionally
// using the provided hints to improve performance.
func (b *readOnlyBucket) Cursor(hints ...kv.CursorHint) (kv.Cursor, error) {

}

// Put should error if the transaction it was called in is not writable.
func (b *readOnlyBucket) Put(key, value []byte) error {
    return kv.ErrTxNotWritable
}

// Delete should error if the transaction it was called in is not writable.
func (b *readOnlyBucket) Delete(key []byte) error {
    return kv.ErrTxNotWritable
}

// ForwardCursor returns a forward cursor from the seek position provided.
// Other options can be supplied to provide direction and hints.
func (b *readOnlyBucket) ForwardCursor(seek []byte, opts ...kv.CursorOption) (kv.ForwardCursor, error) {
    panic("not implement")
}
