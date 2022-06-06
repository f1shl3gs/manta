package kv

import (
	"context"
	"github.com/f1shl3gs/manta"
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
	var (
		org = &manta.Organization{
			Name: req.Organization,
		}

		user = &manta.User{
			Name:   req.Username,
			Status: "",
		}

		err error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		if err = s.createOrganization(ctx, tx, org); err != nil {
			return err
		}

		if err = s.createUser(ctx, tx, user); err != nil {
			return err
		}

		if err = s.setPassword(ctx, tx, user.ID, req.Password); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &manta.OnboardingResult{
		Username: "",
		Org:      nil,
		Auth:     nil,
	}, nil
}
