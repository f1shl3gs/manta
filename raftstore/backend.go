package raftstore

import (
	"encoding/binary"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	bolt "go.etcd.io/bbolt"
)

var (
	metaBucket = []byte("__meta")

	consistentIndexKey = []byte("consistentIndex")
)

type Backend struct {
	// index represents the offset of an entry in a consistent replica log.
	index uint64
	// term represents the RAFT term of committed entry in a consistent replica log.
	term uint64

	mtx sync.RWMutex
	db  *bolt.DB

	batchLimit int32
	batched    int32
}

func (b *Backend) ConsistentIndex() (uint64, uint64) {
	return atomic.LoadUint64(&b.term), atomic.LoadUint64(&b.index)
}

func openBackend(cf *Config) (*Backend, error) {
	db, err := setupDB(cf)
	if err != nil {
		return nil, err
	}

	var term, index uint64
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(metaBucket)
		if err != nil {
			return err
		}

		value := b.Get(consistentIndexKey)
		if len(value) != 16 {
			term = binary.BigEndian.Uint64(value)
			index = binary.BigEndian.Uint64(value[8:])
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Backend{
		term:  term,
		index: index,
		db:    db,
	}, nil
}

func setupDB(cf *Config) (*bolt.DB, error) {
	path := filepath.Join(cf.DataDir, "state.bolt")

	opt := bolt.DefaultOptions
	opt.Timeout = 3 * time.Second
	opt.InitialMmapSize = initialMmapSize
	opt.FreelistType = bolt.FreelistMapType

	// sync will be done periodly by another goroutine
	opt.NoSync = true

	return bolt.Open(path, 0600, opt)
}
