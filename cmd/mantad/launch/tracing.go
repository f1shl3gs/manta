package launch

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerzap "github.com/uber/jaeger-client-go/log/zap"
	jaegerprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
)

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
