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

	ErrOrgNotFound = &Error{
		Code: ENotFound,
		Msg:  "Organization not found",
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

	// FindOrganization returns the first Organization that matches filter
	FindOrganization(ctx context.Context, filter OrganizationFilter) (*Organization, error)

	// FindOrganizations returns a list of Organizations that match filter and the total count of matching Organizations
	// additional options provide pagination & sorting
	FindOrganizations(ctx context.Context, filter OrganizationFilter, opt ...FindOptions) ([]*Organization, int, error)

	// CreateOrganization create a single Organization and sets Organization.id with the new identifier
	CreateOrganization(ctx context.Context, Organization *Organization) error

	// UpdateOrganization updates a single Organization with changeset
	// returns the new Organization state after update
	UpdateOrganization(ctx context.Context, id ID, u OrganizationUpdate) (*Organization, error)

	// DeleteOrganization remove a Organization by ID
	DeleteOrganization(ctx context.Context, id ID) error
}

func (m *Organization) Validate() error {
	if m.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
