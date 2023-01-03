package mvcc

import (
	"bytes"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"math"
	"sync"
)

const (
	// markedRevBytesLen is the byte length of marked revision.
	// The first `revBytesLen` bytes represents a normal revision. The last
	// one byte is the mark.
	markedRevBytesLen      = revBytesLen + 1
	markBytePosition       = markedRevBytesLen - 1
	markTombstone     byte = 't'
)

var (
	restoreChunkKeys = 10000 // non-const for testing

	metaBucketName         = []byte("__meta")
	FinishedCompactKeyName = []byte("finishedCompactRev")
)

type Store struct {
	ReadView
	WriteView

	logger *zap.Logger

	// mtx read locks for txns and write locks for non-txn store changes
	mtx     sync.RWMutex
	db      *bbolt.DB
	kvIndex *treeIndex

	// revMtx protects currentRev and compactMainRev.
	// Locked at end of write txn and released after write txn unlock lock.
	// Locked before locing read txn and released after locking.
	revMtx sync.RWMutex
	// currentRev is the revision of the last completed transaction.
	currentRev int64
	// compactMainRev is the main revision of the last compaction
	compactMainRev int64
}

// New returns a new store. It is useful to create a store inside mvcc pkg.
// It should only be used for testing externally.
func New(logger *zap.Logger) (*Store, error) {
	db, err := bbolt.Open("state.bolt", 0600, nil)
	if err != nil {
		return nil, err
	}

	store := &Store{
		logger:         logger,
		db:             db,
		currentRev:     1,
		compactMainRev: -1,
		kvIndex:        newTreeIndex(),
	}

	store.mtx.Lock()
	defer store.mtx.RUnlock()
	if err := store.restore(); err != nil {
		panic("afiled to recover store from boltdb")
	}

	return store, nil
}

func (s *Store) restore() error {
	min, max := newRevBytes(), newRevBytes()
	revToBytes(revision{main: 1}, min)
	revToBytes(revision{main: math.MaxInt64, sub: math.MaxInt64}, max)

	tx, err := s.db.Begin(false)
	if err != nil {
		return err
	}

	finishedCompact, found, err := readFinishedCompact(tx)
	if err != nil {
		return err
	}

	if found {
		s.revMtx.Lock()
		s.compactMainRev = finishedCompact
		s.revMtx.Unlock()

		s.logger.Info("restored last compact revision",
			zap.ByteString("meta-bucket", metaBucketName),
			zap.Int64("restored-compact-revision", finishedCompact),
		)
	}

	rkvc, revc := restoreIntoIndex(s.logger, s.kvIndex)
	for {
		keys, values := scan(tx, nil, min, max, int64(restoreChunkKeys))
		if len(keys) == 0 {
			break
		}

		// rkvc blocks if the total pending keys exceeds the restore
		// chunk size to keep keys from consuming too much memory.

	}
}

func scan(tx *bbolt.Tx, bucket, start, end []byte, limit int64) ([][]byte, [][]byte) {
	b := tx.Bucket(bucket)

	var isMatch func(b []byte) bool
	if len(end) > 0 {
		isMatch = func(b []byte) bool { return bytes.Compare(b, end) < 0 }
	} else {
		isMatch = func(b []byte) bool { return bytes.Equal(b, start) }
		limit = 1
	}

	var (
		keys   [][]byte
		values [][]byte
	)
	c := b.Cursor()
	for k, v := c.Seek(start); k != nil && isMatch(k); k, v = c.Next() {
		keys = append(keys, k)
		values = append(values, v)
		if limit == int64(len(keys)) {
			break
		}
	}

	return keys, values
}

func readFinishedCompact(tx *bbolt.Tx) (int64, bool, error) {
	b := tx.Bucket(metaBucketName)
	value := b.Get(FinishedCompactKeyName)
	if len(value) != 0 {
		return bytesToRev(value).main, true, nil
	}

	return 0, false, nil
}

type revKeyValue struct {
	key  []byte
	kstr string
	kv   KeyValue
}

func restoreIntoIndex(logger *zap.Logger, idx *treeIndex) (chan<- revKeyValue, <-chan int64) {
	rkvc, revc := make(chan revKeyValue, restoreChunkKeys), make(chan int64, 1)

	go func() {
		currentRev := int64(1)
		defer func() {
			revc <- currentRev
		}()

		// restore the tree index from streaming the unordered index.
		kiCache := make(map[string]*keyIndex, restoreChunkKeys)
		for rkv := range rkvc {
			ki, ok := kiCache[rkv.kstr]
			// purge kiCache if many keys but still missing in the cache
			if !ok && len(kiCache) >= restoreChunkKeys {
				i := 10
				for k := range kiCache {
					delete(kiCache, k)
					if i--; i == 0 {
						break
					}
				}
			}

			// cache miss, fetch from tree index if there
			if !ok {
				ki := &keyIndex{key: rkv.kv.Key}
				if idxKey := idx.keyIndex(ki); idxKey != nil {
					kiCache[rkv.kstr], ki = idxKey, idxKey
					ok = true
				}
			}
			rev := bytesToRev(rkv.key)
			currentRev = rev.main
			if !ok {
				if isTombstone(rkv.key) {
					if err := ki.tombstone(rev.main, rev.sub); err != nil {
						logger.Warn("tombstone encountered error", zap.Error(err))
					}

					continue
				}

				ki.put(rev.main, rev.sub)
			} else if !isTombstone(rkv.key) {
				ki.restore(revision{rkv.kv.CreateRevision, 0}, rev, rkv.kv.Version)
				idx.insert(ki)
				kiCache[rkv.kstr] = ki
			}
		}
	}()

	return rkvc, revc
}

// isTombstone checks whether the revision bytes is a tombstone.
func isTombstone(b []byte) bool {
	return len(b) == markedRevBytesLen && b[markBytePosition] == markTombstone
}
