package scrape

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"

	"go.uber.org/zap"
)

type CoordinatingScrapeService struct {
	manta.ScraperTargetService

	logger          *zap.Logger
	registryService manta.RegistryService
}

func New(logger *zap.Logger, ss manta.ScraperTargetService, rs manta.RegistryService) *CoordinatingScrapeService {
	s := &CoordinatingScrapeService{
		ScraperTargetService: ss,
		registryService:      rs,
		logger:               logger,
	}

	go s.generate()

	return s
}

func (s *CoordinatingScrapeService) generate() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		list, err := s.registryService.Catalog(context.Background())
		if err != nil {
			s.logger.Warn("List registered instance failed", zap.Error(err))
		} else {
			now := time.Now()
			filtered := list[:0]

			for _, ins := range list {
				if ins.ExpiredAt(now) {
					continue
				}

				filtered = append(filtered, ins)
			}

		}

		<-ticker.C
	}
}
