package manta

import (
	"context"
	"errors"

	"github.com/f1shl3gs/manta/token"
)

var (
	ErrInvalidField = errors.New("invalid failed")
)

type OnBoardingRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Org      string `json:"org"`
}

func (r *OnBoardingRequest) Validate() error {
	if r.Username == "" {
		return &Error{
			Code: EInvalid,
			Msg:  "username cannot be empty",
			Op:   "validate",
			Err:  ErrInvalidField,
		}
	}

	if r.Password == "" {
		return &Error{
			Code: EInvalid,
			Msg:  "password cannot be empty",
			Op:   "validate",
			Err:  ErrInvalidField,
		}
	}

	if r.Org == "" {
		return &Error{
			Code: EInvalid,
			Msg:  "org cannot be empty",
			Op:   "validate",
			Err:  ErrInvalidField,
		}
	}

	return nil
}

type OnBoardingService interface {
	Onboarded(ctx context.Context) (bool, error)

	Setup(ctx context.Context, req *OnBoardingRequest) error
}

type onBoardingService struct {
	userService          UserService
	passwordService      PasswordService
	authorizationService AuthorizationService
	organizationService  OrganizationService

	tokenGen TokenGenerator
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
func (o *onBoardingService) Setup(ctx context.Context, req *OnBoardingRequest) error {
	err := o.userService.CreateUser(ctx, &User{
		Name: req.Username,
	})

	if err != nil {
		return err
	}

	user, err := o.userService.FindUser(ctx, UserFilter{
		Name: &req.Username,
	})

	if err != nil {
		return err
	}

	err = o.passwordService.SetPassword(ctx, user.ID, req.Password)
	if err != nil {
		return err
	}

	err = o.organizationService.CreateOrganization(ctx, &Organization{
		Name: req.Org,
	})

	if err != nil {
		return err
	}

	org, err := o.organizationService.FindOrganization(ctx, OrganizationFilter{Name: &req.Org})
	if err != nil {
		return err
	}

	token, err := o.tokenGen.Token()
	if err != nil {
		return err
	}

	return o.authorizationService.CreateAuthorization(ctx, &Authorization{
		UID:         user.ID,
		Status:      "active",
		Token:       token,
		Permissions: append(OwnerPermissions(org.ID), MePermissions(user.ID)...),
	})
}