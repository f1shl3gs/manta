package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
)

var (
	ConfigurationBucket = []byte("configurations")

	ConfigurationOrgIndexBucket = []byte("configurationorgindex")
)

func (s *Service) CreateConfiguration(ctx context.Context, cf *manta.Configuration) error {
	now := time.Now()
	cf.ID = s.idGen.ID()
	cf.Created = now
	cf.Updated = now

	return s.kv.Update(ctx, func(tx Tx) error {
		return putOrgIndexed(tx, cf, ConfigurationBucket, ConfigurationOrgIndexBucket)
	})
}

func (s *Service) GetConfiguration(ctx context.Context, id manta.ID) (*manta.Configuration, error) {
	var (
		cf  *manta.Configuration
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		cf, err = getOrgIndexed[manta.Configuration](tx, id, ConfigurationBucket)
		return err
	})

	return cf, err
}

func (s *Service) FindConfigurations(ctx context.Context, filter manta.ConfigurationFilter) ([]*manta.Configuration, error) {
	var (
		cs  []*manta.Configuration
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		cs, err = findOrgIndexed[manta.Configuration](ctx, tx, filter.OrgID, ConfigurationBucket, ConfigurationOrgIndexBucket)
		return err
	})

	return cs, err
}

func (s *Service) UpdateConfiguration(ctx context.Context, id manta.ID, upd manta.ConfigurationUpdate) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		cf, err := getOrgIndexed[manta.Configuration](tx, id, ConfigurationBucket)
		if err != nil {
			return err
		}

		upd.Apply(cf)
		cf.Updated = time.Now()

		return putOrgIndexed(tx, cf, ConfigurationBucket, ConfigurationOrgIndexBucket)
	})
}

func (s *Service) DeleteConfiguration(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return deleteOrgIndexed[manta.Configuration](tx, id, ConfigurationBucket, ConfigurationOrgIndexBucket)
	})
}
