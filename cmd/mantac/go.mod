module github.com/f1shl3gs/manta/cmd/mantac

go 1.15

require (
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticexporter v0.19.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.19.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/hostobserver v0.19.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/k8sobserver v0.19.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.19.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor v0.19.0
	go.opentelemetry.io/collector v0.19.0
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.0.0-00010101000000-000000000000 => github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.19.0

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.0.0-00010101000000-000000000000 => github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.19.0
