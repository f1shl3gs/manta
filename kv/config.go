package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
)

var (
	ConfigBucket = []byte("configs")

	ConfigOrgIndexBucket = []byte("configorgindex")
)

func (s *Service) CreateConfig(ctx context.Context, cf *manta.Config) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.createConfig(ctx, tx, cf)
	})
}

func (s *Service) createConfig(ctx context.Context, tx Tx, cf *manta.Config) error {
	now := time.Now()
	cf.ID = s.idGen.ID()
	cf.Created = now
	cf.Updated = now

	return putOrgIndexed(tx, cf, ConfigBucket, ConfigOrgIndexBucket)
}

func (s *Service) FindConfigByID(ctx context.Context, id manta.ID) (*manta.Config, error) {
	var (
		cf  *manta.Config
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		cf, err = findByID[manta.Config](tx, id, ConfigBucket)
		return err
	})

	return cf, err
}

func (s *Service) FindConfigs(
	ctx context.Context,
	filter manta.ConfigFilter,
) ([]*manta.Config, error) {
	var (
		cs  []*manta.Config
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		cs, err = findOrgIndexed[manta.Config](ctx, tx, filter.OrgID, ConfigBucket, ConfigOrgIndexBucket)
		return err
	})

	return cs, err
}

func (s *Service) UpdateConfig(
	ctx context.Context,
	id manta.ID,
	upd manta.ConfigUpdate,
) (*manta.Config, error) {
	var (
		cf  *manta.Config
		err error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		cf, err = findByID[manta.Config](tx, id, ConfigBucket)
		if err != nil {
			return err
		}

		upd.Apply(cf)
		cf.Updated = time.Now()

		return putOrgIndexed(tx, cf, ConfigBucket, ConfigOrgIndexBucket)
	})
	if err != nil {
		return nil, err
	}

	return cf, nil
}

func (s *Service) DeleteConfig(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteConfig(tx, id)
	})
}

func (s *Service) deleteConfig(tx Tx, id manta.ID) error {
	return deleteOrgIndexed[manta.Config](tx, id, ConfigBucket, ConfigOrgIndexBucket)
}
