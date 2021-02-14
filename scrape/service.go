package scrape

import (
	"context"
	"github.com/f1shl3gs/manta"
)

type Service struct {
}

func (s *Service) FindScraperTargetByID(ctx context.Context, id manta.ID) (*manta.ScrapeTarget, error) {
	panic("implement me")
}

func (s *Service) FindScraperTargets(ctx context.Context, filter manta.ScraperTargetFilter) ([]*manta.ScrapeTarget, error) {
	panic("implement me")
}

func (s *Service) CreateScraperTarget(ctx context.Context, target *manta.ScrapeTarget) error {
	panic("implement me")
}

func (s *Service) UpdateScraperTarget(ctx context.Context, id manta.ID, u manta.ScraperTargetUpdate) (*manta.ScrapeTarget, error) {
	panic("implement me")
}

func (s *Service) DeleteScraperTarget(ctx context.Context, id manta.ID) error {
	panic("implement me")
}
