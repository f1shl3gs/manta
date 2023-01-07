package authorizer

import (
	"context"
	"github.com/f1shl3gs/manta/errors"

	"github.com/f1shl3gs/manta"
)

type ScrapeTargetService struct {
	scrapeTargetService manta.ScrapeTargetService
}

// FindScrapeTargetByID returns a single ScraperTarget by ID
func (s *ScrapeTargetService) FindScrapeTargetByID(ctx context.Context, id manta.ID) (*manta.ScrapeTarget, error) {
	st, err := s.scrapeTargetService.FindScrapeTargetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeRead(ctx, manta.ScrapesResourceType, id, st.OrgID); err != nil {
		return nil, err
	}

	return st, nil
}

// FindScrapeTargets returns a list of ScraperTargets that match the filter
func (s *ScrapeTargetService) FindScrapeTargets(ctx context.Context, filter manta.ScrapeTargetFilter) ([]*manta.ScrapeTarget, error) {
	ss, err := s.scrapeTargetService.FindScrapeTargets(ctx, filter)
	if err != nil {
		return nil, err
	}

	filtered := ss[:0]
	for _, st := range ss {
		_, _, err := authorizeRead(ctx, manta.ScrapesResourceType, st.ID, st.OrgID)
		if err != nil && errors.ErrorCode(err) != errors.EUnauthorized {
			return nil, err
		}

		if errors.ErrorCode(err) == errors.EUnauthorized {
			continue
		}

		filtered = append(filtered, st)
	}

	return filtered, nil
}

// CreateScrapeTarget create a ScraperTarget
func (s *ScrapeTargetService) CreateScrapeTarget(ctx context.Context, target *manta.ScrapeTarget) error {
	if _, _, err := authorizeCreate(ctx, manta.ScrapesResourceType, target.OrgID); err != nil {
		return err
	}

	return s.scrapeTargetService.CreateScrapeTarget(ctx, target)
}

// UpdateScrapeTarget update a single ScraperTarget with chageset
// returns the new ScraperTarget after update
func (s *ScrapeTargetService) UpdateScrapeTarget(ctx context.Context, id manta.ID, upd manta.ScrapeTargetUpdate) (*manta.ScrapeTarget, error) {
	st, err := s.scrapeTargetService.FindScrapeTargetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeWrite(ctx, manta.ScrapesResourceType, id, st.OrgID); err != nil {
		return nil, err
	}

	return s.scrapeTargetService.UpdateScrapeTarget(ctx, id, upd)
}

// DeleteScrapeTarget delete a single ScraperTarget by ID
func (s *ScrapeTargetService) DeleteScrapeTarget(ctx context.Context, id manta.ID) error {
	st, err := s.scrapeTargetService.FindScrapeTargetByID(ctx, id)
	if err != nil {
		return err
	}

	if _, _, err := authorizeWrite(ctx, manta.ScrapesResourceType, st.ID, st.OrgID); err != nil {
		return err
	}

	return s.scrapeTargetService.DeleteScrapeTarget(ctx, id)
}
