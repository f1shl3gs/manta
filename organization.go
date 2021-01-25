package manta

import (
	"context"
	"errors"
)

var (
	// ErrInvalidOrgID signifies invalid IDs.
	ErrInvalidOrgID = &Error{
		Code: EInvalid,
		Msg:  "invalid Organization ID",
	}

	ErrOrgAlreadyExist = &Error{
		Code: EInvalid,
		Msg:  "Organization already exist",
	}
)

type OrganizationFilter struct {
	Name *string
}

type OrganizationUpdate struct {
	Name        *string
	Description *string
}

type OrganizationService interface {
	FindOrganizationByID(ctx context.Context, id ID) (*Organization, error)

	// returns the first Organization that matches filter
	FindOrganization(ctx context.Context, filter OrganizationFilter) (*Organization, error)

	// returns a list of Organizations that match filter and the total count of matching Organizations
	// additional options provide pagination & sorting
	FindOrganizations(ctx context.Context, filter OrganizationFilter, opt ...FindOptions) ([]*Organization, int, error)

	// Create a single Organization and sets Organization.id with the new identifier
	CreateOrganization(ctx context.Context, Organization *Organization) error

	// Updates a single Organization with changeset
	// returns the new Organization state after update
	UpdateOrganization(ctx context.Context, id ID, u OrganizationUpdate) (*Organization, error)

	// Remove a Organization by ID
	DeleteOrganization(ctx context.Context, id ID) error
}

func (m *Organization) Validate() error {
	if m.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
