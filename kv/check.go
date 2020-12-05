package kv

import (
	"context"
	"sort"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/pkg/errors"
)

var (
	checkBucket         = []byte("check")
	checkOrgIndexBucket = []byte("checkorgindex")
)

var _ manta.CheckService = (*Service)(nil)

func (s *Service) FindCheckByID(ctx context.Context, id manta.ID) (*manta.Check, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		check *manta.Check
		err   error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		check, err = s.findCheckByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return check, nil
}

func (s *Service) findCheckByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Check, error) {
	key, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(checkBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	c := &manta.Check{}
	err = c.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Service) FindChecks(ctx context.Context, filter manta.CheckFilter, opt ...manta.FindOptions) ([]*manta.Check, int, error) {
	var (
		checks []*manta.Check
		err    error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.OrgID != nil {
			checks, err = s.findChecksByOrgID(ctx, tx, *filter.OrgID)
			if err != nil {
				return err
			}

			return nil
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return checks, len(checks), nil
}

func (s *Service) findChecksByOrgID(ctx context.Context, tx Tx, orgID manta.ID) ([]*manta.Check, error) {
	prefix, err := orgID.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(checkOrgIndexBucket)
	if err != nil {
		return nil, err
	}

	cur, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix))
	if err != nil {
		return nil, err
	}

	keys := make([][]byte, 0, 24)
	err = WalkCursor(ctx, cur, func(k, v []byte) error {
		keys = append(keys, v)
		return nil
	})
	if err != nil {
		return nil, err
	}

	b, err = tx.Bucket(checkBucket)
	if err != nil {
		return nil, err
	}

	values, err := b.GetBatch(keys...)
	if err != nil {
		return nil, err
	}

	checks := make([]*manta.Check, 0, len(values))
	for i := 0; i < len(values); i++ {
		c := &manta.Check{}
		if err = c.Unmarshal(values[i]); err != nil {
			return nil, err
		}

		checks = append(checks, c)
	}

	return checks, nil
}

func (s *Service) CreateCheck(ctx context.Context, c *manta.Check) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	span.Finish()

	return s.kv.Update(ctx, func(tx Tx) error {
		if err := s.createCheck(ctx, tx, c); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) createCheck(ctx context.Context, tx Tx, c *manta.Check) error {
	c.ID = s.idGen.ID()
	c.Created = time.Now()
	c.Modified = time.Now()

	org, err := s.findOrganizationByID(ctx, tx, c.OrgID)
	if err != nil {
		return errors.Wrap(err, "cannot find org by id")
	}

	if err := s.putCheck(ctx, tx, c); err != nil {
		return err
	}

	task := &manta.Task{
		Annotations: map[string]string{
			"org":        org.ID.String(),
			"org.name":   org.Name,
			"check":      c.ID.String(),
			"check.name": c.Name,
		},
		Type:    "check",
		Status:  c.Status,
		OwnerID: c.ID,
		OrgID:   c.OrgID,
		Cron:    c.Cron,
	}

	if err := s.createTask(ctx, tx, task); err != nil {
		return err
	}

	return nil
}

func (s *Service) putCheck(ctx context.Context, tx Tx, c *manta.Check) error {
	key, err := c.ID.Encode()
	if err != nil {
		return err
	}

	// todo: uniq (org + name) key

	// org index
	fk, err := c.OrgID.Encode()
	if err != nil {
		return err
	}

	refIdx := IndexKey(fk, key)
	b, err := tx.Bucket(checkOrgIndexBucket)
	if err != nil {
		return err
	}

	err = b.Put(refIdx, key)
	if err != nil {
		return err
	}

	// save check
	b, err = tx.Bucket(checkBucket)
	if err != nil {
		return err
	}

	sort.Slice(c.Conditions, func(i, j int) bool {
		iv := manta.SeverityValue[c.Conditions[i].Status]
		jv := manta.SeverityValue[c.Conditions[j].Status]

		return iv < jv
	})

	data, err := c.Marshal()
	if err != nil {
		return err
	}

	return b.Put(key, data)
}

func (s *Service) UpdateCheck(ctx context.Context, id manta.ID, c *manta.Check) (*manta.Check, error) {
	panic("implement me")
}

func (s *Service) PatchCheck(ctx context.Context, id manta.ID, u manta.CheckUpdate) (*manta.Check, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		check *manta.Check
		err   error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		check, err = s.findCheckByID(ctx, tx, id)
		if err != nil {
			return err
		}

		if err := s.deleteCheck(ctx, tx, id); err != nil {
			return err
		}

		if u.Name != nil {
			check.Name = *u.Name
		}

		if u.Description != nil {
			check.Desc = *u.Description
		}

		if u.Status != nil {
			check.Status = *u.Status
		}

		check.Modified = time.Now()

		return s.putCheck(ctx, tx, check)
	})

	if err != nil {
		return nil, err
	}

	return check, nil
}

func (s *Service) DeleteCheck(ctx context.Context, id manta.ID) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteCheck(ctx, tx, id)
	})
}

func (s *Service) deleteCheck(ctx context.Context, tx Tx, id manta.ID) error {
	c, err := s.findCheckByID(ctx, tx, id)
	if err != nil {
		return err
	}

	pk, _ := id.Encode()
	fk, err := c.OrgID.Encode()
	if err != nil {
		return err
	}

	// delete index
	refIdx := IndexKey(fk, pk)
	b, err := tx.Bucket(checkOrgIndexBucket)
	if err != nil {
		return err
	}

	if err := b.Delete(refIdx); err != nil {
		return err
	}

	// delete check
	b, err = tx.Bucket(checkBucket)
	if err != nil {
		return err
	}

	return b.Delete(pk)
}
