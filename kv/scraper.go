package kv

import (
	"context"
	"errors"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	scraperTargetBucket         = []byte("scrapertargets")
	scraperTargetOrgIndexBucket = []byte("scrapertargetorgindex")

	ErrInvalidFilter = errors.New("invalid filter")
)

func (s *Service) FindScraperTargetByID(ctx context.Context, id manta.ID) (*manta.ScrapeTarget, error) {
	var (
		target *manta.ScrapeTarget
		err    error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		target, err = s.findScraperTargetByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return target, nil
}

func (s *Service) findScraperTargetByID(ctx context.Context, tx Tx, id manta.ID) (*manta.ScrapeTarget, error) {
	key, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(scraperTargetBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	target := &manta.ScrapeTarget{}
	if err = target.Unmarshal(data); err != nil {
		return nil, err
	}

	return target, nil
}

func (s *Service) FindScraperTargets(ctx context.Context, filter manta.ScraperTargetFilter) ([]*manta.ScrapeTarget, error) {
	var (
		targets []*manta.ScrapeTarget
		err     error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		targets, err = s.findScraperTargets(ctx, tx, filter)
		return err
	})

	if err != nil {
		return nil, err
	}

	return targets, nil
}

func (s *Service) findScraperTargets(ctx context.Context, tx Tx, filter manta.ScraperTargetFilter) ([]*manta.ScrapeTarget, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var targets []*manta.ScrapeTarget

	if filter.OrgID != nil {
		prefix, err := filter.OrgID.Encode()
		if err != nil {
			return nil, err
		}

		b, err := tx.Bucket(scraperTargetOrgIndexBucket)
		if err != nil {
			return nil, err
		}

		cur, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix))
		if err != nil {
			return nil, err
		}

		ids := make([][]byte, 0, 16)
		err = WalkCursor(ctx, cur, func(k, v []byte) error {
			ids = append(ids, v)
			return nil
		})

		if err != nil {
			return nil, err
		}

		b, err = tx.Bucket(scraperTargetBucket)
		if err != nil {
			return nil, err
		}

		values, err := b.GetBatch(ids...)
		if err != nil {
			return nil, err
		}

		targets = make([]*manta.ScrapeTarget, 0, len(values))
		for i := 0; i < len(values); i++ {
			target := &manta.ScrapeTarget{}
			if err = target.Unmarshal(values[i]); err != nil {
				return nil, err
			}

			targets = append(targets, target)
		}

		return targets, nil
	}

	return nil, ErrInvalidFilter
}

func (s *Service) CreateScraperTarget(ctx context.Context, target *manta.ScrapeTarget) error {
	target.ID = s.idGen.ID()

	return s.kv.Update(ctx, func(tx Tx) error {
		return s.putScraperTarget(ctx, tx, target)
	})
}

func (s *Service) putScraperTarget(ctx context.Context, tx Tx, target *manta.ScrapeTarget) error {
	pk, err := target.ID.Encode()
	if err != nil {
		return err
	}

	fk, err := target.OrgID.Encode()
	if err != nil {
		return err
	}

	idxKey := IndexKey(fk, pk)

	b, err := tx.Bucket(scraperTargetOrgIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Put(idxKey, pk); err != nil {
		return err
	}

	b, err = tx.Bucket(scraperTargetBucket)
	if err != nil {
		return err
	}

	data, err := target.Marshal()
	if err != nil {
		return err
	}

	return b.Put(pk, data)
}

func (s *Service) UpdateScraperTarget(ctx context.Context, id manta.ID, u manta.ScraperTargetUpdate) (*manta.ScrapeTarget, error) {
	panic("implement me")
}

func (s *Service) updateScraperTarget(ctx context.Context, tx Tx, id manta.ID, u manta.ScraperTargetUpdate) (*manta.ScrapeTarget, error) {
	return nil, nil
}

func (s *Service) DeleteScraperTarget(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteScraperTarget(ctx, tx, id)
	})
}

func (s *Service) deleteScraperTarget(ctx context.Context, tx Tx, id manta.ID) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.LogKV("id", id.String())

	st, err := s.findScraperTargetByID(ctx, tx, id)
	if err != nil {
		return err
	}

	pk, err := id.Encode()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(scraperTargetBucket)
	if err != nil {
		return err
	}

	err = b.Delete(pk)
	if err != nil {
		return err
	}

	// delete orgID index
	fk, _ := st.OrgID.Encode()
	b, err = tx.Bucket(scraperTargetOrgIndexBucket)
	if err != nil {
		return err
	}

	return b.Delete(fk)
}
