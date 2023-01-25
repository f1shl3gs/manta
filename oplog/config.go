package oplog

import (
	"context"
	"encoding/json"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorizer"

	"go.uber.org/zap"
)

type ConfigService struct {
	manta.ConfigService

	logger *zap.Logger
	oplog  manta.OperationLogService
}

func NewConfigService(service manta.ConfigService, oplog manta.OperationLogService, logger *zap.Logger) *ConfigService {
	return &ConfigService{
		ConfigService: service,
		logger:        logger,
		oplog:         oplog,
	}
}

func (s *ConfigService) CreateConfig(ctx context.Context, conf *manta.Config) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.ConfigService.CreateConfig(ctx, conf)
	if err != nil {
		return err
	}

	conf, err = s.ConfigService.FindConfigByID(ctx, conf.ID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Create,
		ResourceID:   conf.ID,
		ResourceType: manta.ConfigsResourceType,
		OrgID:        conf.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add create config oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", conf.ID),
			zap.Stringer("orgID", conf.OrgID))
	}
	return err
}

func (s *ConfigService) UpdateConfig(ctx context.Context, id manta.ID, upd manta.ConfigUpdate) (*manta.Config, error) {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	conf, err := s.ConfigService.UpdateConfig(ctx, id, upd)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Create,
		ResourceID:   conf.ID,
		ResourceType: manta.ConfigsResourceType,
		OrgID:        conf.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add update config oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", conf.ID),
			zap.Stringer("orgID", conf.OrgID))
		return nil, err
	}

	return conf, nil
}

func (s *ConfigService) DeleteConfig(ctx context.Context, id manta.ID) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	conf, err := s.ConfigService.FindConfigByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.ConfigService.DeleteConfig(ctx, id)
	if err != nil {
		return err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Create,
		ResourceID:   conf.ID,
		ResourceType: manta.ConfigsResourceType,
		OrgID:        conf.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: nil,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add delete config oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", conf.ID),
			zap.Stringer("orgID", conf.OrgID))
	}

	return err
}
