package manta

import (
	"context"
	"errors"
)

var (
	ErrPasswordNotMatch = errors.New("password not match")
)

type PasswordService interface {
	// SetPassword overrides the password of a known user
	SetPassword(ctx context.Context, uid ID, password string) error

	// ComparePassword checks if the password matches the password stored
	ComparePassword(ctx context.Context, uid ID, password string) error

	// CompareAndSetPassword checks the password and if they match
	// updates to the new password
	CompareAndSetPassword(ctx context.Context, uid ID, old, new string) error

	// DeletePassword delete password by user id
	DeletePassword(ctx context.Context, uid ID) error
}
