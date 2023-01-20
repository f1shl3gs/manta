package kv

import (
	"context"

	"github.com/f1shl3gs/manta"
)

var (
	// ErrOnboardingNotAllowed occurs when request to onboard comes in and we are not allowing this request
    ErrOnboardingNotAllowed = &manta.Error{
        Code: manta.EConflict,
		Msg:  "onboarding has already been completed",
	}
)

func (s *Service) Onboarded(ctx context.Context) (bool, error) {
	var (
		initialized bool
		err         error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		orgs, err := s.findAllOrganizations(ctx, tx)
		if err != nil {
			return err
		}

		users, err := s.findUsers(ctx, tx, manta.UserFilter{})
		if err != nil {
			return err
		}

		initialized = len(orgs) != 0 && len(users) != 0

		return nil
	})
	if err != nil {
		return false, err
	}

	return initialized, nil
}

func (s *Service) Setup(ctx context.Context, req *manta.OnBoardingRequest) (*manta.OnboardingResult, error) {
	onboarded, err := s.Onboarded(ctx)
	if err != nil {
		return nil, err
	}

	if onboarded {
		return nil, ErrOnboardingNotAllowed
	}

	var (
		org = &manta.Organization{
			Name: req.Organization,
		}

		user = &manta.User{
			Name:   req.Username,
			Status: "",
		}
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		if err = s.createUser(ctx, tx, user); err != nil {
			return err
		}

		if err = s.setPassword(ctx, tx, user.ID, req.Password); err != nil {
			return err
		}

		if org, err = s.createOrganization(ctx, tx, org); err != nil {
			return err
		}

		err = s.createUserResourceMapping(ctx, tx, &manta.UserResourceMapping{
			UserID:       user.ID,
			UserType:     manta.Owner,
			MappingType:  manta.OrgMappingType,
			ResourceType: manta.InstanceResourceType,
			ResourceID:   manta.ID(1), // The instance doesn't have a resource id
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &manta.OnboardingResult{
		User: user,
		Org:  org,
		Auth: nil,
	}, nil
}
