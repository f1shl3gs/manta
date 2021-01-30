package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
)

var (
	sessionBucket = []byte("session")
)

func (s *Service) CreateSession(ctx context.Context, uid manta.ID) (*manta.Session, error) {
	var (
		session *manta.Session
		err     error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		session, err = s.createSession(ctx, tx, uid)
		return err
	})

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) createSession(ctx context.Context, tx Tx, uid manta.ID) (*manta.Session, error) {
	_, err := s.findUserByID(ctx, tx, uid)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session := &manta.Session{
		ID:        s.idGen.ID(),
		Created:   now,
		ExpiresAt: now.Add(time.Hour),
		UID:       uid,
	}

	if err := s.putSession(ctx, tx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) putSession(ctx context.Context, tx Tx, session *manta.Session) error {
	data, err := session.Marshal()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(sessionBucket)
	if err != nil {
		return err
	}

	pk, err := session.ID.Encode()
	if err != nil {
		return err
	}

	return b.Put(pk, data)
}

func (s *Service) FindSession(ctx context.Context, id manta.ID) (*manta.Session, error) {
	var (
		session *manta.Session
		err     error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		session, err = s.findSession(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) findSession(ctx context.Context, tx Tx, id manta.ID) (*manta.Session, error) {
	pk, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(sessionBucket)
	if err != nil {
		return nil, err
	}

	val, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	session := &manta.Session{}
	err = session.Unmarshal(val)
	if err != nil {
		return nil, err
	}

	as, err := s.findAuthorizationsByUser(ctx, tx, session.UID)
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		session.Permissions = append(session.Permissions, a.Permissions...)
	}

	return session, nil
}

func (s *Service) RevokeSession(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteSession(ctx, tx, id)
	})
}

func (s *Service) deleteSession(ctx context.Context, tx Tx, id manta.ID) error {
	pk, err := id.Encode()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(sessionBucket)
	if err != nil {
		return err
	}

	return b.Delete(pk)
}

func (s *Service) RenewSession(ctx context.Context, id manta.ID, expiration time.Time) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		session, err := s.findSession(ctx, tx, id)
		if err != nil {
			return err
		}

		session.ExpiresAt = expiration

		return s.putSession(ctx, tx, session)
	})
}