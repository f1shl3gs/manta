package raftstore

import (
	"context"
	"testing"
	"time"

	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/tests"

	"github.com/stretchr/testify/assert"
)

func initStore(fields tests.KVStoreFields, t *testing.T) (kv.Store, func()) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	store := setupStore(t)
	go store.Run(ctx)

	err := store.CreateBucket(ctx, fields.Bucket)
	assert.NoError(t, err)

	err = store.Update(ctx, func(tx kv.Tx) error {
		b, err := tx.Bucket(fields.Bucket)
		if err != nil {
			return err
		}

		for i := range fields.Pairs {
			err := b.Put(fields.Pairs[i].Key, fields.Pairs[i].Value)
			if err != nil {
				return err
			}
		}

		return nil
	})
	assert.NoError(t, err)

	return store, store.stop
}

func TestKV(t *testing.T) {
	tests.KVStore(initStore, t)
}
