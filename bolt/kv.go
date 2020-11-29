package bolt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/prometheus/client_golang/prometheus"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

// check that *KVStore implement kv.SchemaStore interface.
var _ kv.SchemaStore = (*KVStore)(nil)

// KVStore is a kv.Store backed by boltdb.
type KVStore struct {
	path string
	db   *bolt.DB
	log  *zap.Logger

	noSync bool

	// metrics
	commitSec prometheus.Histogram
	writeSec  prometheus.Histogram

	updateSec prometheus.Histogram
	viewSec   prometheus.Histogram
}

type KVOption func(*KVStore)

// WithNoSync WARNING: this is useful for tests only
// this skips fsyncing on every commit to improve
// write performance in exchange for no guarantees
// that the db will persist.
func WithNoSync(s *KVStore) {
	s.noSync = true
}

// NewKVStore returns an instance of KVStore with the file at
// the provided path.
func NewKVStore(log *zap.Logger, path string, opts ...KVOption) *KVStore {
	store := &KVStore{
		path: path,
		log:  log,
	}

	for _, opt := range opts {
		opt(store)
	}

	store.commitSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "manta",
		Subsystem: "boltdb",
		Name:      "commit_duration_seconds",
		Help:      "The latency distributions of commit called by bbolt",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	store.writeSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "manta",
		Subsystem: "boltdb",
		Name:      "commit_write_duration_seconds",
		Help:      "The latency distributions of commit.write called by bbolt",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	store.updateSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "manta",
		Subsystem: "boltdb",
		Name:      "update_duration_seconds",
		Help:      "The latency distributions of update called by kv",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	store.viewSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "manta",
		Subsystem: "boltdb",
		Name:      "view_duration_seconds",
		Help:      "The latency distributions of view called by kv",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	return store
}

// Open creates boltDB file it doesn't exists and opens it otherwise.
func (s *KVStore) Open(ctx context.Context) error {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	// Ensure the required directory structure exists.
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return fmt.Errorf("unable to create directory %s: %v", s.path, err)
	}

	if _, err := os.Stat(s.path); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Open database file.
	db, err := bolt.Open(s.path, 0600, &bolt.Options{Timeout: 1 * time.Second, InitialMmapSize: 32 * 1024 * 1024})
	if err != nil {
		return fmt.Errorf("unable to open boltdb file %v", err)
	}
	s.db = db

	db.NoSync = s.noSync

	return nil
}

// Close the connection to the bolt database
func (s *KVStore) Close() error {
	if s.db != nil {
		// sync before close
		if err := s.db.Sync(); err != nil {
			return err
		}

		return s.db.Close()
	}
	return nil
}

func (s *KVStore) DB() *bolt.DB {
	return s.db
}

// Flush removes all bolt keys within each bucket.
func (s *KVStore) Flush(ctx context.Context) {
	_ = s.db.Update(
		func(tx *bolt.Tx) error {
			return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
				s.cleanBucket(tx, b)
				return nil
			})
		},
	)
}

func (s *KVStore) cleanBucket(tx *bolt.Tx, b *bolt.Bucket) {
	// nested bucket recursion base case:
	if b == nil {
		return
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		_ = v
		if err := c.Delete(); err != nil {
			// clean out nexted buckets
			s.cleanBucket(tx, b.Bucket(k))
		}
	}
}

// WithDB sets the boltdb on the store.
func (s *KVStore) WithDB(db *bolt.DB) {
	s.db = db
}

// View opens up a view transaction against the store.
func (s *KVStore) View(ctx context.Context, fn func(tx kv.Tx) error) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	start := time.Now()
	tx, err := s.db.Begin(false)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = fn(&Tx{
		tx:  tx,
		ctx: ctx,
	})

	s.viewSec.Observe(time.Since(start).Seconds())

	return err
}

// Update opens up an update transaction against the store.
func (s *KVStore) Update(ctx context.Context, fn func(tx kv.Tx) error) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	start := time.Now()
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}

	defer tx.Commit()

	err = fn(&Tx{
		tx:  tx,
		ctx: ctx,
	})

	s.updateSec.Observe(time.Since(start).Seconds())
	s.commitSec.Observe(time.Since(start).Seconds())
	s.writeSec.Observe(tx.Stats().WriteTime.Seconds())

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// CreateBucket creates a bucket in the underlying boltdb store if it
// does not already exist
func (s *KVStore) CreateBucket(ctx context.Context, name []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
}

// DeleteBucket creates a bucket in the underlying boltdb store if it
// does not already exist
func (s *KVStore) DeleteBucket(ctx context.Context, name []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket(name); err != nil && !errors.Is(err, bolt.ErrBucketNotFound) {
			return err
		}

		return nil
	})
}

// Backup copies all K:Vs to a writer, in BoltDB format.
func (s *KVStore) Backup(ctx context.Context, w io.Writer) error {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.db.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(w)
		return err
	})
}

// Tx is a light wrapper around a boltdb transaction. It implements kv.Tx.
type Tx struct {
	tx  *bolt.Tx
	ctx context.Context
}

// Context returns the context for the transaction.
func (tx *Tx) Context() context.Context {
	return tx.ctx
}

// WithContext sets the context for the transaction.
func (tx *Tx) WithContext(ctx context.Context) {
	tx.ctx = ctx
}

// Bucket retrieves the bucket named b.
func (tx *Tx) Bucket(b []byte) (kv.Bucket, error) {
	bkt := tx.tx.Bucket(b)
	if bkt == nil {
		return nil, fmt.Errorf("bucket %q: %w", string(b), kv.ErrBucketNotFound)
	}
	return &Bucket{
		bucket: bkt,
	}, nil
}

// Bucket implements kv.Bucket.
type Bucket struct {
	bucket *bolt.Bucket
}

// Get retrieves the value at the provided key.
func (b *Bucket) Get(key []byte) ([]byte, error) {
	val := b.bucket.Get(key)
	if len(val) == 0 {
		return nil, kv.ErrKeyNotFound
	}

	return val, nil
}

// GetBatch retrieves the values for the provided keys.
func (b *Bucket) GetBatch(keys ...[]byte) ([][]byte, error) {
	values := make([][]byte, len(keys))
	for idx, key := range keys {
		val := b.bucket.Get(key)
		if len(val) == 0 {
			continue
		}

		values[idx] = val
	}

	return values, nil
}

// Put sets the value at the provided key.
func (b *Bucket) Put(key []byte, value []byte) error {
	err := b.bucket.Put(key, value)
	if err == bolt.ErrTxNotWritable {
		return kv.ErrTxNotWritable
	}
	return err
}

// Delete removes the provided key.
func (b *Bucket) Delete(key []byte) error {
	err := b.bucket.Delete(key)
	if err == bolt.ErrTxNotWritable {
		return kv.ErrTxNotWritable
	}
	return err
}

// ForwardCursor retrieves a cursor for iterating through the entries
// in the key value store in a given direction (ascending / descending).
func (b *Bucket) ForwardCursor(seek []byte, opts ...kv.CursorOption) (kv.ForwardCursor, error) {
	var (
		cursor     = b.bucket.Cursor()
		config     = kv.NewCursorConfig(opts...)
		key, value []byte
	)

	if len(seek) == 0 && config.Direction == kv.CursorDescending {
		seek, _ = cursor.Last()
	}

	key, value = cursor.Seek(seek)

	if config.Prefix != nil && !bytes.HasPrefix(seek, config.Prefix) {
		return nil, fmt.Errorf("seek bytes %q not prefixed with %q: %w", string(seek), string(config.Prefix), kv.ErrSeekMissingPrefix)
	}

	c := &Cursor{
		cursor: cursor,
		config: config,
	}

	// only remember first seeked item if not skipped
	if !config.SkipFirst {
		c.key = key
		c.value = value
	}

	return c, nil
}

// Cursor retrieves a cursor for iterating through the entries
// in the key value store.
func (b *Bucket) Cursor(opts ...kv.CursorHint) (kv.Cursor, error) {
	return &Cursor{
		cursor: b.bucket.Cursor(),
	}, nil
}

// Cursor is a struct for iterating through the entries
// in the key value store.
type Cursor struct {
	cursor *bolt.Cursor

	// previously seeked key/value
	key, value []byte

	config kv.CursorConfig
	closed bool
	seen   int
}

// Close sets the closed to closed
func (c *Cursor) Close() error {
	c.closed = true

	return nil
}

// Seek seeks for the first key that matches the prefix provided.
func (c *Cursor) Seek(prefix []byte) ([]byte, []byte) {
	if c.closed {
		return nil, nil
	}
	k, v := c.cursor.Seek(prefix)
	if len(k) == 0 && len(v) == 0 {
		return nil, nil
	}
	return k, v
}

// First retrieves the first key value pair in the bucket.
func (c *Cursor) First() ([]byte, []byte) {
	if c.closed {
		return nil, nil
	}
	k, v := c.cursor.First()
	if len(k) == 0 && len(v) == 0 {
		return nil, nil
	}
	return k, v
}

// Last retrieves the last key value pair in the bucket.
func (c *Cursor) Last() ([]byte, []byte) {
	if c.closed {
		return nil, nil
	}
	k, v := c.cursor.Last()
	if len(k) == 0 && len(v) == 0 {
		return nil, nil
	}
	return k, v
}

// Next retrieves the next key in the bucket.
func (c *Cursor) Next() (k []byte, v []byte) {
	if c.closed ||
		c.atLimit() ||
		(c.key != nil && c.missingPrefix(c.key)) {
		return nil, nil
	}

	// get and unset previously seeked values if they exist
	k, v, c.key, c.value = c.key, c.value, nil, nil
	if len(k) > 0 || len(v) > 0 {
		c.seen++
		return
	}

	next := c.cursor.Next
	if c.config.Direction == kv.CursorDescending {
		next = c.cursor.Prev
	}

	k, v = next()
	if (len(k) == 0 && len(v) == 0) || c.missingPrefix(k) {
		return nil, nil
	}

	c.seen++

	return k, v
}

// Prev retrieves the previous key in the bucket.
func (c *Cursor) Prev() (k []byte, v []byte) {
	if c.closed ||
		c.atLimit() ||
		(c.key != nil && c.missingPrefix(c.key)) {
		return nil, nil
	}

	// get and unset previously seeked values if they exist
	k, v, c.key, c.value = c.key, c.value, nil, nil
	if len(k) > 0 && len(v) > 0 {
		c.seen++
		return
	}

	prev := c.cursor.Prev
	if c.config.Direction == kv.CursorDescending {
		prev = c.cursor.Next
	}

	k, v = prev()
	if (len(k) == 0 && len(v) == 0) || c.missingPrefix(k) {
		return nil, nil
	}

	c.seen++

	return k, v
}

func (c *Cursor) missingPrefix(key []byte) bool {
	return c.config.Prefix != nil && !bytes.HasPrefix(key, c.config.Prefix)
}

func (c *Cursor) atLimit() bool {
	return c.config.Limit != nil && c.seen >= *c.config.Limit
}

// Err always returns nil as nothing can go wrong™ during iteration
func (c *Cursor) Err() error {
	return nil
}
