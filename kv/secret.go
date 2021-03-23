package kv

import (
	"context"

	"github.com/f1shl3gs/manta"
)

var (
	secretBucket = []byte("secrets")
)

func (s *Service) FindSecret(ctx context.Context, orgID manta.ID, k string) (string, error) {
	var (
		value string
		err   error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		value, err = s.findSecret(ctx, tx, orgID, k)
		return err
	})

	if err != nil {
		return "", err
	}

	return value, err
}

func (s *Service) findSecret(ctx context.Context, tx Tx, orgID manta.ID, k string) (string, error) {
	prefix, err := orgID.Encode()
	if err != nil {
		return "", err
	}

	b, err := tx.Bucket(secretBucket)
	if err != nil {
		return "", err
	}

	v, err := b.Get(IndexKey(prefix, []byte(k)))
	if err != nil {
		return "", err
	}

	return string(v), nil
}

func (s *Service) PutSecret(ctx context.Context, orgID manta.ID, k, v string) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.putSecret(ctx, tx, orgID, k, v)
	})
}

func (s *Service) putSecret(ctx context.Context, tx Tx, orgID manta.ID, k, v string) error {
	b, err := tx.Bucket(secretBucket)
	if err != nil {
		return err
	}

	pk, err := orgID.Encode()
	if err != nil {
		return err
	}

	return b.Put(IndexKey(pk, []byte(k)), []byte(v))
}

func (s *Service) DeleteSecret(ctx context.Context, orgID manta.ID, keys ...string) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteSecret(ctx, tx, orgID, keys...)
	})
}

func (s *Service) deleteSecret(ctx context.Context, tx Tx, orgID manta.ID, keys ...string) error {
	prefix, err := orgID.Encode()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(secretBucket)
	if err != nil {
		return err
	}

	for _, k := range keys {
		_ = b.Delete(IndexKey(prefix, []byte(k)))
	}

	return nil
}
