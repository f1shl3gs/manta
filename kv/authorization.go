package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	authorizationBucket           = []byte("authorization")
	authorizationTokenIndexBucket = []byte("authorizationtokenindex")

	// todo
	authorizationUserIndexBucket = []byte("authorizationuserindex")
)

var _ manta.AuthorizationService = (*Service)(nil)

func authTokenIndexKey(token string) []byte {
	return []byte(token)
}

func (s *Service) FindAuthorizationByID(ctx context.Context, id manta.ID) (*manta.Authorization, error) {
	var (
		auth *manta.Authorization
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		auth, err = s.findAuthorizationByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return auth, err
}

func (s *Service) findAuthorizationByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Authorization, error) {
	key, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(authorizationBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	auth := &manta.Authorization{}
	if err = auth.Unmarshal(data); err != nil {
		return nil, err
	}

	return auth, nil
}

func (s *Service) FindAuthorizationByToken(ctx context.Context, token string) (*manta.Authorization, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		auth *manta.Authorization
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		auth, err = s.findAuthorizationByToken(ctx, tx, token)
		return err
	})

	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (s *Service) findAuthorizationByToken(ctx context.Context, tx Tx, token string) (*manta.Authorization, error) {
	key := authTokenIndexKey(token)

	b, err := tx.Bucket(authorizationTokenIndexBucket)
	if err != nil {
		return nil, err
	}

	pk, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	b, err = tx.Bucket(authorizationBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	a := &manta.Authorization{}
	if err = a.Unmarshal(data); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Service) CreateAuthorization(ctx context.Context, a *manta.Authorization) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.createAuthorization(ctx, tx, a)
	})
}

func (s *Service) createAuthorization(ctx context.Context, tx Tx, a *manta.Authorization) error {
	var err error

	a.Token, err = s.tokenGen.Token()
	if err != nil {
		return err
	}

	a.ID = s.idGen.ID()
	a.Created = time.Now()
	a.Modified = time.Now()

	return s.putAuthorization(ctx, tx, a)
}

func (s *Service) putAuthorization(ctx context.Context, tx Tx, a *manta.Authorization) error {
	pk, err := a.ID.Encode()
	if err != nil {
		return err
	}

	// token index
	idx := []byte(a.Token)
	b, err := tx.Bucket(authorizationTokenIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Put(idx, pk); err != nil {
		return err
	}

	// user index
	fk, err := a.UID.Encode()
	if err != nil {
		return err
	}

	b, err = tx.Bucket(authorizationUserIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Put(fk, pk); err != nil {
		return err
	}

	// save auth
	b, err = tx.Bucket(authorizationBucket)
	if err != nil {
		return err
	}

	data, err := a.Marshal()
	if err != nil {
		return err
	}

	return b.Put(pk, data)
}

func (s *Service) UpdateAuthorization(ctx context.Context, id manta.ID, u manta.UpdateAuthorization) (*manta.Authorization, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		a   *manta.Authorization
		err error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		a, err = s.updateAuthorization(ctx, tx, id, u)
		return err
	})

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Service) updateAuthorization(ctx context.Context, tx Tx, id manta.ID, u manta.UpdateAuthorization) (*manta.Authorization, error) {
	a, err := s.findAuthorizationByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if err = s.deleteAuthorization(ctx, tx, id); err != nil {
		return nil, err
	}

	if u.Token != nil {
		a.Token = *u.Token
	}

	if u.Status != nil {
		a.Status = *u.Status
	}

	a.Modified = time.Now()

	if err = s.putAuthorization(ctx, tx, a); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Service) DeleteAuthorization(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteAuthorization(ctx, tx, id)
	})
}

func (s *Service) deleteAuthorization(ctx context.Context, tx Tx, id manta.ID) error {
	a, err := s.findAuthorizationByID(ctx, tx, id)
	if err != nil {
		return err
	}

	// delete token index
	tk := authTokenIndexKey(a.Token)
	b, err := tx.Bucket(authorizationTokenIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Delete(tk); err != nil {
		return err
	}

	// delete authorization
	pk, _ := id.Encode()
	b, err = tx.Bucket(authorizationBucket)
	if err != nil {
		return err
	}

	return b.Delete(pk)
}
