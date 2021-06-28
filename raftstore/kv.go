package raftstore

import (
	"context"
	"io"

	"github.com/f1shl3gs/manta/kv"
)

type bucket struct {
	ctx  context.Context
	name []byte
}

func rawKey(bucket, key []byte) []byte {
	rk := make([]byte, len(bucket)+len(key)+1)
	copy(rk, bucket)
	rk[len(bucket)+1] = '/'
	copy(rk[len(bucket)+1:], key)

	return rk
}

func (b *bucket) Get(key []byte) ([]byte, error) {
	panic("implement me")
}

func (b *bucket) GetBatch(keys ...byte) ([][]byte, error) {
	panic("implement me")
}

func (b *bucket) Cursor(hints ...kv.CursorHint) (kv.Cursor, error) {
	panic("implement me")
}

func (b *bucket) Put(key, value []byte) error {
	panic("implement me")
}

func (b *bucket) Delete(key []byte) error {
	panic("implement me")
}

func (b *bucket) ForwardCursor(seek []byte, opts ...kv.CursorOption) (kv.ForwardCursor, error) {
	panic("implement me")
}

type txn struct {
	readOnly bool
	ctx      context.Context

	// TODO: add readSet and writeSet
}

func (s *Store) View(ctx context.Context, f func(kv.Tx) error) error {
	panic("implement me")
}

func (s *Store) Update(ctx context.Context, f func(kv.Tx) error) error {
	panic("implement me")
}

func (s *Store) Backup(ctx context.Context, w io.Writer) error {
	panic("implement me")
}

func (s *Store) CreateBucket(ctx context.Context, bucket []byte) error {
	panic("implement me")
}

func (s *Store) DeleteBucket(ctx context.Context, bucket []byte) error {
	panic("implement me")
}
