package raftstore

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"reflect"
	"unsafe"

	bolt "go.etcd.io/bbolt"

	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/raftstore/kvpb"
)

type valueItem struct {
	version int64
	value   []byte
}

type readSet map[string]valueItem

func (rs readSet) add(key, value []byte, version int64) {
	rs[unsafeBytesToString(key)] = valueItem{
		version: version,
		value:   value,
	}
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

func unsafeBytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func unsafeStringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

type writeOp struct {
	value    []byte
	deletion bool
}

type writeSet map[string]writeOp

func (s writeSet) get(key []byte) []byte {
	op, exist := s[unsafeBytesToString(key)]
	if !exist {
		return nil
	}

	if op.deletion {
		return nil
	}

	return op.value
}

type writeTx struct {
	ctx context.Context

	tx *bolt.Tx

	// rset holds read key values and versions
	rset map[string]readSet
	// wset holds overwritten keys and their values
	wset map[string]writeSet
}

func (tx *writeTx) txn() *kvpb.Txn {
	txn := &kvpb.Txn{}

	for b, rset := range tx.rset {
		for key, item := range rset {
			txn.Compares = append(txn.Compares, &kvpb.Compare{
				Bucket:  unsafeStringToBytes(b),
				Key:     unsafeStringToBytes(key),
				Version: item.version,
			})
		}
	}

	for b, wset := range tx.wset {
		for key, item := range wset {
			op := &kvpb.Operation{
				Bucket: unsafeStringToBytes(b),
				Key:    unsafeStringToBytes(key),
			}

			if item.deletion {
				op.Deletion = true
			} else {
				op.Value = item.value
			}

			txn.Successes = append(txn.Successes, op)
		}
	}

	return txn
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
	name []byte

	// readOnly bucket
	bucket *bolt.Bucket
	rset   readSet
	wset   writeSet
}

// Get returns a key within this bucket. Errors if key does not exist.
func (b *bucket) Get(key []byte) ([]byte, error) {
	sk := unsafeBytesToString(key)
	if value := b.wset.get(key); value != nil {
		return value, nil
	}

	if item, exist := b.rset[sk]; exist {
		return item.value, nil
	}

	value := b.bucket.Get(key)
	b.rset.add(key, value, 0)

	return value, nil
}

// GetBatch returns a corresponding set of values for the provided
// set of keys. If a value cannot be found for any provided key its
// value will be nil at the same index for the provided key.
func (b *bucket) GetBatch(keys ...[]byte) ([][]byte, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	values := make([][]byte, len(keys))

	for _, key := range keys {
		sk := unsafeBytesToString(key)
		op, exist := b.wset[sk]
		if exist {
			if op.deletion {
				values = append(values, nil)
			} else {
				values = append(values, op.value)
			}

			continue
		}

		item, exist := b.rset[sk]
		if exist {
			values = append(values, item.value)
			continue
		}

		value := b.bucket.Get(key)
		b.rset.add(key, value, 0)
	}

	return values, nil
}

// Cursor returns a cursor at the beginning of this bucket optionally
// using the provided hints to improve performance.
func (b *bucket) Cursor(hints ...kv.CursorHint) (kv.Cursor, error) {
	return &cursor{
		cursor: b.bucket.Cursor(),
	}, nil
}

// Put should error if the transaction it was called in is not writable.
func (b *bucket) Put(key, value []byte) error {
	b.wset[unsafeBytesToString(key)] = writeOp{
		value:    value,
		deletion: false,
	}

	return nil
}

// Delete should error if the transaction it was called in is not writable.
func (b *bucket) Delete(key []byte) error {
	b.wset[unsafeBytesToString(key)] = writeOp{
		value:    nil,
		deletion: true,
	}
	return nil
}

// ForwardCursor returns a forward cursor from the seek position provided.
// Other options can be supplied to provide direction and hints.
func (b *bucket) ForwardCursor(seek []byte, opts ...kv.CursorOption) (kv.ForwardCursor, error) {
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
		return nil, fmt.Errorf(
			"seek bytes %q not prefixed with %q: %w",
			string(seek),
			string(config.Prefix),
			kv.ErrSeekMissingPrefix,
		)
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
