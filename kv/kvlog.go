package kv

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"time"

	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	kvlogBucket = []byte("kvlog")

	ErrInvalidKey = errors.New("invalid key")
)

func encodeLogEntryPrefix(k []byte) []byte {
	h := sha1.New()
	h.Write(k)
	return h.Sum(nil)
}

func encodeLogEntryKey(k []byte, t time.Time) []byte {
	prefix := encodeLogEntryPrefix(k)
	ek := make([]byte, 20+8)
	copy(prefix, k)

	binary.BigEndian.PutUint64(ek[20:], uint64(t.UnixNano()))

	return ek
}

func decodeLogEntryTimestamp(key []byte) (time.Time, error) {
	if len(key) != 28 {
		return time.Time{}, ErrInvalidKey
	}

	t := binary.BigEndian.Uint64(key[20:])
	return time.Unix(0, int64(t)), nil
}

func (s *Service) AddLogEntry(ctx context.Context, k, v []byte, t time.Time) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.addLogEntry(ctx, tx, k, v, t)
	})
}

func (s *Service) addLogEntry(ctx context.Context, tx Tx, k, v []byte, t time.Time) error {
	ek := encodeLogEntryKey(k, t)
	b, err := tx.Bucket(kvlogBucket)
	if err != nil {
		return err
	}

	return b.Put(ek, v)
}

func (s *Service) ForEachLogEntry(ctx context.Context, k []byte, fn func(v []byte, t time.Time) error) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.kv.View(ctx, func(tx Tx) error {
		return s.forEachLogEntry(ctx, tx, k, fn)
	})
}

func (s *Service) forEachLogEntry(ctx context.Context, tx Tx, k []byte, fn func(v []byte, t time.Time) error) error {
	b, err := tx.Bucket(kvlogBucket)
	if err != nil {
		return err
	}

	prefix := encodeLogEntryPrefix(k)
	c, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix))
	if err != nil {
		return err
	}

	return WalkCursor(ctx, c, func(k, v []byte) error {
		ts, err := decodeLogEntryTimestamp(k)
		if err != nil {
			return err
		}

		return fn(v, ts)
	})
}
