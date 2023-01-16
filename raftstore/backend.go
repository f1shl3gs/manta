package raftstore

import (
	"encoding/binary"
    "fmt"
    "github.com/gogo/protobuf/proto"
    "go.etcd.io/raft/v3/raftpb"
    "go.uber.org/zap"
    "path/filepath"
	"sync/atomic"
	"time"

	bolt "go.etcd.io/bbolt"
)



type Backend struct {
    logger *zap.Logger

    db atomic.Pointer[bolt.DB]

	batchLimit int32
	batched    int32

    term, index uint64
}

func (b *Backend) DB() *bolt.DB {
    return b.db.Load()
}

// openDB open a boltdb with default options, and setup(if none) meta buckets
// to store membership and consistent index.
func openDB(cf *Config) (*bolt.DB, error) {
	path := filepath.Join(cf.DataDir, "state.bolt")

	opt := bolt.DefaultOptions
	opt.Timeout = 3 * time.Second
	opt.InitialMmapSize = initialMmapSize
	opt.FreelistType = bolt.FreelistMapType

	// sync will be done periodly by another goroutine
	opt.NoSync = true

	db, err := bolt.Open(path, 0600, opt)
    if err != nil {
        return nil, err
    }

    err = db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists(metaBucket)
        return err
    })
    if err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
}

func (b *Backend) consistentIndex() (term, index uint64) {
    db := b.db.Load()

    err := db.View(func(tx *bolt.Tx) error {
        term, index = readConsistentIndex(tx)
        return nil
    })

    if err != nil {
        panic(fmt.Sprintf("read consistent index failed, %s", err))
    }

    return term, index
}


func MustUnmarshal(um proto.Unmarshaler, data []byte) {
    if err := um.Unmarshal(data); err != nil {
        panic(fmt.Sprintf("unmarshal should never fail (%v)", err))
    }
}