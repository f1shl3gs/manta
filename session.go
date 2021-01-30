package manta

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
)

type SessionService interface {
	// CreateSession create a new session
	CreateSession(ctx context.Context, uid ID) (*Session, error)

	// FindSession find session by key
	FindSession(ctx context.Context, id ID) (*Session, error)

	// RevokeSession delete the session
	RevokeSession(ctx context.Context, id ID) error

	//
	RenewSession(ctx context.Context, id ID, expiration time.Time) error
}

func (s *Session) Identifier() ID {
	return s.ID
}

func (s *Session) GetUserID() ID {
	return s.UID
}

func (s *Session) Kind() string {
	return "session"
}

func (s *Session) PermissionSet() PermissionSet {
	return s.Permissions
}