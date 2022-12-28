package scrape

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/prometheus/scrape"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/multitsdb"
)

type CoordinatingScrapeService struct {
	logger *zap.Logger

	// services
	scrapeTargetService manta.ScrapeTargetService
	tenantStorage       multitsdb.TenantStorage

	mtx      sync.Mutex
	scrapers map[manta.ID]*Scraper
}

func New(
	ctx context.Context,
	logger *zap.Logger,
	orgService manta.OrganizationService,
	scraperTargetService manta.ScrapeTargetService,
	tenantStorage multitsdb.TenantStorage,
) (*CoordinatingScrapeService, error) {
	orgs, _, err := orgService.FindOrganizations(ctx, manta.OrganizationFilter{})
	if err != nil {
		return nil, nil
	}

	scrapers := make(map[manta.ID]*Scraper)
	for _, org := range orgs {
		app, err := tenantStorage.Appendable(ctx, org.ID)
		if err != nil {
			return nil, err
		}

		scrapers[org.ID] = newScraper(logger, org.ID, app, scraperTargetService)
	}

	return &CoordinatingScrapeService{
		logger:              logger,
		scrapeTargetService: scraperTargetService,
		tenantStorage:       tenantStorage,
		scrapers:            scrapers,
	}, nil
}

func (s *CoordinatingScrapeService) syncScraper(orgID manta.ID) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.mtx.Lock()
	defer s.mtx.Unlock()

	scraper, exist := s.scrapers[orgID]
	if !exist {
		app, err := s.tenantStorage.Appendable(ctx, orgID)
		if err != nil {
			s.logger.Warn("sync scraper failed")
			return
		}

		scraper = newScraper(s.logger, orgID, app, s.scrapeTargetService)

		s.scrapers[orgID] = scraper
	}

	go scraper.syncTargets()
}

// FindScrapeTargetByID returns a single ScraperTarget by ID
func (s *CoordinatingScrapeService) FindScrapeTargetByID(ctx context.Context, id manta.ID) (*manta.ScrapeTarget, error) {
	return s.scrapeTargetService.FindScrapeTargetByID(ctx, id)
}

// FindScrapeTargets returns a list of ScraperTargets that match the filter
func (s *CoordinatingScrapeService) FindScrapeTargets(ctx context.Context, filter manta.ScrapeTargetFilter) ([]*manta.ScrapeTarget, error) {
	return s.scrapeTargetService.FindScrapeTargets(ctx, filter)
}

func (s *CoordinatingScrapeService) CreateScrapeTarget(ctx context.Context, target *manta.ScrapeTarget) error {
	if err := s.scrapeTargetService.CreateScrapeTarget(ctx, target); err != nil {
		return err
	}

	s.syncScraper(target.OrgID)

	return nil
}

// UpdateScrapeTarget update a single ScraperTarget with chageset
// returns the new ScraperTarget after update
func (s *CoordinatingScrapeService) UpdateScrapeTarget(ctx context.Context, id manta.ID, upd manta.ScrapeTargetUpdate) (*manta.ScrapeTarget, error) {
	target, err := s.scrapeTargetService.UpdateScrapeTarget(ctx, id, upd)
	if err != nil {
		return nil, err
	}

	s.syncScraper(target.OrgID)

	return target, nil
}

// DeleteScrapeTarget delete a single ScraperTarget by ID
func (s *CoordinatingScrapeService) DeleteScrapeTarget(ctx context.Context, id manta.ID) error {
	target, err := s.scrapeTargetService.FindScrapeTargetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.scrapeTargetService.DeleteScrapeTarget(ctx, id)
	if err != nil {
		return err
	}

	s.syncScraper(target.OrgID)

	return nil
}

// TargetsActive implement TenantTargetRetriever
func (s *CoordinatingScrapeService) TargetsActive(id manta.ID) map[string][]*scrape.Target {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	scraper, exist := s.scrapers[id]
	if !exist {
		return nil
	}

	return scraper.mgr.TargetsActive()
}

// TargetsDropped implement TenantTargetRetriever
func (s *CoordinatingScrapeService) TargetsDropped(id manta.ID) map[string][]*scrape.Target {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	scraper, exist := s.scrapers[id]
	if !exist {
		return nil
	}

	return scraper.mgr.TargetsDropped()
}
