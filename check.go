package manta

import (
	"context"
)

type CheckFilter struct {
	OrgID *ID
}

type CheckUpdate struct {
	Name   *string `json:"name"`
	Desc   *string `json:"desc"`
	Status *string `json:"status"`
}

func (upd *CheckUpdate) Validate() error {
	if upd.Name != nil {
		if *upd.Name == "" {
			return &Error{Code: EInvalid, Op: "validate check's Name", Msg: "Name cannot be empty"}
		}
	}

	if upd.Status != nil {
		if *upd.Status != "active" && *upd.Status != "inactive" {
			return &Error{Code: EInvalid, Op: "validate check's Status", Msg: "status is not active nor inactive"}
		}
	}

	return nil
}

func (upd *CheckUpdate) Apply(check *Check) {
	if upd.Name != nil {
		check.Name = *upd.Name
	}

	if upd.Desc != nil {
		check.Desc = *upd.Desc
	}

	if upd.Status != nil {
		check.Status = *upd.Status
	}
}

type CheckService interface {
	// FindCheckByID returns a check by id
	FindCheckByID(ctx context.Context, id ID) (*Check, error)

	// FindChecks returns a list of checks that match the filter and total count of matching checks
	// Additional options provide pagination & sorting.
	FindChecks(ctx context.Context, filter CheckFilter, opt ...FindOptions) ([]*Check, int, error)

	// CreateCheck creates a new and set its id with new identifier
	CreateCheck(ctx context.Context, c *Check) error

	// UpdateCheck updates the whole check
	// Returns the new check after update
	UpdateCheck(ctx context.Context, id ID, c *Check) (*Check, error)

	// PatchCheck updates a single check with changeset
	// Returns the new check after patch
	PatchCheck(ctx context.Context, id ID, u CheckUpdate) (*Check, error)

	// DeleteCheck delete a single check by ID
	DeleteCheck(ctx context.Context, id ID) error
}
