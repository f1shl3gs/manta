package kv

import (
	"bytes"
	"context"
)

var (
	keyringBucket = []byte("keyring")
)

func (s *Service) AddKey(ctx context.Context, key []byte) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		b, err := tx.Bucket(keyringBucket)
		if err != nil {
			return err
		}

		id := s.idGen.ID()
		sk, err := id.Encode()
		if err != nil {
			return err
		}

		return b.Put(sk, key)
	})
}

func (s *Service) PrimaryKey(ctx context.Context) ([]byte, error) {
	var (
		key []byte
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		b, err := tx.Bucket(keyringBucket)
		if err != nil {
			return err
		}

		c, err := b.Cursor()
		if err != nil {
			return err
		}

		_, key = c.Last()
		if key == nil {
			return ErrKeyNotFound
		}

		return nil
	})

	return key, err
}

func (s *Service) Keys(ctx context.Context) ([][]byte, error) {
	var (
		keys [][]byte
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		b, err := tx.Bucket(keyringBucket)
		if err != nil {
			return err
		}

		c, err := b.Cursor()
		if err != nil {
			return err
		}

		for k, v := c.First(); k != nil; k, v = c.Next() {
			keys = append(keys, v)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *Service) RemoveKey(ctx context.Context, key []byte) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		b, err := tx.Bucket(keyringBucket)
		if err != nil {
			return err
		}

		c, err := b.Cursor()
		if err != nil {
			return err
		}

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Equal(key, v) {
				return b.Delete(k)
			}
		}

		return ErrKeyNotFound
	})
}
