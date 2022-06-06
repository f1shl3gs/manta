package authz

import (
	"context"
	"github.com/f1shl3gs/manta"
)

type OrganizationService struct {
	s manta.OrganizationService
}

func (o *OrganizationService) FindOrganizationByID(ctx context.Context, id manta.ID) (*manta.Organization, error) {
	panic("implement me")
}

func (o *OrganizationService) FindOrganization(ctx context.Context, filter manta.OrganizationFilter) (*manta.Organization, error) {
	panic("implement me")
}

func (o *OrganizationService) FindOrganizations(ctx context.Context, filter manta.OrganizationFilter, opt ...manta.FindOptions) ([]*manta.Organization, int, error) {
	panic("implement me")
}

func (o *OrganizationService) CreateOrganization(ctx context.Context, Organization *manta.Organization) error {
	panic("implement me")
}

func (o *OrganizationService) UpdateOrganization(ctx context.Context, id manta.ID, u manta.OrganizationUpdate) (*manta.Organization, error) {
	panic("implement me")
}

func (o *OrganizationService) DeleteOrganization(ctx context.Context, id manta.ID) error {
	panic("implement me")
}
