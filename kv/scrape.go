package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
)

var (
	ScraperBucket        = []byte("scrapes")
	ScrapeOrgIndexBucket = []byte("scrapeorgindex")
)

func (s *Service) FindScrapeTargetByID(ctx context.Context, id manta.ID) (*manta.ScrapeTarget, error) {
	var (
		st  *manta.ScrapeTarget
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		st, err = findByID[manta.ScrapeTarget](tx, id, ScraperBucket)
		return err
	})

	return st, err
}

func (s *Service) FindScrapeTargets(ctx context.Context, filter manta.ScrapeTargetFilter) ([]*manta.ScrapeTarget, error) {
	var (
		list []*manta.ScrapeTarget
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		list, err = findOrgIndexed[manta.ScrapeTarget](ctx, tx, *filter.OrgID, ScraperBucket, ScrapeOrgIndexBucket)
		return err
	})

	return list, err
}

func (s *Service) CreateScrapeTarget(ctx context.Context, target *manta.ScrapeTarget) error {
	now := time.Now()
	target.ID = s.idGen.ID()
	target.Created = now
	target.Updated = now

	return s.kv.Update(ctx, func(tx Tx) error {
		return putOrgIndexed(tx, target, ScraperBucket, ScrapeOrgIndexBucket)
	})
}

func (s *Service) UpdateScrapeTarget(ctx context.Context, id manta.ID, u manta.ScrapeTargetUpdate) (*manta.ScrapeTarget, error) {
	var (
		st  *manta.ScrapeTarget
		err error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		st, err = findByID[manta.ScrapeTarget](tx, id, ScraperBucket)

		u.Apply(st)
		st.Updated = time.Now()

		return err
	})

	return st, err
}

func (s *Service) DeleteScrapeTarget(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return deleteOrgIndexed[manta.ScrapeTarget](tx, id, ScraperBucket, ScrapeOrgIndexBucket)
	})
}
