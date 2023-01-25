package oplog

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorizer"

	"go.uber.org/zap"
)

type SecretService struct {
	manta.SecretService

	logger *zap.Logger
	oplog  manta.OperationLogService
}

func NewSecretService(service manta.SecretService, oplog manta.OperationLogService, logger *zap.Logger) *SecretService {
	return &SecretService{
		SecretService: service,
		logger:        logger,
		oplog:         oplog,
	}
}

// PutSecret creates or updates a secret and return the desensitized secret
func (s *SecretService) PutSecret(ctx context.Context, secret *manta.Secret) (*manta.Secret, error) {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	secret, err = s.SecretService.PutSecret(ctx, secret)
	if err != nil {
		return nil, err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Put,
		ResourceID:   manta.UniqueKeyToID(secret.Key),
		ResourceType: manta.SecretsResourceType,
		OrgID:        secret.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: nil,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add put secret oplog failed",
			zap.Error(err),
			zap.Stringer("orgID", secret.OrgID),
			zap.String("key", secret.Key))
		return nil, err
	}

	return secret, nil
}

// DeleteSecret deletes secrets by keys
func (s *SecretService) DeleteSecret(ctx context.Context, orgID manta.ID, keys ...string) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.DeleteSecret(ctx, orgID, keys...)
	if err != nil {
		return err
	}

	for _, key := range keys {
		// TODO: optimize with batch?
		err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
			Type:         manta.Delete,
			ResourceID:   manta.UniqueKeyToID(key),
			ResourceType: manta.SecretsResourceType,
			OrgID:        orgID,
			UserID:       auth.GetUserID(),
			ResourceBody: nil,
			Time:         now,
		})
		if err != nil {
			s.logger.Error("add delete secret oplog failed",
				zap.Error(err),
				zap.Stringer("orgID", orgID),
				zap.String("key", key))
			return err
		}
	}

	return nil
}
