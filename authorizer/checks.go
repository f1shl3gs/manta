package authorizer

import (
	"context"
	"github.com/f1shl3gs/manta"
)

type CheckService struct {
	cs manta.CheckService
}

// FindCheckByID returns a check by id
func (s *CheckService) FindCheckByID(ctx context.Context, id manta.ID) (*manta.Check, error) {
	panic("not implement")
}

// FindChecks returns a list of checks that match the filter and total count of matching checks
// Additional options provide pagination & sorting.
func (s *CheckService) FindChecks(ctx context.Context, filter manta.CheckFilter, opt ...manta.FindOptions) ([]*manta.Check, int, error) {
	panic("not implement")
}

// CreateCheck creates a new and set its id with new identifier
func (s *CheckService) CreateCheck(ctx context.Context, c *manta.Check) error {
	panic("not implement")
}

// UpdateCheck updates the whole check
// Returns the new check after update
func (s *CheckService) UpdateCheck(ctx context.Context, id manta.ID, c *manta.Check) (*manta.Check, error) {
	panic("not implement")
}

// PatchCheck updates a single check with changeset
// Returns the new check after patch
func (s *CheckService) PatchCheck(ctx context.Context, id manta.ID, u manta.CheckUpdate) (*manta.Check, error) {
	panic("not implement")
}

// DeleteCheck delete a single check by ID
func (s *CheckService) DeleteCheck(ctx context.Context, id manta.ID) error {
	panic("not implement")
}
