package launcher

import (
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/f1shl3gs/manta/pkg/cgroups"
	"github.com/f1shl3gs/manta/pkg/profiler"
	tsdb3 "github.com/f1shl3gs/manta/pkg/tsdb"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/scrape"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/bolt"
	"github.com/f1shl3gs/manta/checks"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/kv/migration"
	"github.com/f1shl3gs/manta/log"
	"github.com/f1shl3gs/manta/pkg/signals"
	"github.com/f1shl3gs/manta/task/backend"
	"github.com/f1shl3gs/manta/task/backend/coordinator"
	"github.com/f1shl3gs/manta/task/backend/executor"
	"github.com/f1shl3gs/manta/task/backend/middleware"
	"github.com/f1shl3gs/manta/task/backend/scheduler"
	"github.com/f1shl3gs/manta/web"
)

type Launcher struct {
	//
	Listen string

	// log
	LogLevel string
	LogFile  string

	// http service
	AccessLog bool

	// bolt store
	BoltPath string

	// opentracing
	Opentracing bool

	// scheduler
	noopSchedule bool
	WorkerLimit  int

	// storage
	StorageDir string

	// pprof
	ProfileDir       string
	ProfileInterval  string
	ProfileRetention string

	// grpc
	grpcServer *grpc.Server

	// cluster

}

func (l *Launcher) Options() []Option {
	return []Option{
		{
			DestP:   &l.Listen,
			Flag:    "listen",
			Default: ":8088",
		},
		{
			DestP:   &l.AccessLog,
			Flag:    "http.access_log",
			Default: false,
		},
		{
			DestP:   &l.BoltPath,
			Flag:    "bolt.path",
			Default: "manta.bolt",
		},
		{
			DestP:   &l.Opentracing,
			Flag:    "opentracing",
			Default: true,
		},
		{
			DestP:   &l.LogLevel,
			Flag:    "log.level",
			Default: "info",
		},
		{
			DestP:   &l.WorkerLimit,
			Flag:    "scheduler.worker",
			Default: 128,
		},
		{
			DestP:   &l.noopSchedule,
			Flag:    "scheduler.noop",
			Default: false,
		},
		{
			DestP:   &l.StorageDir,
			Flag:    "storage.dir",
			Default: "data",
			Desc:    "storage is disabled by default",
		},
		{
			DestP:   &l.ProfileDir,
			Flag:    "profile.dir",
			Default: "",
			EnvVar:  "MANTA_PROFILE_DIR",
		},
	}
}

func (l *Launcher) Logger() (*zap.Logger, error) {
	var lvl zapcore.Level
	if err := lvl.Set(l.LogLevel); err != nil {
		return nil, err
	}

	lcf := log.Config{Level: lvl}

	var w io.Writer = os.Stdout
	if l.LogFile != "" {
		f, err := os.OpenFile(l.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}

		w = f
	}

	return lcf.New(w)
}

func (l *Launcher) Run() error {
	// setup logger
	var lvl zapcore.Level
	if err := lvl.Set(l.LogLevel); err != nil {
		return err
	}

	lcf := log.Config{Level: lvl}
	logger, err := lcf.New(os.Stdout)
	if err != nil {
		return err
	}

	defer func() {
		err = logger.Sync()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "flush logs failed\n")
		}
	}()

	CPUToUse := adjustMaxProcs()
	if CPUToUse != 0 {
		logger.Info("starting mantad",
			zap.Int("GOMAXPROCS", CPUToUse))
	}

	if l.Opentracing {
		logger.Info("Opentracing is enabled")
		closer, otErr := setupOpentracing(logger)
		if otErr != nil {
			return otErr
		}

		defer closer.Close()
	}

	// init tsdb storage
	// for now only local TenantStorage is available, aka MultiTSDB
	var tenantStorage tsdb3.TenantStorage
	{
		tsdbOpts := &tsdb.Options{
			MinBlockDuration:  int64(2 * time.Hour / time.Millisecond),
			MaxBlockDuration:  int64(2 * time.Hour / time.Millisecond),
			RetentionDuration: int64(15 * 24 * time.Hour / time.Millisecond),
			NoLockfile:        false,
			WALCompression:    true,
		}

		multitsdb := tsdb3.NewMultiTSDB(l.StorageDir, logger, prometheus.DefaultRegisterer, tsdbOpts, nil, false)
		err = multitsdb.Open()
		if err != nil {
			return err
		}

		defer func() {
			logger.Info("Starting flush storage")
			start := time.Now()
			err = multitsdb.Flush()
			if err != nil {
				logger.Error("Flush storage failed",
					zap.Error(err))
			} else {
				logger.Error("Flush storage success",
					zap.Duration("time", time.Since(start)))
			}

			err = multitsdb.Close()
			if err != nil {
				logger.Error("Close storage failed",
					zap.Error(err))
			}
		}()

		tenantStorage = multitsdb
	}

	// starting services
	ctx := signals.WithStandardSignals(context.Background())
	group, ctx := errgroup.WithContext(ctx)

	kvStore := bolt.NewKVStore(logger, l.BoltPath)
	if err = kvStore.Open(ctx); err != nil {
		return err
	}
	defer kvStore.Close()

	migrator := migration.New(logger, kvStore, migration.All...)
	err = migrator.Up(ctx)
	if err != nil {
		return err
	}

	service := kv.NewService(logger, kvStore)

	prometheus.MustRegister(kvStore)

	var scrapeManager *scrape.Manager
	// scrape service
	{
		var scrapeTargetService manta.ScraperTargetService = service
		syncTargetsCh := func(orgID manta.ID, mgr *scrape.Manager) chan map[string][]*targetgroup.Group {
			ch := make(chan map[string][]*targetgroup.Group)

			go func() {
				ticker := time.NewTicker(time.Minute)
				defer ticker.Stop()

				for {
					// sync tset as soon as possible
					targets, err := scrapeTargetService.FindScraperTargets(ctx, manta.ScraperTargetFilter{OrgID: &orgID})
					if err != nil {
						logger.Warn("Find scrape targets failed",
							zap.Error(err))
					} else {
						scf := &config.Config{
							GlobalConfig: config.GlobalConfig{
								ScrapeInterval: model.Duration(15 * time.Second),
								ScrapeTimeout:  model.Duration(10 * time.Second),
							},
						}

						tset := make(map[string][]*targetgroup.Group)
						for _, tg := range targets {
							scf.ScrapeConfigs = append(scf.ScrapeConfigs, &config.ScrapeConfig{
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
							logger.Warn("Apply scrape config failed",
								zap.Error(err))
						} else {
							logger.Debug("Apply scrape config success")
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

		orgs, _, err := service.FindOrganizations(ctx, manta.OrganizationFilter{})
		if err != nil {
			return err
		}

		for _, org := range orgs {
			app, err := tenantStorage.Appendable(ctx, org.ID)
			if err != nil {
				return err
			}

			kl := log.NewZapToGokitLogAdapter(logger.With(
				zap.String("service", "scrape"),
				zap.String("tenant", org.ID.String())))

			orgID := org.ID
			group.Go(func() error {
				mgr := scrape.NewManager(kl, app)
				errCh := make(chan error)
				go func() {
					scrapeManager = mgr
					errCh <- mgr.Run(syncTargetsCh(orgID, mgr))
				}()

				select {
				case <-ctx.Done():
					mgr.Stop()
					return nil
				case err := <-errCh:
					return err
				}
			})
		}
	}

	// checks
	checker := checks.NewChecker(
		logger.With(zap.String("service", "checker")),
		service, service, tenantStorage)

	var (
		taskService        manta.TaskService          = service
		taskControlService backend.TaskControlService = service
	)

	// scheduler
	ex := executor.NewExecutor(logger, taskService, taskControlService, checker.Process)
	var sch scheduler.Scheduler = &scheduler.NoopScheduler{}
	if !l.noopSchedule {
		tsch, sm, err := scheduler.NewScheduler(ex, backend.NewSchedulableTaskService(service),
			scheduler.WithMaxConcurrentWorkers(l.WorkerLimit),
			scheduler.WithOnErrorFn(func(ctx context.Context, taskID scheduler.ID, scheduledFor time.Time, err error) {
				logger.Warn("Schedule task failed",
					zap.String("task", manta.ID(taskID).String()),
					zap.Time("scheduledFor", scheduledFor),
					zap.Error(err))
			}))

		if err != nil {
			return errors.Wrap(err, "init scheduler failed")
		}

		defer tsch.Stop()

		sch = tsch
		prometheus.MustRegister(sm.PrometheusCollectors()...)
	}

	coord := coordinator.NewCoordinator(logger, sch, nil)
	checkService := middleware.NewCheckService(service, service, coord)

	if err = backend.NotifyCoordinatorOfExisting(ctx, logger, service, coord); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", l.Listen)
	if err != nil {
		return err
	}

	m := cmux.New(listener)

	group.Go(m.Serve)

	group.Go(func() error {
		// todo: make profile path configurable
		path := "profile"

		err = os.MkdirAll(path, 0700)
		if err != nil {
			return err
		}

		hp := profiler.NewHeapProfiler(path, 0)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
			}

			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)

			err = hp.MaybeTakeProfile(ctx, ms.HeapAlloc)
			if err != nil {
				logger.Warn("take heap profile failed",
					zap.Error(err))
			}
		}
	})

	{
		// http service
		hl := logger.With(zap.String("service", "http"))
		handler := web.New(hl, &web.Backend{
			BackupService:               kvStore,
			OrganizationService:         service,
			CheckService:                checkService,
			TaskService:                 service,
			TemplateService:             service,
			UserService:                 service,
			PasswordService:             service,
			AuthorizationService:        service,
			OtclService:                 service,
			DashboardService:            service,
			Keyring:                     service,
			SessionService:              service,
			ScrapeService:               service,
			TenantStorage:               tenantStorage,
			NotificationEndpointService: service,
			VariableService:             service,
			Flusher:                     kvStore,
		}, l.AccessLog, scrapeManager)

		_ = m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
		httpListener := m.Match(cmux.HTTP1Fast(http.MethodPatch))

		group.Go(func() error {
			server := &http.Server{
				Handler: handler,
			}

			errCh := make(chan error)
			go func() {
				logger.Info("Start HTTP service",
					zap.String("listen", l.Listen))
				errCh <- server.Serve(httpListener)
			}()

			select {
			case <-ctx.Done():
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err = server.Shutdown(ctx)
				if err != nil {
					logger.Error("Shutdown http server failed",
						zap.Error(err))
				} else {
					logger.Info("Shutdown http server success")
				}

				return nil
			case err := <-errCh:
				logger.Error("HTTP service exit on error",
					zap.Error(err))
				return err
			}
		})
	}

	return group.Wait()
}

// adjustMaxProcs set GOMAXPROCS (if not overridden by env variables)
// to be the CPU limit of the current cgroup, if running inside a cgroup
// with a cpu limit lower than runtime.NumCPU(). This is preferable to
// letting it fall back to Go's default, which is runtime.NumCPU(), as
// the Go scheduler would be running more OS-level threads than can ever
// be concurrently scheduled.
func adjustMaxProcs() int {
	var defaultValue = runtime.NumCPU()

	if _, set := os.LookupEnv("GOMAXPROCS"); set {
		return defaultValue
	}

	cpuInfo, err := cgroups.GetCgroupCPU()
	if err != nil {
		return defaultValue
	}

	numCPUToUse := int(math.Ceil(cpuInfo.CPUShares()))
	if numCPUToUse < runtime.NumCPU() && numCPUToUse > 0 {
		runtime.GOMAXPROCS(numCPUToUse)
	}

	return numCPUToUse
}
