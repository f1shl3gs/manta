package manta

import (
	"context"

	"github.com/f1shl3gs/manta/token"
)

type OnBoardingRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
}

func (r *OnBoardingRequest) Validate() error {
	if r.Username == "" {
		return &Error{
			Code: EInvalid,
			Msg:  "username cannot be empty",
		}
	}

	if r.Password == "" {
		return &Error{
			Code: EInvalid,
			Msg:  "password cannot be empty",
		}
	}

	if r.Organization == "" {
		return &Error{
			Code: EInvalid,
			Msg:  "organization cannot be empty",
		}
	}

	return nil
}

type OnboardingResult struct {
	User *User          `json:"user"`
	Org  *Organization  `json:"org"`
	Auth *Authorization `json:"auth"`
}

type OnBoardingService interface {
	Onboarded(ctx context.Context) (bool, error)

	Setup(ctx context.Context, req *OnBoardingRequest) (*OnboardingResult, error)
}

type onBoardingService struct {
	userService          UserService
	passwordService      PasswordService
	authorizationService AuthorizationService
	organizationService  OrganizationService

	tokenGen token.Generator
}

func NewOnBoardingService(
	userService UserService,
	passwordService PasswordService,
	authService AuthorizationService,
	orgService OrganizationService,
) OnBoardingService {
	return &onBoardingService{
		userService:          userService,
		passwordService:      passwordService,
		authorizationService: authService,
		organizationService:  orgService,
		tokenGen:             token.NewGenerator(0),
	}
}

func (o *onBoardingService) Onboarded(ctx context.Context) (bool, error) {
	orgs, _, err := o.organizationService.FindOrganizations(ctx, OrganizationFilter{})
	if err != nil {
		return false, err
	}

	return len(orgs) != 0, err
}

// Setup setup the initial organization, user and authorization
// TODO: Setup should be protect by transaction
func (o *onBoardingService) Setup(ctx context.Context, req *OnBoardingRequest) (*OnboardingResult, error) {
	err := o.userService.CreateUser(ctx, &User{
		Name: req.Username,
	})

	if err != nil {
		return nil, err
	}

	user, err := o.userService.FindUser(ctx, UserFilter{
		Name: &req.Username,
	})

	if err != nil {
		return nil, err
	}

	err = o.passwordService.SetPassword(ctx, user.ID, req.Password)
	if err != nil {
		return nil, err
	}

	err = o.organizationService.CreateOrganization(ctx, &Organization{
		Name: req.Organization,
	})

	if err != nil {
		return nil, err
	}

	org, err := o.organizationService.FindOrganization(ctx, OrganizationFilter{Name: &req.Organization})
	if err != nil {
		return nil, err
	}

	tk, err := o.tokenGen.Token()
	if err != nil {
		return nil, err
	}

	err = o.authorizationService.CreateAuthorization(ctx, &Authorization{
		UID:         user.ID,
		Status:      "active",
		Token:       tk,
		Permissions: append(OwnerPermissions(org.ID), MePermissions(user.ID)...),
	})
	if err != nil {
		return nil, err
	}

	return &OnboardingResult{
		User: user,
		Org:  org,
	}, nil
}
