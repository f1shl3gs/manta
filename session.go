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

type Session struct {
	ID          ID           `json:"id"`
	Created     time.Time    `json:"created"`
	LastSeen    time.Time    `json:"lastSeen"`
	ExpiresAt   time.Time    `json:"expiresAt"`
	UserID      ID           `json:"userId"`
	Permissions []Permission `json:"permissions,omitempty"`
}

func (s *Session) Expired() bool {
	return time.Now().After(s.ExpiresAt)
}

type SessionService interface {
	// CreateSession create a new session
	CreateSession(ctx context.Context, uid ID) (*Session, error)

	// FindSession find session by key
	FindSession(ctx context.Context, id ID) (*Session, error)

	// RevokeSession delete the session, if the session does not
	// exist then nothing is done and a nil error is returned.
	RevokeSession(ctx context.Context, id ID) error

	// RenewSession renew the session and update the ExpireAt
	RenewSession(ctx context.Context, id ID, expiration time.Time) error

	// todo: clean up the sessions which will never be used,
	//   whose ExpireAt > Now + TTL
}

func (s *Session) Identifier() ID {
	return s.ID
}

func (s *Session) GetUserID() ID {
	return s.UserID
}

func (s *Session) Kind() string {
	return "session"
}

func (s *Session) PermissionSet() (PermissionSet, error) {
	if s.Expired() {
		return nil, &Error{
			Code: EForbidden,
			Msg:  "session has expired",
		}
	}

	return s.Permissions, nil
}
