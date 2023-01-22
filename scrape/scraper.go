package scrape

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/log"

	"github.com/prometheus/common/model"
	promconfig "github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/scrape"
	"github.com/prometheus/prometheus/storage"
	"go.uber.org/zap"
)

type Scraper struct {
	orgID  manta.ID
	logger *zap.Logger

	mgr                 *scrape.Manager
	syncCh              chan map[string][]*targetgroup.Group
	scrapeTargetService manta.ScrapeTargetService
}

func newScraper(
	logger *zap.Logger,
	orgID manta.ID,
	appendable storage.Appendable,
	scrapeTargetService manta.ScrapeTargetService,
) *Scraper {
	logger = logger.With(zap.String("scraper", "scrape"), zap.String("org", orgID.String()))
	kl := log.NewZapToGokitLogAdapter(logger)

	mgr := scrape.NewManager(nil, kl, appendable)

	ch := newScrapPool(context.Background(), logger, orgID, mgr, scrapeTargetService)

	scraper := &Scraper{
		orgID:  orgID,
		logger: logger,

		mgr:                 mgr,
		syncCh:              ch,
		scrapeTargetService: scrapeTargetService,
	}

	go func() {
		err := mgr.Run(ch)
		if err != nil {
			logger.Error("scrape manager run failed", zap.Error(err))
		}
	}()

	go scraper.syncTargets()

	return scraper
}

func (s *Scraper) syncTargets() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	targets, err := s.scrapeTargetService.FindScrapeTargets(ctx, manta.ScrapeTargetFilter{OrgID: &s.orgID})
	if err != nil {
		s.logger.Warn("find targets failed",
			zap.Error(err))
		return
	}

	scf := &promconfig.Config{
		GlobalConfig: promconfig.GlobalConfig{
			ScrapeInterval: model.Duration(15 * time.Second),
			ScrapeTimeout:  model.Duration(10 * time.Second),
		},
	}

	tset := make(map[string][]*targetgroup.Group)
	for _, tg := range targets {
		scf.ScrapeConfigs = append(scf.ScrapeConfigs, &promconfig.ScrapeConfig{
			JobName:        tg.Name,
			ScrapeInterval: model.Duration(15 * time.Second),
			ScrapeTimeout:  model.Duration(10 * time.Second),
			MetricsPath:    "/metrics",
			Scheme:         "http",
		})

		labelSet := model.LabelSet{}
		for k, v := range tg.Labels {
			labelSet[model.LabelName(k)] = model.LabelValue(v)
		}

		targetLs := make([]model.LabelSet, 0, len(tg.Targets))
		for _, addr := range tg.Targets {
			targetLs = append(targetLs, model.LabelSet{
				model.AddressLabel: model.LabelValue(addr),
			})
		}

		tset[tg.Name] = append(tset[tg.Name], &targetgroup.Group{
			Source:  tg.ID.String(),
			Labels:  labelSet,
			Targets: targetLs,
		})
	}

	err = s.mgr.ApplyConfig(scf)
	if err != nil {
		s.logger.Warn("Apply scrape config failed", zap.Error(err))
	} else {
		s.logger.Debug("Apply scrape config success", zap.Int("jobs", len(scf.ScrapeConfigs)))
	}
}

func newScrapPool(
	ctx context.Context,
	logger *zap.Logger,
	orgID manta.ID,
	mgr *scrape.Manager,
	scrapeTargetService manta.ScrapeTargetService,
) chan map[string][]*targetgroup.Group {
	ch := make(chan map[string][]*targetgroup.Group)

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			// sync tset as soon as possible
			targets, err := scrapeTargetService.FindScrapeTargets(ctx, manta.ScrapeTargetFilter{OrgID: &orgID})
			if err != nil {
				logger.Warn("Find scrape targets failed", zap.Error(err))
			} else {
				scf := &promconfig.Config{
					GlobalConfig: promconfig.GlobalConfig{
						ScrapeInterval: model.Duration(15 * time.Second),
						ScrapeTimeout:  model.Duration(10 * time.Second),
					},
				}

				tset := make(map[string][]*targetgroup.Group)
				for _, tg := range targets {
					scf.ScrapeConfigs = append(scf.ScrapeConfigs, &promconfig.ScrapeConfig{
						JobName:        tg.Name,
						ScrapeInterval: model.Duration(15 * time.Second),
						ScrapeTimeout:  model.Duration(10 * time.Second),
						MetricsPath:    "/metrics",
						Scheme:         "http",
					})

					labelSet := model.LabelSet{}
					for k, v := range tg.Labels {
						labelSet[model.LabelName(k)] = model.LabelValue(v)
					}

					targetLs := make([]model.LabelSet, 0, len(tg.Targets))
					for _, addr := range tg.Targets {
						targetLs = append(targetLs, model.LabelSet{
							model.AddressLabel: model.LabelValue(addr),
						})
					}

					tset[tg.Name] = append(tset[tg.Name], &targetgroup.Group{
						Source:  tg.ID.String(),
						Labels:  labelSet,
						Targets: targetLs,
					})
				}

				err := mgr.ApplyConfig(scf)
				if err != nil {
					logger.Warn("Apply scrape config failed", zap.Error(err))
				} else {
					logger.Debug("Apply scrape config success", zap.Int("jobs", len(scf.ScrapeConfigs)))
				}

				select {
				case <-ctx.Done():
					return
				case ch <- tset:
				}
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()

	return ch
}
