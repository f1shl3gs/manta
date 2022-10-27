package launch

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

	"github.com/f1shl3gs/manta/bolt"
	httpservice "github.com/f1shl3gs/manta/http"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/kv/migration"
	"github.com/f1shl3gs/manta/pkg/cgroups"
	"github.com/f1shl3gs/manta/pkg/log"
	"github.com/f1shl3gs/manta/pkg/signals"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerzap "github.com/uber/jaeger-client-go/log/zap"
	jaegerprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
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
		err = logger.Sync()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "flush logs failed\n")
		}
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

	{
		// http service
		listener, err := net.Listen("tcp", l.Listen)
		if err != nil {
			return err
		}

		hl := logger.With(zap.String("service", "http"))
		handler := httpservice.New(hl, &httpservice.Backend{
			OnBoardingService:    service,
			BackupService:        kvStore,
			OrganizationService:  service,
			UserService:          service,
			PasswordService:      service,
			AuthorizationService: service,
			DashboardService:     service,
			SessionService:       service,
			Flusher:              kvStore,
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
