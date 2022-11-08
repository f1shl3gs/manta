package manta

import (
	"context"
	"encoding/json"
	"time"
)

type Configuration struct {
	ID      ID        `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	OrgID ID     `json:"orgID"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Data  []byte `json:"data"`
}

type ConfigurationUpdate struct {
	Name *string `json:"name,omitempty"`
	Desc *string `json:"desc,omitempty"`
	Data *[]byte `json:"data,omitempty"`
}

func (upd *ConfigurationUpdate) Apply(cf *Configuration) {
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

type ConfigurationFilter struct {
	OrgID ID
}

type ConfigurationService interface {
	CreateConfiguration(ctx context.Context, cf Configuration) error

	GetConfiguration(ctx context.Context, id ID) (*Configuration, error)

	FindConfigurations(ctx context.Context, filter ConfigurationFilter) ([]*Configuration, error)

	UpdateConfiguration(ctx context.Context, id ID, upd ConfigurationUpdate) error

	DeleteConfiguration(ctx context.Context, id ID) error
}

func (cs *Configuration) Marshal() ([]byte, error) {
	return json.Marshal(cs)
}

func (cs *Configuration) Unmarshal(data []byte) error {
	return json.Unmarshal(data, cs)
}
