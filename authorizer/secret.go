package authorizer

import (
	"context"

	"github.com/f1shl3gs/manta"
)

type SecretService struct {
	service manta.SecretService
}

var _ manta.SecretService = &SecretService{}

func NewSecretService(service manta.SecretService) *SecretService {
	return &SecretService{
		service: service,
	}
}

// LoadSecret retrieves the secret value v found at key k for organization orgID
func (s *SecretService) LoadSecret(ctx context.Context, orgID manta.ID, k string) (*manta.Secret, error) {
	if _, _, err := authorizeOrgReadResource(ctx, manta.SecretsResourceType, orgID); err != nil {
		return nil, err
	}

	return s.service.LoadSecret(ctx, orgID, k)
}

// GetSecrets retrieves desensitized secrets of 'orgID'
func (s *SecretService) GetSecrets(ctx context.Context, orgID manta.ID) ([]manta.Secret, error) {
	if _, _, err := authorizeOrgReadResource(ctx, manta.SecretsResourceType, orgID); err != nil {
		return nil, err
	}

	return s.service.GetSecrets(ctx, orgID)
}

// PutSecret creates or updates a secret and return the desensitized secret
func (s *SecretService) PutSecret(ctx context.Context, secret *manta.Secret) (*manta.Secret, error) {
	if _, _, err := authorizeCreate(ctx, manta.SecretsResourceType, secret.OrgID); err != nil {
		return nil, err
	}

	return s.service.PutSecret(ctx, secret)
}

// DeleteSecret deletes secrets by keys
func (s *SecretService) DeleteSecret(ctx context.Context, orgID manta.ID, keys ...string) error {
	if _, _, err := authorizeOrgWriteResource(ctx, manta.SecretsResourceType, orgID); err != nil {
		return err
	}

	return s.service.DeleteSecret(ctx, orgID, keys...)
}
