package kv

import (
	"context"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	userBucket          = []byte("user")
	userNameIndexBucket = []byte("usernameindex")
)

func (s *Service) FindUserByID(ctx context.Context, id manta.ID) (*manta.User, error) {
	var (
		user *manta.User
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		user, err = s.findUserByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) findUserByID(ctx context.Context, tx Tx, id manta.ID) (*manta.User, error) {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	key, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(key)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	u := &manta.User{}
	if err := u.Unmarshal(data); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) FindUser(ctx context.Context, filter manta.UserFilter) (*manta.User, error) {
	var (
		user *manta.User
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.ID != nil {
			user, err = s.findUserByID(ctx, tx, *filter.ID)
			return err
		}

		if filter.Name != nil {
			user, err = s.findUserByName(ctx, tx, *filter.Name)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) findUserByName(ctx context.Context, tx Tx, name string) (*manta.User, error) {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	fk := []byte(name)

	b, err := tx.Bucket(userNameIndexBucket)
	if err != nil {
		return nil, err
	}

	pk, err := b.Get(fk)
	if err != nil {
		return nil, err
	}

	b, err = tx.Bucket(userBucket)
	if err != nil {
		return nil, err
	}

	val, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	u := &manta.User{}
	if err = u.Unmarshal(val); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) FindUsers(ctx context.Context, filter manta.UserFilter, opts ...manta.FindOptions) ([]*manta.User, error) {
	return nil, nil
}

func (s *Service) CreateUser(ctx context.Context, user *manta.User) error {
	panic("implement me")
}

func (s *Service) UpdateUser(ctx context.Context, id manta.ID, udp manta.UserUpdate) (*manta.User, error) {
	panic("implement me")
}

func (s *Service) DeleteUser(ctx context.Context, id manta.ID) error {
	panic("implement me")
}

func (s *Service) deleteUser(ctx context.Context, tx Tx, id manta.ID) error {
	panic("implement me")
}
