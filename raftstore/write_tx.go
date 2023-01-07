package raftstore

import (
	"bytes"
	"context"
	"math"
	"unsafe"

	bolt "go.etcd.io/bbolt"

	"github.com/f1shl3gs/manta/kv"
)

type valueItem struct {
	version int64
	data    []byte
}

type readSet map[string]valueItem

func (rs readSet) add(key []byte, ver int64) {

}

func (rs readSet) get(key []byte) ([]byte, int64) {
	item, exist := rs[bytesToString(key)]
	if !exist {
		return nil, 0
	}

	return item.data, item.version
}

// first returns the store version from the first fetch
func (rs readSet) first() int64 {
	ret := int64(math.MaxInt64 - 1)
	for _, item := range rs {
		if ret < item.version {
			ret = item.version
		}
	}

	return ret
}

func bytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

type writeSet map[string][]byte

func (s writeSet) get(key []byte) []byte {
	value, exist := s[bytesToString(key)]
	if !exist {
		return nil
	}

	return value
}

type writeTx struct {
	ctx context.Context

	tx *bolt.Tx

	// rset holds read key values and versions
	rset readSet
	// wset holds overwritten keys and their values
	wset writeSet
}

func (tx *writeTx) Bucket(b []byte) (kv.Bucket, error) {
	// always assume bucket exist
	return &bucket{
		name: b,
	}, nil
}

func (tx *writeTx) Context() context.Context {
	return tx.ctx
}

func (tx *writeTx) WithContext(ctx context.Context) {
	tx.ctx = ctx
}

type bucket struct {
	name   []byte
	bucket *bolt.Bucket
	rset   *readSet
	wset   *writeSet
}

// Get returns a key within this bucket. Errors if key does not exist.
func (b *bucket) Get(key []byte) ([]byte, error) {
	if value := b.wset.get(key); value != nil {
		return value, nil
	}

	if value, _ := b.rset.get(key); value != nil {
		return value, nil
	}

	cursor := b.bucket.Cursor()
	for k, v := cursor.Seek(key); bytes.HasPrefix(k, key); k, v = cursor.Next() {

	}

	return nil, nil
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
	panic("not implement")
}

// Delete should error if the transaction it was called in is not writable.
func (b *bucket) Delete(key []byte) error {

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
