package manta

import (
	"context"
	"fmt"
	"github.com/f1shl3gs/manta/pkg/slices"
)

type DatasourceUpdate struct {
	Name *string
	Desc *string
	Type *string
	URL  *string
}

type DatasourceFilter struct {
	Default *bool
	Name    *string
	Type    *string
}

type DatasourceService interface {
	FindDatasourceByID(ctx context.Context, id ID) (*Datasource, error)

	FindDatasource(ctx context.Context, filter DatasourceFilter) (*Datasource, error)

	FindDatasources(ctx context.Context, filter DatasourceFilter) ([]*Datasource, error)

	CreateDatasource(ctx context.Context, ds *Datasource) error

	UpdateDatasource(ctx context.Context, id ID, udp DatasourceUpdate) (*Datasource, error)

	DeleteDatasource(ctx context.Context, id ID) error
}

var (
	SupportDatasourceTypes = []string{"prometheus"}
)

type DatasourceType string

const (
	DatasourcePrometheus DatasourceType = "prometheus"
)

func (m *Datasource) Validate() error {
	if m.Name == "" {
		return invalidField("name", ErrFieldMustBeSet)
	}

	if !slices.Contain(SupportDatasourceTypes, m.Type) {
		return invalidField("type", fmt.Errorf("datasource is not support"))
	}

	return nil
}
