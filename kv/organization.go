package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	organizationBucket          = []byte("organizations")
	organizationNameIndexBucket = []byte("organizationnameindex")
)

func orgNameIndexKey(n string) []byte {
	return []byte(n)
}

func (s *Service) FindOrganizationByID(ctx context.Context, id manta.ID) (*manta.Organization, error) {
	var org *manta.Organization

	err := s.kv.View(ctx, func(tx Tx) error {
		tmp, err := s.findOrganizationByID(ctx, tx, id)
		if err != nil {
			return err
		}

		org = tmp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return org, nil
}

func (s *Service) findOrganizationByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Organization, error) {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	b, err := tx.Bucket(organizationBucket)
	if err != nil {
		return nil, err
	}

	key, err := id.Encode()
	if err != nil {
		return nil, err
	}

	v, err := b.Get(key)
	if err != nil {
		if err == ErrKeyNotFound {
			return nil, manta.ErrOrgNotFound
		}

		return nil, err
	}

	org := &manta.Organization{}
	if err = org.Unmarshal(v); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *Service) findOrganizationByName(ctx context.Context, tx Tx, n string) (*manta.Organization, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	b, err := tx.Bucket(organizationNameIndexBucket)
	if err != nil {
		return nil, err
	}

	v, err := b.Get(orgNameIndexKey(n))
	if err != nil {
		return nil, err
	}

	var id manta.ID
	err = id.Decode(v)
	if err != nil {
		return nil, err
	}

	o, err := s.findOrganizationByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (s *Service) FindOrganization(ctx context.Context, filter manta.OrganizationFilter) (*manta.Organization, error) {
	var (
		org *manta.Organization
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.Name != nil {
			org, err = s.findOrganizationByName(ctx, tx, *filter.Name)
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return org, nil
}

func (s *Service) FindOrganizations(ctx context.Context, filter manta.OrganizationFilter, opt ...manta.FindOptions) ([]*manta.Organization, int, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		orgs []*manta.Organization
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.Name != nil {
			org, err := s.findOrganizationByName(ctx, tx, *filter.Name)
			if err != nil {
				return err
			}

			orgs = append(orgs, org)
			return nil
		}

		orgs, err = s.findAllOrganizations(ctx, tx)
		return err
	})

	if err != nil {
		return nil, 0, err
	}

	return orgs, len(orgs), nil
}

func (s *Service) findAllOrganizations(ctx context.Context, tx Tx) ([]*manta.Organization, error) {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	b, err := tx.Bucket(organizationBucket)
	if err != nil {
		return nil, err
	}

	c, err := b.Cursor()
	if err != nil {
		return nil, err
	}

	orgs := make([]*manta.Organization, 0)
	for k, v := c.First(); k != nil; k, v = c.Next() {

		o := &manta.Organization{}
		if err := o.Unmarshal(v); err != nil {
			return nil, err
		}

		orgs = append(orgs, o)
	}

	return orgs, nil
}

func (s *Service) CreateOrganization(ctx context.Context, org *manta.Organization) error {
	err := s.kv.Update(ctx, func(tx Tx) error {
		if err := s.createOrganization(ctx, tx, org); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) createOrganization(ctx context.Context, tx Tx, org *manta.Organization) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	org.ID = s.idGen.ID()
	org.Created = time.Now()
	org.Updated = time.Now()

	return s.putOrganization(ctx, tx, org)
}

func (s *Service) putOrganization(ctx context.Context, tx Tx, org *manta.Organization) error {
	pk, err := org.ID.Encode()
	if err != nil {
		return err
	}

	// name index
	fk := []byte(org.Name)
	b, err := tx.Bucket(organizationNameIndexBucket)
	if err != nil {
		return err
	}

	// check name conflict
	if _, err = b.Get(fk); err != ErrKeyNotFound {
		return manta.ErrOrgAlreadyExist
	}

	if err = b.Put(fk, pk); err != nil {
		return err
	}

	// organization
	data, err := org.Marshal()
	if err != nil {
		return err
	}

	b, err = tx.Bucket(organizationBucket)
	if err != nil {
		return err
	}

	return b.Put(pk, data)
}

func (s *Service) UpdateOrganization(ctx context.Context, id manta.ID, u manta.OrganizationUpdate) (*manta.Organization, error) {
	var resp *manta.Organization

	err := s.kv.Update(ctx, func(tx Tx) error {
		current, err := s.findOrganizationByID(ctx, tx, id)
		if err != nil {
			return err
		}

		if err := s.deleteOrganization(ctx, tx, id); err != nil {
			return err
		}

		if u.Name != nil {
			current.Name = *u.Name
		}

		if u.Description != nil {
			current.Desc = *u.Description
		}

		current.Updated = time.Now()
		err = s.createOrganization(ctx, tx, current)
		if err != nil {
			return err
		}

		resp = current

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) DeleteOrganization(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		if err := s.deleteOrganization(ctx, tx, id); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) deleteOrganization(ctx context.Context, tx Tx, id manta.ID) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	org, err := s.findOrganizationByID(ctx, tx, id)
	if err != nil {
		return err
	}

	pk, err := org.ID.Encode()
	if err != nil {
		return err
	}

	// name index
	fk := []byte(org.Name)
	nameIdx := IndexKey(fk, pk)
	b, err := tx.Bucket(organizationNameIndexBucket)
	if err != nil {
		return nil
	}

	return b.Delete(nameIdx)
}
