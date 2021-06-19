package manta

import (
	"context"
	"errors"
)

type SecretField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *SecretField) Validate() error {
	if len(s.Key) == 0 {
		return errors.New("secret key cannot be empty")
	}

	for i, b := range s.Key {
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || (b >= '0' && b <= '9' && i > 0)) {
			return errors.New("invalid char for key")
		}
	}

	if s.Value == "" {
		return errors.New("secret value cannot be empty")
	}

	return nil
}

// SecretService a service for storing and retrieving secrets
type SecretService interface {
	// GetSecretKeys retrieves all secret keys that are stored for the organization orgID
	GetSecretKeys(ctx context.Context, orgID ID) ([]string, error)

	// FindSecret retrieves the secret value v found at key k for organization orgID
	FindSecret(ctx context.Context, orgID ID, k string) (string, error)

	// PutSecret stores the secret pair(k, v) for the organization orgID
	PutSecret(ctx context.Context, orgID ID, k, v string) error

	// DeleteSecret removes secrets from the secret store
	DeleteSecret(ctx context.Context, orgID ID, keys ...string) error
}
