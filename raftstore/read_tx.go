package raftstore

import (
	"bytes"
	"context"
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/f1shl3gs/manta/kv"
)

type readTx struct {
	*bolt.Tx

	ctx context.Context
}

// Bucket possibly creates and returns bucket, b.
func (tx *readTx) Bucket(name []byte) (kv.Bucket, error) {
	b := tx.Tx.Bucket(name)
	if b == nil {
		return nil, kv.ErrBucketNotFound
	}

	return &readOnlyBucket{
		bucket: b,
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
	bucket *bolt.Bucket
}

// Get returns a key within this bucket. Errors if key does not exist.
func (b *readOnlyBucket) Get(key []byte) ([]byte, error) {
	return b.bucket.Get(key), nil
}

// GetBatch returns a corresponding set of values for the provided
// set of keys. If a value cannot be found for any provided key its
// value will be nil at the same index for the provided key.
func (b *readOnlyBucket) GetBatch(keys ...[]byte) ([][]byte, error) {
	values := make([][]byte, len(keys))
	for idx, key := range keys {
		values[idx] = b.bucket.Get(key)
	}

	return values, nil
}

// cursor is a struct for iterating through the entries
// in the key value store.
type cursor struct {
	cursor *bolt.Cursor
	// previously seeked key/value
	key, value []byte

	config kv.CursorConfig
	closed bool
	seen   int
}

// Seek moves the cursor forward until reaching prefix in the key name.
func (c *cursor) Seek(prefix []byte) ([]byte, []byte) {
	if c.closed {
		return nil, nil
	}

	k, v := c.cursor.Seek(prefix)
	if len(k) == 0 && len(v) == 0 {
		return nil, nil
	}

	return k, v
}

// First moves the cursor to the first key in the bucket.
func (c *cursor) First() ([]byte, []byte) {
	if c.closed {
		return nil, nil
	}

	k, v := c.cursor.First()
	if len(k) == 0 && len(v) == 0 {
		return nil, nil
	}

	return k, v
}

// Last moves the cursor to the last key in the bucket.
func (c *cursor) Last() ([]byte, []byte) {
	if c.closed {
		return nil, nil
	}

	k, v := c.cursor.Last()
	if len(k) == 0 && len(v) == 0 {
		return nil, nil
	}

	return k, v
}

// Next moves the cursor to the next key in the bucket.
func (c *cursor) Next() (k []byte, v []byte) {
	if c.closed || c.atLimit() || (c.key != nil && c.missingPrefix(c.key)) {
		return nil, nil
	}

	// get and unset previously seeked values if they exist
	k, v, c.key, c.value = c.key, c.value, nil, nil
	if len(k) > 0 || len(v) > 0 {
		c.seen += 1
		return
	}

	if c.config.Direction == kv.CursorDescending {
		k, v = c.cursor.Prev()
	} else {
		k, v = c.cursor.Next()
	}

	if (len(k) == 0 && len(v) == 0) || c.missingPrefix(k) {
		return nil, nil
	}

	c.seen += 1

	return k, v
}

func (c *cursor) missingPrefix(key []byte) bool {
	return c.config.Prefix != nil && !bytes.HasPrefix(key, c.config.Prefix)
}

func (c *cursor) atLimit() bool {
	return c.config.Limit != nil && c.seen >= *c.config.Limit
}

// Prev moves the cursor to the prev key in the bucket.
func (c *cursor) Prev() (k []byte, v []byte) {
	if c.closed || c.atLimit() || (c.key != nil && c.missingPrefix(c.key)) {
		return nil, nil
	}

	// get and unset previously seeked values if they exist
	k, v, c.key, c.value = c.key, c.value, nil, nil
	if len(k) > 0 && len(v) > 0 {
		c.seen++
		return
	}

	if c.config.Direction == kv.CursorDescending {
		k, v = c.cursor.Next()
	} else {
		k, v = c.cursor.Prev()
	}

	if (len(k) == 0 && len(v) == 0) || c.missingPrefix(k) {
		return nil, nil
	}

	c.seen++

	return k, v
}

// Cursor returns a cursor at the beginning of this bucket optionally
// using the provided hints to improve performance.
func (b *readOnlyBucket) Cursor(hints ...kv.CursorHint) (kv.Cursor, error) {
	return &cursor{
		cursor: b.bucket.Cursor(),
	}, nil
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
	var (
		c          = b.bucket.Cursor()
		config     = kv.NewCursorConfig(opts...)
		key, value []byte
	)

	if len(seek) == 0 && config.Direction == kv.CursorDescending {
		seek, _ = c.Last()
	}

	key, value = c.Seek(seek)

	if config.Prefix != nil && !bytes.HasPrefix(seek, config.Prefix) {
		return nil, fmt.Errorf("seek bytes %q not prefixed with %q: %w", string(seek), string(config.Prefix), kv.ErrSeekMissingPrefix)
	}

	fc := &cursor{
		cursor: c,
		config: config,
	}

	// only remember first seeked item if not skipped
	if !config.SkipFirst {
		fc.key = key
		fc.value = value
	}

	return fc, nil
}

// Err always returns nil as nothing can go wrong™ during iteration
func (c *cursor) Err() error {
	return nil
}

// Close sets the closed to closed
func (c *cursor) Close() error {
	c.closed = true
	return nil
}
