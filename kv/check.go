package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
)

var (
	ChecksBucket        = []byte("checks")
	CheckOrgIndexBucket = []byte("checkorgindex")
)

// FindCheckByID returns a check by id
func (s *Service) FindCheckByID(ctx context.Context, id manta.ID) (*manta.Check, error) {
	var (
		ck  *manta.Check
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		ck, err = getOrgIndexed[manta.Check](tx, id, CheckOrgIndexBucket)
		return err
	})

	return ck, err
}

// FindChecks returns a list of checks that match the filter and total count of matching checks
// Additional options provide pagination & sorting.
func (s *Service) FindChecks(ctx context.Context, filter manta.CheckFilter, opt ...manta.FindOptions) ([]*manta.Check, int, error) {
	// TODO
	return nil, 0, ErrNotImplement
}

// CreateCheck creates a new and set its id with new identifier
func (s *Service) CreateCheck(ctx context.Context, c *manta.Check) error {
	now := time.Now()

	c.ID = s.idGen.ID()
	c.Created = now
	c.Updated = now

	return s.kv.Update(ctx, func(tx Tx) error {
		return putOrgIndexed(tx, c, ChecksBucket, ConfigurationOrgIndexBucket)
	})
}

// UpdateCheck updates the whole check returns the new check after update
func (s *Service) UpdateCheck(ctx context.Context, id manta.ID, c *manta.Check) (*manta.Check, error) {
	err := s.kv.Update(ctx, func(tx Tx) error {
		err := deleteOrgIndexed[manta.Check](tx, id, ChecksBucket, CheckOrgIndexBucket)
		if err != nil {
			return err
		}

		c.Updated = time.Now()

		return putOrgIndexed(tx, c, ChecksBucket, CheckOrgIndexBucket)
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}

// PatchCheck updates a single check with changeset
// Returns the new check after patch
func (s *Service) PatchCheck(ctx context.Context, id manta.ID, u manta.CheckUpdate) (*manta.Check, error) {
	var (
		c   *manta.Check
		err error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		c, err = getOrgIndexed[manta.Check](tx, id, ChecksBucket)
		if err != nil {
			return err
		}

		u.Apply(c)
		c.Updated = time.Now()

		return putOrgIndexed(tx, c, ChecksBucket, CheckOrgIndexBucket)
	})

	return c, err
}

// DeleteCheck delete a single check by ID
func (s *Service) DeleteCheck(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return deleteOrgIndexed[manta.Check](tx, id, ChecksBucket, CheckOrgIndexBucket)
	})
}
