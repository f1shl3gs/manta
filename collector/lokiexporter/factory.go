package lokiexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	typeStr = "loki"
)

func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithLogs(createLogsExporter),
	)
}

func createDefaultConfig() configmodels.Exporter {
	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
	}
}

func createLogsExporter(_ context.Context, params component.ExporterCreateParams, config configmodels.Exporter) (component.LogsExporter, error) {
	cfg := config.(*Config)
	cli, err := cfg.ToClient()
	if err != nil {
		return nil, err
	}

	return newLogsExporter(cfg, params.Logger, cli)
}
