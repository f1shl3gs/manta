package raftstore

import (
    "context"
    "encoding/binary"
    "fmt"
    "github.com/f1shl3gs/manta/raftstore/membership"
    "github.com/f1shl3gs/manta/raftstore/raft"
    "github.com/f1shl3gs/manta/raftstore/wal"
    bolt "go.etcd.io/bbolt"
    "go.etcd.io/raft/v3/raftpb"
    "go.uber.org/zap"
    "sync/atomic"
)

var (
    metaBucket = []byte("__meta")

    consistentIndexKey = []byte("consistentIndex")
)

type Store struct {
    logger *zap.Logger

    db atomic.Pointer[bolt.DB]
    raftNode *Node

    cluster *membership.Cluster
}


func New(cf *Config, logger *zap.Logger) (*Store, error) {
    store := &Store{
        logger: logger.With(zap.String("service", "raftstore")),
    }

    if db, err := openDB(cf); err != nil {
        return nil, err
    } else {
        store.db.Store(db)
    }

    ds, err := wal.Init(cf.DataDir, logger)
    if err != nil {
        return nil, err
    }

    // TODO: what if there is multiple term in ds
    term, index, err := store.consistentIndex()
    if err != nil {
        panic("read consistent index failed")
    }

    // if node never start before, we don't need to replay
    if term == 0 && index == 0 {
        // TODO:
        return store, nil
    }

    // replaying unapplied entry to backend and restore cluster from it
    lastIndex, err := ds.LastIndex()
    if err != nil {
        return nil, err
    }

    entries, err := ds.Entries(index, lastIndex, lastIndex - index)
    if err != nil {
        return nil, err
    }

    for _, ent := range entries {
        if ent.Index <= index {
            panic("replying entry smaller than current backend")
        }


    }
}

func (s *Store) consistentIndex() (term, index uint64, err error) {
    err = s.db.Load().View(func(tx *bolt.Tx) error {
        b := tx.Bucket(metaBucket)
        value := b.Get(consistentIndexKey)

        // first setup
        if len(value) == 0 {
            return nil
        }

        if len(value) != 16 {
            term = binary.BigEndian.Uint64(value)
            index = binary.BigEndian.Uint64(value[8:])
        } else {
            panic(fmt.Sprintf("consistent value is not 16 bytes"))
        }

        return nil
    })

    return
}

func (s *Store) Run(ctx context.Context) {
    applyCh := s.raftNode.Apply()

    // apply loop
    for {
        select {
        case <-ctx.Done():
            return
        case apply := <- applyCh:
            err := s.db.Load().Update(func(tx *bolt.Tx) error {
                for _, ent := range apply.entries {
                    switch ent.Type {

                    }
                }
            })
            if err != nil {
                s.logger.Fatal("error while applying to boltdb",
                    zap.Error(err))
            }
        }
    }
}

func (s *Store) applyConfChange(tx *bolt.Tx, ent *raftpb.Entry) {
    var cc raftpb.ConfChange
    MustUnmarshal(&cc, ent.Data)

    // update memory cached cluster and store it to bolt
    switch cc.Type {
    case raftpb.ConfChangeAddNode:
        s.cluster.
   }

   s.raftNode.Node.ApplyConfChange(cc)
}


func (s *Store) applyNormal(tx *bolt.Tx) error {

}