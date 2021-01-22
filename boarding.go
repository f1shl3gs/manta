package manta

import "context"

type OnboardRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Org      string `json:"org"`
}

type BoardingService interface {
	Onboarded(ctx context.Context) (bool, error)

	Setup(ctx context.Context, req *OnboardRequest) error
}

type OnBoarding struct {
	UserService
	OrganizationService
}

func (o *OnBoarding) Onboarded(ctx context.Context) (bool, error) {
	orgs, _, err := o.OrganizationService.FindOrganizations(ctx, OrganizationFilter{})
	if err != nil {
		return false, err
	}

	return len(orgs) != 0, err
}

func (o *OnBoarding) Setup(ctx context.Context, req *OnboardRequest) error {
	panic("implement me")
}
