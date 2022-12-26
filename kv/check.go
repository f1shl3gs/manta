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
		ck, err = findByID[manta.Check](tx, id, ChecksBucket)
		return err
	})

	return ck, err
}

// FindChecks returns a list of checks that match the filter and total count of matching checks
// Additional options provide pagination & sorting.
func (s *Service) FindChecks(ctx context.Context, filter manta.CheckFilter, opt ...manta.FindOptions) ([]*manta.Check, int, error) {
	var (
		list []*manta.Check
		err  error
	)

	if filter.OrgID == nil {
		return nil, 0, ErrOrgIDRequired
	}

	err = s.kv.View(ctx, func(tx Tx) error {
		list, err = findOrgIndexed[manta.Check](ctx, tx, *filter.OrgID, ChecksBucket, CheckOrgIndexBucket)
		return err
	})
	if err != nil {
		return nil, 0, err
	}

	return list, len(list), nil
}

// CreateCheck creates a new and set its id with new identifier
func (s *Service) CreateCheck(ctx context.Context, check *manta.Check) error {
	now := time.Now()

	check.ID = s.idGen.ID()
	check.Created = now
	check.Updated = now

	task := &manta.Task{
		ID:      s.idGen.ID(),
		Created: now,
		Updated: now,
		Type:    "check",
		Status:  manta.TaskStatus(check.Status),
		OwnerID: check.ID,
		OrgID:   check.OrgID,
		Cron:    check.Cron,
	}

	check.TaskID = task.ID

	return s.kv.Update(ctx, func(tx Tx) error {
		err := putOrgIndexed(tx, task, TasksBucket, TaskOrgIndexBucket)
		if err != nil {
			return err
		}

		return putOrgIndexed(tx, check, ChecksBucket, CheckOrgIndexBucket)
	})
}

// UpdateCheck updates the whole check returns the new check after update
func (s *Service) UpdateCheck(ctx context.Context, id manta.ID, check *manta.Check) (*manta.Check, error) {
	err := s.kv.Update(ctx, func(tx Tx) error {
		err := deleteOrgIndexed[manta.Check](tx, id, ChecksBucket, CheckOrgIndexBucket)
		if err != nil {
			return err
		}

		check.Updated = time.Now()

		err = putOrgIndexed(tx, check, ChecksBucket, CheckOrgIndexBucket)
		if err != nil {
			return err
		}

		// patch task
		status := manta.TaskStatus(check.Status)
		taskUpd := manta.TaskUpdate{
			Desc:   &check.Desc,
			Status: &status,
			Cron:   &check.Cron,
		}

		_, err = updateTask(tx, check.TaskID, taskUpd)
		return err
	})

	if err != nil {
		return nil, err
	}

	return check, nil
}

// PatchCheck updates a single check with changeset
// Returns the new check after patch
func (s *Service) PatchCheck(ctx context.Context, id manta.ID, upd manta.CheckUpdate) (*manta.Check, error) {
	var (
		check *manta.Check
		err   error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		check, err = findByID[manta.Check](tx, id, ChecksBucket)
		if err != nil {
			return err
		}

		upd.Apply(check)
		check.Updated = time.Now()

		err = putOrgIndexed(tx, check, ChecksBucket, CheckOrgIndexBucket)
		if err != nil {
			return err
		}

		// patch task
		taskUpd := manta.TaskUpdate{
			Desc: upd.Desc,
		}

		if upd.Status != nil {
			status := manta.TaskStatus(*upd.Status)
			taskUpd.Status = &status
		}

		_, err = updateTask(tx, check.TaskID, taskUpd)
		return err
	})

	return check, err
}

// DeleteCheck delete a single check by ID
func (s *Service) DeleteCheck(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		check, err := findByID[manta.Check](tx, id, ChecksBucket)
		if err != nil {
			return err
		}

		if err = deleteTask(tx, check.TaskID); err != nil {
			return err
		}

		runs, _, err := findRuns(tx, manta.RunFilter{Task: check.TaskID})
		if err != nil {
			return err
		}

		for _, run := range runs {
			err = deleteRun(tx, check.TaskID, run.ID)
			if err != nil {
				return err
			}
		}

		return deleteOrgIndexed[manta.Check](tx, id, ChecksBucket, CheckOrgIndexBucket)
	})
}
