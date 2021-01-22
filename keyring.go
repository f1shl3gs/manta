package manta

import "context"

// todo: retention info
type Keyring interface {
	// AddKey add the a new to the set, and
	// it become the primary key
	AddKey(ctx context.Context, key []byte) error

	// PrimaryKey returns the key which is the only key
	// for encryption and the first key tried for decryption
	PrimaryKey(ctx context.Context) ([]byte, error)

	// Keys returns all the keys, primary key is included
	Keys(ctx context.Context) ([][]byte, error)

	// RemoveKey drops the key from the keyring
	RemoveKey(ctx context.Context, key []byte) error
}
