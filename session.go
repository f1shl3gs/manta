package manta

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
)

type Session struct {
	ID          ID           `json:"id,omitempty"`
	Created     time.Time    `json:"created"`
	ExpiresAt   time.Time    `json:"expiresAt"`
	UID         ID           `json:"userId,omitempty"`
	Permissions []Permission `json:"permissions"`
}

func (s *Session) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *Session) Marshal() ([]byte, error) {
	return json.Marshal(s)
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
	return s.UID
}

func (s *Session) Kind() string {
	return "session"
}

func (s *Session) PermissionSet() PermissionSet {
	return s.Permissions
}
