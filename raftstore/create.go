package raftstore

import (
	"context"
	"github.com/f1shl3gs/manta/raftstore/raft"
	"sync"
	"time"
)

const (
	// initialMmapSize is the initial size of the mmapped region. Setting this
	// larger than the potential max db size can prevent writer from blocking reader.
	//
	// This only works for linux
	initialMmapSize = 64 * 1024 * 1024 // 64MiB
)

func NewV1(cf *Config) error {
	db, err := setupDB(cf)
	if err != nil {
		return err
	}

	raftNode, err := raft.New(db)
	if err != nil {
		return err
	}

	kv := &KV{
		db:    db,
		raft:  raftNode,
		idGen: newGenerator(0, time.Now()),
		wait:  newWait(),
	}

	return nil
}

func (s *KV) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.logger.Info("compact")
			}
		}
	}()

	wg.Wait()

	return nil
}
