package manta

import (
	"context"
	"errors"
)

type VariableFilter struct {
	ID    *ID
	OrgID *ID
}

// VariableUpdate describes a set of changes that can be applied to a Variable
type VariableUpdate struct {
	Name *string `json:"name"`
	Desc *string `json:"desc"`
}

func (udp *VariableUpdate) Apply(v *Variable) {
	if udp.Name != nil {
		v.Name = *udp.Name
	}

	if udp.Desc != nil {
		v.Desc = *udp.Desc
	}
}

type VariableService interface {
	// FindVariableByID finds a single variable from the store by its ID
	FindVariableByID(ctx context.Context, id ID) (*Variable, error)

	// FindVariables returns all variables in the store
	FindVariables(ctx context.Context, filter VariableFilter) ([]*Variable, error)

	// CreateVariable creates a new variable and assigns it an ID
	CreateVariable(ctx context.Context, v *Variable) error

	// PatchVariable updates a single variable with a changeset
	PatchVariable(ctx context.Context, id ID, udp *VariableUpdate) (*Variable, error)

	// UpdateVariable replaces a single variable
	UpdateVariable(ctx context.Context, v *Variable) error

	// DeleteVariable removes a variable from the store
	DeleteVariable(ctx context.Context, id ID) error
}

func (m *Variable) Validate() error {
	if m.OrgID == 0 {
		return ErrInvalidOrgID
	}

	if m.Name == "" {
		return errors.New("name is required")
	}

	if m.Type != "static" && m.Type != "query" {
		return errors.New("unexpected variable type")
	}

	return nil
}
