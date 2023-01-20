package authorizer

import (
	"context"

	"github.com/f1shl3gs/manta"
)

type CheckService struct {
	service manta.CheckService
}

func NewCheckService(service manta.CheckService) *CheckService {
	return &CheckService{
		service: service,
	}
}

// FindCheckByID returns a check by id
func (s *CheckService) FindCheckByID(ctx context.Context, id manta.ID) (*manta.Check, error) {
	check, err := s.service.FindCheckByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeRead(ctx, manta.ChecksResourceType, id, check.OrgID); err != nil {
		return nil, err
	}

	return check, nil
}

// FindChecks returns a list of checks that match the filter and total count of matching checks
// Additional options provide pagination & sorting.
func (s *CheckService) FindChecks(ctx context.Context, filter manta.CheckFilter, opt ...manta.FindOptions) ([]*manta.Check, int, error) {
	checks, _, err := s.service.FindChecks(ctx, filter, opt...)
	if err != nil {
		return nil, 0, err
	}

	filtered := checks[:0]
	for _, c := range checks {
		_, _, err := authorizeRead(ctx, manta.ChecksResourceType, c.ID, c.OrgID)
		if err != nil && manta.ErrorCode(err) != manta.EUnauthorized {
			return nil, 0, err
		}

		if manta.ErrorCode(err) == manta.EUnauthorized {
			continue
		}

		filtered = append(filtered, c)
	}

	return filtered, len(filtered), nil
}

// CreateCheck creates a new and set its id with new identifier
func (s *CheckService) CreateCheck(ctx context.Context, c *manta.Check) error {
	if _, _, err := authorizeCreate(ctx, manta.ChecksResourceType, c.OrgID); err != nil {
		return err
	}

	return s.service.CreateCheck(ctx, c)
}

// UpdateCheck updates the whole check
// Returns the new check after update
func (s *CheckService) UpdateCheck(ctx context.Context, id manta.ID, c *manta.Check) (*manta.Check, error) {
	if _, _, err := authorizeWrite(ctx, manta.ChecksResourceType, c.OrgID, id); err != nil {
		return nil, err
	}

	return s.service.UpdateCheck(ctx, id, c)
}

// PatchCheck updates a single check with changeset
// Returns the new check after patch
func (s *CheckService) PatchCheck(ctx context.Context, id manta.ID, u manta.CheckUpdate) (*manta.Check, error) {
	check, err := s.service.FindCheckByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeWrite(ctx, manta.ChecksResourceType, check.OrgID, id); err != nil {
		return nil, err
	}

	return s.service.PatchCheck(ctx, id, u)
}

// DeleteCheck delete a single check by ID
func (s *CheckService) DeleteCheck(ctx context.Context, id manta.ID) error {
	check, err := s.service.FindCheckByID(ctx, id)
	if err != nil {
		return err
	}

	if _, _, err := authorizeWrite(ctx, manta.ChecksResourceType, check.OrgID, id); err != nil {
		return err
	}

	return s.service.DeleteCheck(ctx, id)
}
