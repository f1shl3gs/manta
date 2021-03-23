package manta

import "context"

// SecretService a service for storing and retrieving secrets
type SecretService interface {
	// FindSecret retrieves the secret value v found at key k for organization orgID
	FindSecret(ctx context.Context, orgID ID, k string) (string, error)

	// PutSecret stores the secret pair(k, v) for the organization orgID
	PutSecret(ctx context.Context, orgID ID, k, v string) error

	// DeleteSecret removes secrets from the secret store
	DeleteSecret(ctx context.Context, orgID ID, keys ...string) error
}
