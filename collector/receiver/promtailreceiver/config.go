package promtailreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

type Config struct {
	configmodels.ReceiverSettings `mapstructure:",squash"`

	Config string `mapstructure:"config"`
}

const (
	typeStr = "promtail"
)

func NewFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		receiverhelper.WithLogs(createLogsReceiver),
	)
}

func createDefaultConfig() configmodels.Receiver {
	return &Config{
		ReceiverSettings: configmodels.ReceiverSettings{
			TypeVal: "promtail",
			NameVal: "promtail",
		},
	}
}

func createLogsReceiver(
	_ context.Context,
	params component.ReceiverCreateParams,
	cfg configmodels.Receiver,
	consumer consumer.LogsConsumer,
) (component.LogsReceiver, error) {
	rcf := cfg.(*Config)

	return newPromtailReceiver(params.Logger, rcf.Config, consumer)
}
