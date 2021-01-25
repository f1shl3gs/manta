package launcher

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/bolt"
	"github.com/f1shl3gs/manta/checks"
	"github.com/f1shl3gs/manta/control"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/log"
	"github.com/f1shl3gs/manta/pkg/signals"
	"github.com/f1shl3gs/manta/task/backend"
	"github.com/f1shl3gs/manta/task/backend/coordinator"
	"github.com/f1shl3gs/manta/task/backend/executor"
	"github.com/f1shl3gs/manta/task/backend/middleware"
	"github.com/f1shl3gs/manta/task/backend/scheduler"
	"github.com/f1shl3gs/manta/task/mock"
	"github.com/f1shl3gs/manta/web"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerzap "github.com/uber/jaeger-client-go/log/zap"
	jaegerprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type Launcher struct {
	// misc
	LogLevel string
	LogFile  string

	// service
	GrpcAddress string

	// http service
	HTTPAddress string
	AccessLog   bool

	// bolt store
	BoltPath string

	// opentracing
	Opentracing bool

	// event retention
	EventRetention time.Duration

	// scheduler
	noopSchedule bool
	WorkerLimit  int
}

func (l *Launcher) Options() []Option {
	return []Option{
		{
			DestP:   &l.GrpcAddress,
			Flag:    "grpc.address",
			Default: ":8081",
		},
		{
			DestP:   &l.HTTPAddress,
			Flag:    "http.address",
			Default: ":8080",
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
			DestP:   &l.EventRetention,
			Flag:    "event.retention",
			Default: 14 * 24 * time.Hour,
		},
		{
			DestP:   &l.LogLevel,
			Flag:    "file.level",
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

	defer logger.Sync()

	if l.Opentracing {
		logger.Info("opentracing is enabled")
		closer, otErr := setupOpentracing(logger)
		if otErr != nil {
			return otErr
		}

		defer closer.Close()
	}

	// starting services
	ctx := signals.WithStandardSignals(context.Background())
	group, ctx := errgroup.WithContext(ctx)

	store := bolt.NewKVStore(logger, l.BoltPath)
	if err = store.Open(ctx); err != nil {
		return err
	}
	defer store.Close()

	if err = kv.Initial(ctx, store); err != nil {
		return err
	}

	service := kv.NewService(logger, store)

	prometheus.MustRegister(store)

	// checks
	checker := checks.NewChecker(
		logger.With(zap.String("service", "checker")),
		service, service, service)

	// scheduler
	ex := executor.NewExecutor(logger, service, mock.NewTaskControlService(), checker.Process)
	var sch scheduler.Scheduler = &scheduler.NoopScheduler{}
	if !l.noopSchedule {
		tsch, sm, err := scheduler.NewScheduler(ex, backend.NewSchedulableTaskService(service),
			scheduler.WithMaxConcurrentWorkers(l.WorkerLimit),
			scheduler.WithOnErrorFn(func(ctx context.Context, taskID scheduler.ID, scheduledFor time.Time, err error) {
				logger.Warn("schedule task failed",
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

	{
		gl := logger.With(zap.String("service", "grpc"))
		// grpc service
		// Shared options for the logger, with a custom gRPC code to file level function.
		opts := []grpc_recovery.Option{
			grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
				gl.Error("grpc handler panicked",
					zap.Any("error", p))
				fmt.Println(string(debug.Stack()))
				return status.Errorf(codes.Unknown, "panic triggered: %v", p)
			}),
		}

		svr := grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				grpc_recovery.UnaryServerInterceptor(opts...),
				grpc_prometheus.UnaryServerInterceptor,
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				grpc_validator.UnaryServerInterceptor(),
			),
			grpc_middleware.WithStreamServerChain(
				grpc_recovery.StreamServerInterceptor(opts...),
				grpc_prometheus.StreamServerInterceptor,
				otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer()),
				grpc_validator.StreamServerInterceptor(),
			),
		)

		grpc_prometheus.Register(svr)

		ctl := control.New(logger)
		ctl.NodeService = service
		ctl.OrganizationService = service
		ctl.CheckService = checkService
		manta.RegisterControlServer(svr, ctl)

		// todo: use our own healthz implement
		hc := health.NewServer()
		hc.SetServingStatus("manta", healthpb.HealthCheckResponse_SERVING)
		healthpb.RegisterHealthServer(svr, hc)

		listener, err := net.Listen("tcp", l.GrpcAddress)
		if err != nil {
			return err
		}

		group.Go(func() error {
			errCh := make(chan error)
			go func() {
				logger.Info("start grpc service",
					zap.String("listen", l.GrpcAddress))
				errCh <- svr.Serve(listener)
			}()

			select {
			case <-ctx.Done():
				svr.GracefulStop()
				logger.Info("shutdown grpc server success")
				return nil
			case err := <-errCh:
				return errors.Wrap(err, "grpc server exit")
			}
		})
	}

	{
		// http service
		hl := logger.With(zap.String("service", "http"))
		handler := web.New(hl, &web.Backend{
			BackupService:        store,
			NodeService:          service,
			OrganizationService:  service,
			CheckService:         checkService,
			TaskService:          service,
			DatasourceService:    service,
			TemplateService:      service,
			UserService:          service,
			PasswordService:      service,
			AuthorizationService: service,
			OtclService:          service,
			DashboardService:     service,
		})

		group.Go(func() error {
			server := &http.Server{
				Addr:    l.HTTPAddress,
				Handler: handler,
			}

			errCh := make(chan error)
			go func() {
				logger.Info("start http service",
					zap.String("listen", l.HTTPAddress))
				errCh <- server.ListenAndServe()
			}()

			select {
			case <-ctx.Done():
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err = server.Shutdown(ctx)
				if err != nil {
					logger.Error("shutdown http server failed",
						zap.Error(err))
				} else {
					logger.Info("shutdown http server success")
				}

				return nil
			case err := <-errCh:
				logger.Error("http service exit on error",
					zap.Error(err))
				return err
			}
		})
	}

	return group.Wait()
}

func setupOpentracing(logger *zap.Logger) (io.Closer, error) {
	cf, err := jaegercfg.FromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "create jaeger config failed")
	}

	jmf := jaegerprom.New(jaegerprom.WithRegisterer(prometheus.DefaultRegisterer))
	jaegerZapLogger := jaegerzap.NewLogger(logger.With(zap.String("service", "opentracing")))

	tracer, closer, err := cf.NewTracer(
		jaegercfg.Logger(jaegerZapLogger),
		jaegercfg.Metrics(jmf),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create jaeger tracer failed")
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
