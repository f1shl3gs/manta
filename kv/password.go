package kv

import (
	"context"

	"github.com/f1shl3gs/manta"
)

var (
	passwordBucket = []byte("passwords")
)

func (s *Service) SetPassword(ctx context.Context, uid manta.ID, password string) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		pk, err := uid.Encode()
		if err != nil {
			return err
		}

		b, err := tx.Bucket(passwordBucket)
		if err != nil {
			return err
		}

		return b.Put(pk, []byte(password))
	})
}

func (s *Service) ComparePassword(ctx context.Context, uid manta.ID, password string) error {
	return s.kv.View(ctx, func(tx Tx) error {
		pk, err := uid.Encode()
		if err != nil {
			return err
		}

		b, err := tx.Bucket(passwordBucket)
		if err != nil {
			return err
		}

		v, err := b.Get(pk)
		if err != nil {
			return err
		}

		if string(v) == password {
			return nil
		}

		return manta.ErrPasswordNotMatch
	})
}

func (s *Service) CompareAndSetPassword(ctx context.Context, uid manta.ID, old, new string) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		pk, err := uid.Encode()
		if err != nil {
			return err
		}

		b, err := tx.Bucket(passwordBucket)
		if err != nil {
			return err
		}

		v, err := b.Get(pk)
		if err != nil {
			return err
		}

		if string(v) == old {
			return nil
		}

		return b.Put(pk, []byte(new))
	})
}

func (s *Service) DeletePassword(ctx context.Context, uid manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		pk, err := uid.Encode()
		if err != nil {
			return err
		}

		b, err := tx.Bucket(passwordBucket)
		if err != nil {
			return err
		}

		return b.Delete(pk)
	})
}
