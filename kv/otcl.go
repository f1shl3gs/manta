package kv

import (
	"context"
	"errors"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	otclBucket   = []byte("otcls")
	otclOrgIndex = []byte("otclorgindex")

	otclOrgIndexMapping = NewIndexMapping(otclBucket, otclOrgIndex, func(data []byte) ([]byte, error) {
		var otcl manta.Otcl
		if err := otcl.Unmarshal(data); err != nil {
			return nil, err
		}

		id, _ := otcl.OrgID.Encode()
		return id, nil
	})
)

func (s *Service) CreateOtcl(ctx context.Context, o *manta.Otcl) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	err := s.kv.Update(ctx, func(tx Tx) error {
		return s.createOtcl(ctx, tx, o)
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) createOtcl(ctx context.Context, tx Tx, o *manta.Otcl) error {
	var (
		id  = s.idGen.ID()
		now = time.Now()
	)

	o.ID = id
	o.Created = now
	o.Modified = now

	return s.putOtcl(ctx, tx, o)
}

func (s *Service) putOtcl(ctx context.Context, tx Tx, o *manta.Otcl) error {
	pk, err := o.ID.Encode()
	if err != nil {
		return err
	}

	fk, err := o.OrgID.Encode()
	if err != nil {
		return manta.ErrInvalidOrgID
	}

	// put data
	data, err := o.Marshal()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(otclBucket)
	if err != nil {
		return err
	}

	if err = b.Put(pk, data); err != nil {
		return err
	}

	// save index
	b, err = tx.Bucket(otclOrgIndex)
	if err != nil {
		return err
	}

	if err = b.Put(indexKey(fk, pk), pk); err != nil {
		return err
	}

	return nil
}

func (s *Service) FindOtclByID(ctx context.Context, id manta.ID) (*manta.Otcl, error) {
	var (
		otcl *manta.Otcl
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		otcl, err = s.findOtclByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return otcl, err
}

func (s *Service) findOtclByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Otcl, error) {
	pk, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(otclBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(pk)
	if err != nil {
		if err == ErrKeyNotFound {
			return nil, manta.ErrOtclNotFound
		}

		return nil, err
	}

	otcl := &manta.Otcl{}
	if err = otcl.Unmarshal(data); err != nil {
		return nil, err
	}

	return otcl, nil
}

func (s *Service) FindOtcls(ctx context.Context, filter manta.OtclFilter) ([]*manta.Otcl, error) {
	var (
		otcls []*manta.Otcl
		err   error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.OrgID != nil {
			otcls, err = s.findOtclsByOrg(ctx, tx, *filter.OrgID)
			return err
		}

		return errors.New("OrgID must be specified")
	})

	if err != nil {
		return nil, err
	}

	return otcls, nil
}

func (s *Service) findOtclsByOrg(ctx context.Context, tx Tx, orgID manta.ID) ([]*manta.Otcl, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.LogKV("OrgID", orgID)

	var (
		keys [][]byte
		err  error
	)

	err = s.findResourceByOrg(ctx, tx, orgID, otclOrgIndex, func(k, v []byte) error {
		keys = append(keys, v)
		return nil
	})

	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(otclBucket)
	if err != nil {
		return nil, err
	}

	values, err := b.GetBatch(keys...)
	if err != nil {
		return nil, err
	}

	otcls := make([]*manta.Otcl, 0, len(values))
	for _, val := range values {
		otcl := &manta.Otcl{}
		if err = otcl.Unmarshal(val); err != nil {
			return nil, err
		}

		otcls = append(otcls, otcl)
	}

	return otcls, nil
}

func (s *Service) findResourceByOrg(ctx context.Context, tx Tx, orgID manta.ID, idxBucket []byte, appender func(k, v []byte) error) error {
	fk, err := orgID.Encode()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(idxBucket)
	if err != nil {
		return err
	}

	c, err := b.ForwardCursor(fk, WithCursorPrefix(fk))
	if err != nil {
		return err
	}

	defer c.Close()

	for {
		k, v := c.Next()
		if k == nil {
			break
		}

		if err = appender(k, v); err != nil {
			return err
		}
	}

	return c.Err()
}

func (s *Service) PatchOtcl(ctx context.Context, id manta.ID, u manta.OtclPatch) (*manta.Otcl, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		otcl *manta.Otcl
		err  error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		otcl, err = s.findOtclByID(ctx, tx, id)
		if err != nil {
			return err
		}

		if u.Name != nil {
			otcl.Name = *u.Name
		}

		if u.Description != nil {
			otcl.Desc = *u.Description
		}

		if u.Content != nil {
			otcl.Content = *u.Content
		}

		otcl.Modified = time.Now()

		return s.putOtcl(ctx, tx, otcl)
	})

	if err != nil {
		return nil, err
	}

	return otcl, nil
}

func (s *Service) DeleteOtcl(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteOtcl(ctx, tx, id)
	})
}

func (s *Service) deleteOtcl(ctx context.Context, tx Tx, id manta.ID) error {
	pk, err := id.Encode()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(otclBucket)
	if err != nil {
		return err
	}

	val, err := b.Get(pk)
	if err != nil {
		if err == ErrKeyNotFound {
			return manta.ErrOtclNotFound
		}

		return nil
	}

	// delete otcls
	if err = b.Delete(pk); err != nil {
		return err
	}

	// delete otcls's org index
	var otcl manta.Otcl
	err = otcl.Unmarshal(val)
	if err != nil {
		return err
	}

	fk, err := otcl.OrgID.Encode()
	if err != nil {
		return err
	}

	b, err = tx.Bucket(otclOrgIndex)
	if err != nil {
		return err
	}

	return b.Delete(indexKey(fk, pk))
}
