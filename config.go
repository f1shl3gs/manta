package manta

import (
	"context"
	"time"
)

type Config struct {
	ID      ID        `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	OrgID ID     `json:"orgID"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Data  string `json:"data"`
}

func (c *Config) GetID() ID {
	return c.ID
}

func (c *Config) GetOrgID() ID {
	return c.OrgID
}

type ConfigUpdate struct {
	Name *string `json:"name,omitempty"`
	Desc *string `json:"desc,omitempty"`
	Data *string `json:"data,omitempty"`
}

func (upd *ConfigUpdate) Apply(cf *Config) {
	if upd.Name != nil {
		cf.Name = *upd.Name
	}

	if upd.Desc != nil {
		cf.Desc = *upd.Desc
	}

	if upd.Data != nil {
		cf.Data = *upd.Data
	}
}

type ConfigFilter struct {
	OrgID ID
}

type ConfigService interface {
	CreateConfig(ctx context.Context, conf *Config) error

	FindConfigByID(ctx context.Context, id ID) (*Config, error)

	FindConfigs(ctx context.Context, filter ConfigFilter) ([]*Config, error)

	UpdateConfig(ctx context.Context, id ID, upd ConfigUpdate) (*Config, error)

	DeleteConfig(ctx context.Context, id ID) error
}
