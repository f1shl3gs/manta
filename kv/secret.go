package kv

import (
	"context"
	"encoding/json"
	"time"

	"github.com/f1shl3gs/manta"
)

var (
	SecretsBucket = []byte("secrets")
)

func (s *Service) LoadSecret(ctx context.Context, orgID manta.ID, k string) (*manta.Secret, error) {
	var (
		secret *manta.Secret
		err    error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		secret, err = loadSecret(tx, orgID, k)
		return err
	})

	if err != nil {
		return nil, err
	}

	return secret, nil
}

func loadSecret(tx Tx, orgID manta.ID, k string) (*manta.Secret, error) {
	key, err := secretKey(orgID, k)
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(SecretsBucket)
	if err != nil {
		return nil, err
	}

	value, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	secret := &manta.Secret{}
	err = json.Unmarshal(value, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (s *Service) GetSecrets(ctx context.Context, orgID manta.ID) ([]manta.Secret, error) {
	var secrets = make([]manta.Secret, 0)

	err := s.kv.View(ctx, func(tx Tx) error {
		b, err := tx.Bucket(SecretsBucket)
		if err != nil {
			return err
		}

		fk, err := orgID.Encode()
		if err != nil {
			return err
		}

		cursor, err := b.Cursor(WithCursorHintPrefix(string(fk)))
		if err != nil {
			return err
		}

		for k, value := cursor.Seek(fk); k != nil; k, value = cursor.Next() {
			secret := manta.Secret{}
			err = json.Unmarshal(value, &secret)
			if err != nil {
				return err
			}

			secret.Desensitize()
			secrets = append(secrets, secret)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func (s *Service) putSecret(tx Tx, secret *manta.Secret) error {
	key, err := secretKey(secret.OrgID, secret.Key)
	if err != nil {
		return err
	}

	value, err := json.Marshal(secret)
	if err != nil {
		return err
	}

	b, err := tx.Bucket(SecretsBucket)
	if err != nil {
		return err
	}

	return b.Put(key, value)
}

func (s *Service) PutSecret(ctx context.Context, secret *manta.Secret) (*manta.Secret, error) {
	secret.Updated = time.Now()

	err := s.kv.Update(ctx, func(tx Tx) error {
		return s.putSecret(tx, secret)
	})

	if err != nil {
		return nil, err
	}

	secret.Desensitize()

	return secret, nil
}

func (s *Service) DeleteSecret(ctx context.Context, orgID manta.ID, keys ...string) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		b, err := tx.Bucket(SecretsBucket)
		if err != nil {
			return err
		}

		for _, k := range keys {
			key, err := secretKey(orgID, k)
			if err != nil {
				return err
			}

			if err = b.Delete(key); err != nil {
				return err
			}
		}

		return nil
	})
}

func secretKey(orgID manta.ID, k string) ([]byte, error) {
	o, err := orgID.Encode()
	if err != nil {
		return nil, err
	}

	return IndexKey(o, []byte(k)), nil
}
