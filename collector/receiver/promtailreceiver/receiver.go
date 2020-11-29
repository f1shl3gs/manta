package promtailreceiver

import (
	"context"
	"errors"
	"time"

	kitlog "github.com/go-kit/kit/log"
	ptclient "github.com/grafana/loki/pkg/promtail/client"
	promtailCfg "github.com/grafana/loki/pkg/promtail/config"
	"github.com/grafana/loki/pkg/promtail/targets"
	"github.com/grafana/loki/pkg/util/flagext"
	"github.com/prometheus/common/model"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/f1shl3gs/manta/collector/receiver/promtailreceiver/internal"
)

type pReceiver struct {
	logger        *zap.Logger
	next          consumer.LogsConsumer
	targetManager *targets.TargetManagers
	ctx           context.Context
}

// implement EntryHandler
func (p *pReceiver) Handle(labels model.LabelSet, time time.Time, entry string) error {
	ctx := p.ctx
	logs := convert(labels, time, entry)

	return p.next.ConsumeLogs(ctx, logs)
}

type shutdownable struct {
	*pReceiver
}

func (s *shutdownable) Shutdown() {
	panic("implement me")
}

func (p *pReceiver) Start(ctx context.Context, host component.Host) error {

	return nil
}

func (p *pReceiver) Shutdown(ctx context.Context) error {
	return nil
}

func newPromtailReceiver(logger *zap.Logger, text string, next consumer.LogsConsumer) (*pReceiver, error) {
	var cfg promtailCfg.Config

	err := yaml.UnmarshalStrict([]byte(text), &cfg)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(cfg.ScrapeConfig); i++ {
		if cfg.ScrapeConfig[i].PushConfig != nil {
			cfg.ScrapeConfig[i].PushConfig.Server.LogLevel = cfg.ServerConfig.LogLevel
			cfg.ScrapeConfig[i].PushConfig.Server.LogFormat = cfg.ServerConfig.LogFormat
		}
	}

	receiver := &pReceiver{
		logger: logger,
		next:   next,
	}
	kl := internal.NewZapToGokitLogAdapter(logger)
	cli, err := newMulti(kl, cfg.ClientConfig.ExternalLabels, cfg.ClientConfigs...)
	if err != nil {
		return nil, err
	}

	tms, err := targets.NewTargetManagers(
		&shutdownable{pReceiver: receiver},
		kl,
		cfg.PositionsConfig,
		cli,
		cfg.ScrapeConfig,
		&cfg.TargetConfig)

	if err != nil {
		return nil, err
	}

	receiver.targetManager = tms

	return receiver, nil
}

func newMulti(logger kitlog.Logger, externalLabels flagext.LabelSet, cfgs ...ptclient.Config) (ptclient.Client, error) {
	if len(cfgs) == 0 {
		return nil, errors.New("at least one client config should be provided")
	}

	var clients ptclient.MultiClient
	for _, cfg := range cfgs {
		cfg.ExternalLabels = flagext.LabelSet{
			LabelSet: externalLabels.Merge(cfg.ExternalLabels.LabelSet),
		}

		cli, err := ptclient.New(cfg, logger)
		if err != nil {
			return nil, err
		}

		clients = append(clients, cli)
	}

	return clients, nil
}
