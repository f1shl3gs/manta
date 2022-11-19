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

func (s *Service) FindScraperTargetByID(ctx context.Context, id manta.ID) (*manta.ScrapeTarget, error) {
	var (
		st  *manta.ScrapeTarget
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		st, err = getOrgIndexed[manta.ScrapeTarget](tx, id, ScraperBucket)
		return err
	})

	return st, err
}

func (s *Service) FindScraperTargets(ctx context.Context, filter manta.ScraperTargetFilter) ([]*manta.ScrapeTarget, error) {
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

func (s *Service) CreateScraperTarget(ctx context.Context, target *manta.ScrapeTarget) error {
	now := time.Now()
	target.ID = s.idGen.ID()
	target.Created = now
	target.Updated = now

	return s.kv.Update(ctx, func(tx Tx) error {
		return putOrgIndexed(tx, target, ScraperBucket, ScrapeOrgIndexBucket)
	})
}

func (s *Service) UpdateScraperTarget(ctx context.Context, id manta.ID, u manta.ScraperTargetUpdate) (*manta.ScrapeTarget, error) {
	var (
		st  *manta.ScrapeTarget
		err error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		st, err = getOrgIndexed[manta.ScrapeTarget](tx, id, ScraperBucket)

		u.Apply(st)
		st.Updated = time.Now()

		return err
	})

	return st, err
}

func (s *Service) DeleteScraperTarget(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return deleteOrgIndexed[manta.ScrapeTarget](tx, id, ScraperBucket, ScrapeOrgIndexBucket)
	})
}
