package lokiexporter

import (
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configmodels"
)

type Config struct {
	configmodels.ExporterSettings `mapstructure:",squash"`

	confighttp.HTTPClientSettings `mapstructure:",squash"`
}
