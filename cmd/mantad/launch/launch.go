package launch

import (
	"context"
	"math"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/tsdb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorizer"
	"github.com/f1shl3gs/manta/bolt"
	"github.com/f1shl3gs/manta/checks"
	httpservice "github.com/f1shl3gs/manta/http"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/kv/migration"
	"github.com/f1shl3gs/manta/multitsdb"
	"github.com/f1shl3gs/manta/pkg/cgroups"
	"github.com/f1shl3gs/manta/pkg/log"
	"github.com/f1shl3gs/manta/pkg/signals"
	"github.com/f1shl3gs/manta/scrape"
	"github.com/f1shl3gs/manta/task/backend"
	"github.com/f1shl3gs/manta/task/backend/coordinator"
	"github.com/f1shl3gs/manta/task/backend/executor"
	"github.com/f1shl3gs/manta/task/backend/middleware"
	"github.com/f1shl3gs/manta/task/backend/scheduler"
)

type Launcher struct {
	//
	Listen string

	// log
	LogLevel string

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
}

func (l *Launcher) Options() []Option {
	return []Option{
		{
			DestP:   &l.Listen,
			Flag:    "listen",
			Default: ":8088",
		},
		{
			DestP:   &l.BoltPath,
			Flag:    "bolt.path",
			Default: "manta.bolt",
		},
		{
			DestP:   &l.Opentracing,
			Flag:    "opentracing",
			Default: false,
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

func (l *Launcher) logger() (*zap.Logger, error) {
	var lvl zapcore.Level
	if err := lvl.Set(l.LogLevel); err != nil {
		return nil, err
	}

	lcf := log.Config{Level: lvl}

	return lcf.New(os.Stdout)
}

func (l *Launcher) run() error {
	// setup logger
	logger, err := l.logger()
	if err != nil {
		return err
	}

	defer func() {
		// Calling fsync on os.Stdout cause EINVAL on Linux platform, but it's fine
		// on MacOS.
		//
		// See https://github.com/uber-go/zap/issues/328#issuecomment-284337436
		_ = logger.Sync()
	}()

	CPUToUse := adjustMaxProcs()
	if CPUToUse != 0 {
		logger.Info("Starting mantad",
			zap.Int("GOMAXPROCS", CPUToUse))
	}

	if l.Opentracing {
		logger.Info("Opentracing is enabled")
		closer, otErr := setupOpentracing(logger)
		if otErr != nil {
			return otErr
		}

		defer closer.Close()
	} else {
		logger.Info("OpenTracing is disabled")
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
		return errors.Wrap(err, "migrate failed")
	}

	service := kv.NewService(logger, kvStore)

	var (
		orgService manta.OrganizationService = service
	)

	prometheus.MustRegister(kvStore)

	var tenantStorage multitsdb.TenantStorage
	{
		tsdbOpts := &tsdb.Options{
			MinBlockDuration:  int64(2 * time.Hour / time.Millisecond),
			MaxBlockDuration:  int64(2 * time.Hour / time.Millisecond),
			RetentionDuration: int64(15 * 24 * time.Hour / time.Millisecond),
			NoLockfile:        false,
			WALCompression:    true,
		}

		mtsdb := multitsdb.NewMultiTSDB(l.StorageDir, logger, prometheus.DefaultRegisterer, tsdbOpts, nil, false)
		if err := mtsdb.Open(); err != nil {
			return err
		}

		defer func() {
			logger.Info("Staring flush storage")

			start := time.Now()
			if err = mtsdb.Flush(); err != nil {
				logger.Error("Flush storage failed", zap.Error(err))
			} else {
				logger.Error("Flush storage success", zap.Duration("duration", time.Since(start)))
			}

			if err = mtsdb.Close(); err != nil {
				logger.Error("Close storage failed", zap.Error(err))
			}
		}()

		tenantStorage = mtsdb
	}

	scrapeTargetService, err := scrape.New(ctx, logger, orgService, service, tenantStorage)
	if err != nil {
		return errors.Wrap(err, "create scrape service failed")
	}

	var targetRetrievers multitsdb.TenantTargetRetriever = scrapeTargetService

	var checkService manta.CheckService
	var taskService manta.TaskService = service
	{
		// checks
		var (
			taskControlService backend.TaskControlService = service
			checker                                       = checks.NewChecker(
				logger.With(zap.String("service", "check")),
				service, tenantStorage,
			)
			sch scheduler.Scheduler = &scheduler.NoopScheduler{}
		)

		ex := executor.NewExecutor(logger, taskService, taskControlService, checker.Process)
		if !l.noopSchedule {
			tsch, sm, err := scheduler.NewScheduler(
				ex,
				backend.NewSchedulableTaskService(service),
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
		checkService = middleware.NewCheckService(service, taskService, coord)

		if err = backend.NotifyCoordinatorOfExisting(ctx, logger, service, coord); err != nil {
			return err
		}
	}

	{
		// http service
		listener, err := net.Listen("tcp", l.Listen)
		if err != nil {
			return err
		}

		hl := logger.With(zap.String("service", "http"))
		handler := httpservice.New(hl, &httpservice.Backend{
			OnBoardingService:           service,
			BackupService:               kvStore,
			CheckService:                authorizer.NewCheckService(checkService),
			TaskService:                 taskService,
			OrganizationService:         service,
			UserService:                 service,
			PasswordService:             service,
			AuthorizationService:        service,
			DashboardService:            authorizer.NewDashboardService(service),
			SessionService:              service,
			Flusher:                     kvStore,
			ConfigurationService:        service,
			ScrapeTargetService:         authorizer.NewScrapeTargetService(scrapeTargetService),
			RegistryService:             service,
			SecretService:               service,
			NotificationEndpointService: service,

			TenantStorage:         tenantStorage,
			TenantTargetRetriever: targetRetrievers,
		})

		group.Go(func() error {
			server := &http.Server{
				Handler: handler,
			}

			errCh := make(chan error)
			go func() {
				logger.Info("Start HTTP service",
					zap.String("listen", l.Listen))
				errCh <- server.Serve(listener)
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
