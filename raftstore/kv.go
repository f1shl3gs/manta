package raftstore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"unsafe"

	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/raftstore/pb"

	bolt "go.etcd.io/bbolt"
)

// CreateBucket creates a bucket on the underlying store if it does not exist
func (s *Store) CreateBucket(ctx context.Context, bucket []byte) error {
	return s.propose(ctx, pb.InternalRequest{
		CreateBucket: &pb.CreateBucket{
			Name: bucket,
		},
	})
}

// DeleteBucket deletes a bucket on the underlying store if it exists
func (s *Store) DeleteBucket(ctx context.Context, bucket []byte) error {
	return s.propose(ctx, pb.InternalRequest{
		DeleteBucket: &pb.DeleteBucket{
			Name: bucket,
		},
	})
}

// View opens up a transaction that will not write to any value. Implementing interfaces
// should take care to ensure that all view transactions do not mutate any value.
func (s *Store) View(ctx context.Context, fn func(kv.Tx) error) error {
	err := s.linearizableReadNotify(ctx)
	if err != nil {
		return err
	}

	db := s.db.Load()
	tx, err := db.Begin(false)
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
func (s *Store) Update(ctx context.Context, fn func(kv.Tx) error) error {
	// TODO: some transication do not need to read anything, so we can move
	// this when read happends
	err := s.linearizableReadNotify(ctx)
	if err != nil {
		return err
	}

	// write operation is cached and it will be propose throught raft,
	// the it will be applied, so we don't need this tx to be writable.
	db := s.db.Load()
	tx, err := db.Begin(false)
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

	return s.propose(ctx, pb.InternalRequest{
		Txn: txn,
	})
}

// Backup copies all K:Vs to a writer, file format determined by implementation.
func (s *Store) Backup(ctx context.Context, w io.Writer) error {
	panic("not implement")
}

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
	val := b.bucket.Get(key)
	if len(val) == 0 {
		return nil, kv.ErrKeyNotFound
	}

	return val, nil
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
		c.seen++
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

	c.seen++

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
		return nil, fmt.Errorf("seek bytes %q not prefixed with %q: %w",
			string(seek), string(config.Prefix), kv.ErrSeekMissingPrefix)
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

type valueItem struct {
	version int64
	value   []byte
}

type readSet map[string]valueItem

func newReadSet() readSet {
	return make(map[string]valueItem)
}

func (rs readSet) add(key, value []byte, version int64) {
	rs[unsafeBytesToString(key)] = valueItem{
		version: version,
		value:   value,
	}
}

// first returns the store version from the first fetch
// func (rs readSet) first() int64 {
// 	ret := int64(math.MaxInt64 - 1)
// 	for _, item := range rs {
// 		if ret < item.version {
// 			ret = item.version
// 		}
// 	}
//
// 	return ret
// }

func unsafeBytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func unsafeStringToBytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(&s))
}

type writeOp struct {
	value    []byte
	deletion bool
}

type writeSet map[string]writeOp

func newWriteSet() writeSet {
	return make(map[string]writeOp)
}

type writeTx struct {
	ctx context.Context

	tx *bolt.Tx

	// rset holds read key values and versions
	rset map[string]readSet
	// wset holds overwritten keys and their values
	wset map[string]writeSet
}

func (tx *writeTx) txn() *pb.Txn {
	txn := &pb.Txn{}

	for b, rset := range tx.rset {
		for key, item := range rset {
			txn.Compares = append(txn.Compares, &pb.Compare{
				Bucket:  unsafeStringToBytes(b),
				Key:     unsafeStringToBytes(key),
				Version: item.version,
			})
		}
	}

	for b, wset := range tx.wset {
		for key, item := range wset {
			op := &pb.Operation{
				Bucket: unsafeStringToBytes(b),
				Key:    unsafeStringToBytes(key),
			}

			if item.deletion {
				op.Type = pb.Delete
			} else {
				op.Value = item.value
			}

			txn.Successes = append(txn.Successes, op)
		}
	}

	return txn
}

func (tx *writeTx) Bucket(b []byte) (kv.Bucket, error) {
	bs := unsafeBytesToString(b)

	ro := tx.tx.Bucket(b)
	if ro == nil {
		return nil, kv.ErrBucketNotFound
	}

	rset, ok := tx.rset[bs]
	if !ok {
		rset = newReadSet()
		tx.rset[bs] = rset
	}

	wset, ok := tx.wset[bs]
	if !ok {
		wset = newWriteSet()
		tx.wset[bs] = wset
	}

	// always assume bucket exist
	return &bucket{
		name:   b,
		bucket: ro,
		rset:   rset,
		wset:   wset,
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
	if op, exist := b.wset[sk]; exist {
		if op.deletion {
			return nil, kv.ErrKeyNotFound
		}

		return op.value, nil
	}

	if item, exist := b.rset[sk]; exist {
		return item.value, nil
	}

	val := b.bucket.Get(key)
	if len(val) == 0 {
		return nil, kv.ErrKeyNotFound
	}

	b.rset.add(key, val, 0)

	return val, nil
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
	sk := unsafeBytesToString(key)
	b.wset[sk] = writeOp{
		value:    nil,
		deletion: true,
	}

	delete(b.rset, sk)

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
