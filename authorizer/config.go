package authorizer

import (
	"context"

	"github.com/f1shl3gs/manta"
)

type ConfigService struct {
	service manta.ConfigService
}

var _ manta.ConfigService = &ConfigService{}

func NewConfigService(service manta.ConfigService) *ConfigService {
	return &ConfigService{
		service: service,
	}
}

func (s *ConfigService) CreateConfig(ctx context.Context, conf *manta.Config) error {
	if _, _, err := authorizeCreate(ctx, manta.ConfigsResourceType, conf.OrgID); err != nil {
		return err
	}

	return s.service.CreateConfig(ctx, conf)
}

func (s *ConfigService) FindConfigByID(ctx context.Context, id manta.ID) (*manta.Config, error) {
	conf, err := s.service.FindConfigByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeRead(ctx, manta.ConfigsResourceType, id, conf.OrgID); err != nil {
		return nil, err
	}

	return conf, nil
}

func (s *ConfigService) FindConfigs(ctx context.Context, filter manta.ConfigFilter) ([]*manta.Config, error) {
	confs, err := s.service.FindConfigs(ctx, filter)
	if err != nil {
		return nil, err
	}

	filtered := confs[:0]
	for _, conf := range confs {
		_, _, err := authorizeRead(ctx, manta.ConfigsResourceType, conf.ID, conf.OrgID)
		if err != nil && manta.ErrorCode(err) != manta.EUnauthorized {
			return nil, err
		}

		if manta.ErrorCode(err) == manta.EUnauthorized {
			continue
		}

		filtered = append(filtered, conf)
	}

	return filtered, nil
}

func (s *ConfigService) UpdateConfig(ctx context.Context, id manta.ID, upd manta.ConfigUpdate) (*manta.Config, error) {
	conf, err := s.service.FindConfigByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeWrite(ctx, manta.ConfigsResourceType, id, conf.OrgID); err != nil {
		return nil, err
	}

	return s.service.UpdateConfig(ctx, id, upd)
}

func (s *ConfigService) DeleteConfig(ctx context.Context, id manta.ID) error {
	conf, err := s.service.FindConfigByID(ctx, id)
	if err != nil {
		return err
	}

	if _, _, err := authorizeWrite(ctx, manta.ConfigsResourceType, id, conf.OrgID); err != nil {
		return err
	}

	return s.service.DeleteConfig(ctx, id)
}
