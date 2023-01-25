package launch

import (
	"context"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorizer"
	"github.com/f1shl3gs/manta/bolt"
	"github.com/f1shl3gs/manta/checks"
	httpservice "github.com/f1shl3gs/manta/http"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/kv/migration"
	"github.com/f1shl3gs/manta/multitsdb"
	"github.com/f1shl3gs/manta/oplog"
	"github.com/f1shl3gs/manta/pkg/cgroups"
	"github.com/f1shl3gs/manta/pkg/log"
	"github.com/f1shl3gs/manta/pkg/signals"
	"github.com/f1shl3gs/manta/raftstore"
	"github.com/f1shl3gs/manta/raftstore/pb"
	"github.com/f1shl3gs/manta/scrape"
	"github.com/f1shl3gs/manta/task/backend"
	"github.com/f1shl3gs/manta/task/backend/coordinator"
	"github.com/f1shl3gs/manta/task/backend/executor"
	"github.com/f1shl3gs/manta/task/backend/middleware"
	"github.com/f1shl3gs/manta/task/backend/scheduler"
	"github.com/f1shl3gs/manta/telemetry/prom"

	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Launcher struct {
	Listen string

	// log
	LogLevel string

	// opentracing
	Opentracing bool

	// scheduler
	noopSchedule bool
	WorkerLimit  int

	// Store used to specify store type, the available value could be "bolt" or "raftstore"
	Store string
	// where to store boltdb file or raftstore's files, e.g. WAL, snapshot, and FSM
	StorePath string

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
			DestP:   &l.Store,
			Flag:    "store",
			Default: "bolt",
		},
		{
			DestP:   &l.StorePath,
			Flag:    "store.path",
			Default: "manta",
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

	listener, err := net.Listen("tcp", l.Listen)
	if err != nil {
		return err
	}

	muxer := cmux.New(listener)
	muxer.HandleError(func(err error) bool {
		logger.Error("cmux handle failed",
			zap.Error(err))
		return true
	})

	// starting services
	ctx := signals.WithStandardSignals(context.Background())
	group, ctx := errgroup.WithContext(ctx)

	var (
		kvStore        kv.SchemaStore
		flusher        httpservice.Flusher
		clusterService raftstore.ClusterService
		promRegistry   = prom.NewRegistry(logger)

		grpcSvr = grpc.NewServer()
	)

	switch l.Store {
	case "bolt":
		bs := bolt.NewKVStore(logger, filepath.Join(l.StorePath, "manta.bolt"))
		if err = bs.Open(ctx); err != nil {
			return err
		}

		kvStore = bs
		flusher = bs
		promRegistry.MustRegister(bs.Collectors()...)

		defer bs.Close()
	case "raftstore":
		rs, err := raftstore.New(&raftstore.Config{
			DataDir:      filepath.Join(l.StorePath, "raft"),
			Listen:       l.Listen,
			DefragOnBoot: false,
		}, logger)
		if err != nil {
			return err
		}

		kvStore = rs
		clusterService = rs
		promRegistry.MustRegister(rs.Collectors()...)
		pb.RegisterRaftServer(grpcSvr, rs)

		group.Go(func() error {
			rs.Run(ctx)
			return nil
		})
	default:
		return errors.Errorf("unknown store type %q", l.Store)
	}

	kvStore = kv.NewMetricService(kvStore, promRegistry)
	migrator := migration.New(logger, kvStore, migration.All...)
	err = migrator.Up(ctx)
	if err != nil {
		return errors.Wrap(err, "migrate failed")
	}

	service := kv.NewService(logger, kvStore)

	var (
		orgService       manta.OrganizationService = service
		oplogService     manta.OperationLogService = service
		dashboardService manta.DashboardService    = service
		taskService      manta.TaskService         = service
		configService    manta.ConfigService       = service
	)

	var tenantStorage multitsdb.TenantStorage
	{
		tsdbOpts := &tsdb.Options{
			MinBlockDuration:  int64(2 * time.Hour / time.Millisecond),
			MaxBlockDuration:  int64(2 * time.Hour / time.Millisecond),
			RetentionDuration: int64(15 * 24 * time.Hour / time.Millisecond),
			NoLockfile:        false,
			WALCompression:    true,
		}

		mtsdb := multitsdb.NewMultiTSDB(l.StorageDir, logger, promRegistry, tsdbOpts, nil, false)
		if err := mtsdb.Open(); err != nil {
			return err
		}

		defer func() {
			logger.Info("Staring flush storage")

			start := time.Now()
			if err = mtsdb.Flush(); err != nil {
				logger.Error("Flush tsdb failed", zap.Error(err))
			} else {
				logger.Error("Flush tsdb success", zap.Duration("duration", time.Since(start)))
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
			promRegistry.MustRegister(sm.Collectors()...)
		}

		coord := coordinator.NewCoordinator(logger, sch, nil)
		checkService = middleware.NewCheckService(service, taskService, coord)

		if err = backend.NotifyCoordinatorOfExisting(ctx, logger, service, coord); err != nil {
			return err
		}
	}

	{
		// grpc service
		grpcListener := muxer.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))

		group.Go(func() error {
			logger.Info("start grpc service")
			return grpcSvr.Serve(grpcListener)
		})
	}

	{
		// http service
		httpListener := muxer.Match(cmux.HTTP1Fast(http.MethodPatch))

		checkService = oplog.NewCheckService(checkService, oplogService, logger)
		configService = oplog.NewConfigService(configService, oplogService, logger)
		dashboardService = oplog.NewDashboardService(dashboardService, oplogService, logger)

		hl := logger.With(zap.String("service", "http"))
		handler := httpservice.New(hl, &httpservice.Backend{
			PromRegistry:                promRegistry,
			OnBoardingService:           service,
			BackupService:               kvStore,
			CheckService:                authorizer.NewCheckService(checkService),
			TaskService:                 taskService,
			OrganizationService:         service,
			UserService:                 service,
			PasswordService:             service,
			AuthorizationService:        service,
			DashboardService:            authorizer.NewDashboardService(dashboardService),
			SessionService:              service,
			Flusher:                     flusher,
			ConfigService:               authorizer.NewConfigService(configService),
			ScrapeTargetService:         authorizer.NewScrapeTargetService(scrapeTargetService),
			RegistryService:             service,
			SecretService:               service,
			NotificationEndpointService: service,
			OperationLogService:         oplogService,
			TenantStorage:               tenantStorage,
			TenantTargetRetriever:       targetRetrievers,
			ClusterService:              clusterService,
		})

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

	group.Go(func() error {
		errCh := make(chan error)

		go func() {
			errCh <- muxer.Serve()
		}()

		select {
		case <-ctx.Done():
			muxer.Close()
			return nil

		case err := <-errCh:
			if err == os.ErrClosed {
				return nil
			}

			if err != nil {
				logger.Error("muxer serve error",
					zap.Error(err))
			}

			return err
		}
	})

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
