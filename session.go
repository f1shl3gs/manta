package manta

import (
	"context"
	baseErrors "errors"
	"time"

	"github.com/f1shl3gs/manta/errors"
)

var (
	ErrSessionNotFound = baseErrors.New("session not found")
	ErrSessionExpired  = baseErrors.New("session has expired")
)

type Session struct {
	ID          ID           `json:"id,omitempty"`
	Created     time.Time    `json:"created"`
	ExpiresAt   time.Time    `json:"expiresAt"`
	UserID      ID           `json:"userId,omitempty"`
	Permissions []Permission `json:"permissions"`
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
		return nil, &errors.Error{
			Code: errors.EForbidden,
			Msg:  "session has expired",
		}
	}

	return s.Permissions, nil
}
