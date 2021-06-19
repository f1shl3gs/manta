package kv

import (
	"context"

	"github.com/f1shl3gs/manta"
)

var (
	secretBucket         = []byte("secrets")
	secretOrgIndexBucket = []byte("secretorgindex")
)

func (s *Service) GetSecretKeys(ctx context.Context, orgID manta.ID) ([]string, error) {
	var (
		keys []string
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		keys, err = s.getSecretKeys(tx, orgID)
		return err
	})

	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *Service) getSecretKeys(tx Tx, orgID manta.ID) ([]string, error) {
	b, err := tx.Bucket(secretOrgIndexBucket)
	if err != nil {
		return nil, err
	}

	prefix, err := orgID.Encode()
	if err != nil {
		return nil, err
	}

	cursor, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix))
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	keys := make([]string, 0, 8)
	for {
		k, v := cursor.Next()
		if k == nil {
			break
		}

		keys = append(keys, string(v))
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

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
		_, err := s.findOrganizationByID(ctx, tx, orgID)
		if err != nil {
			return &manta.Error{
				Code: manta.ENotFound,
				Op:   "PutSecret",
				Err:  err,
			}
		}

		return s.putSecret(ctx, tx, orgID, k, v)
	})
}

// putSecret save the sec
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
