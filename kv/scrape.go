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

func (s *Service) FindScrapeTargets(
	ctx context.Context,
	filter manta.ScrapeTargetFilter,
) ([]*manta.ScrapeTarget, error) {
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
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.createScrapeTarget(ctx, tx, target)
	})
}

func (s *Service) createScrapeTarget(ctx context.Context, tx Tx, scrape *manta.ScrapeTarget) error {
	now := time.Now()
	scrape.ID = s.idGen.ID()
	scrape.Created = now
	scrape.Updated = now

	return putOrgIndexed(tx, scrape, ScraperBucket, ScrapeOrgIndexBucket)
}

func (s *Service) UpdateScrapeTarget(
	ctx context.Context,
	id manta.ID,
	upd manta.ScrapeTargetUpdate,
) (*manta.ScrapeTarget, error) {
	var (
		st  *manta.ScrapeTarget
		err error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		st, err = findByID[manta.ScrapeTarget](tx, id, ScraperBucket)

		upd.Apply(st)
		st.Updated = time.Now()

		return err
	})

	return st, err
}

func (s *Service) DeleteScrapeTarget(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteScrapeTarget(tx, id)
	})
}

func (s *Service) deleteScrapeTarget(tx Tx, id manta.ID) error {
	return deleteOrgIndexed[manta.ScrapeTarget](tx, id, ScraperBucket, ScrapeOrgIndexBucket)
}
