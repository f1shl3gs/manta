package kv

import (
	"context"
	"encoding/json"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	userBucket          = []byte("users")
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

	return findByID[manta.User](tx, id, userBucket)
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
	if err = json.Unmarshal(val, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) FindUsers(ctx context.Context, filter manta.UserFilter, opts ...manta.FindOptions) ([]*manta.User, error) {
	var (
		users []*manta.User
		err   error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		users, err = s.findUsers(ctx, tx, filter)
		return err
	})

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) findUsers(ctx context.Context, tx Tx, filter manta.UserFilter) ([]*manta.User, error) {
	b, err := tx.Bucket(userBucket)
	if err != nil {
		return nil, err
	}
	c, err := b.ForwardCursor(nil)
	if err != nil {
		return nil, err
	}

	users := make([]*manta.User, 0, 8)
	err = WalkCursor(ctx, c, func(k, v []byte) error {
		user := &manta.User{}
		err = json.Unmarshal(v, user)
		if err != nil {
			return err
		}

		users = append(users, user)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) CreateUser(ctx context.Context, user *manta.User) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.createUser(ctx, tx, user)
	})
}

func (s *Service) createUser(ctx context.Context, tx Tx, user *manta.User) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	// check user name index
	b, err := tx.Bucket(userNameIndexBucket)
	if err != nil {
		return err
	}

	if _, err = b.Get([]byte(user.Name)); err == nil {
		return manta.ErrUserAlreadyExist
	} else if err != ErrKeyNotFound {
		return err
	}

	// initial user
	user.ID = s.idGen.ID()
	now := time.Now()
	user.Created = now
	user.Updated = now

	return s.putUser(ctx, tx, user)
}

func (s *Service) putUser(ctx context.Context, tx Tx, user *manta.User) error {
	// create user
	pk, err := user.ID.Encode()
	if err != nil {
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	b, err := tx.Bucket(userBucket)
	if err != nil {
		return err
	}

	err = b.Put(pk, data)
	if err != nil {
		return err
	}

	// create name index
	fk := []byte(user.Name)
	b, err = tx.Bucket(userNameIndexBucket)
	if err != nil {
		return err
	}

	return b.Put(fk, pk)
}

func (s *Service) UpdateUser(ctx context.Context, id manta.ID, udp manta.UserUpdate) (*manta.User, error) {
	var (
		user *manta.User
		err  error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		user, err = s.updateUser(ctx, tx, id, udp)
		return err
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) updateUser(ctx context.Context, tx Tx, id manta.ID, udp manta.UserUpdate) (*manta.User, error) {
	prev, err := s.findUserByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	user := *prev
	udp.Apply(&user)
	user.Updated = time.Now()

	err = s.putUser(ctx, tx, &user)
	if err != nil {
		return nil, err
	}

	// update user name index
	if prev.Name == user.Name {
		return &user, nil
	}

	b, err := tx.Bucket(userNameIndexBucket)
	if err != nil {
		return nil, err
	}

	if err = b.Delete([]byte(prev.Name)); err != nil {
		return nil, err
	}

	pk, _ := id.Encode()
	if err = b.Put([]byte(user.Name), pk); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteUser(ctx, tx, id)
	})
}

func (s *Service) deleteUser(ctx context.Context, tx Tx, id manta.ID) error {
	user, err := s.findUserByID(ctx, tx, id)
	if err != nil {
		return err
	}

	// delete user
	b, err := tx.Bucket(userBucket)
	if err != nil {
		return err
	}

	pk, err := id.Encode()
	if err != nil {
		return err
	}

	if err = b.Delete(pk); err != nil {
		return err
	}

	// delete user name index
	b, err = tx.Bucket(userNameIndexBucket)
	if err != nil {
		return err
	}

	return b.Delete([]byte(user.Name))
}
