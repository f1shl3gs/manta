module github.com/f1shl3gs/manta/collector/receiver/promtailreceiver

go 1.15

require (
	github.com/go-kit/kit v0.10.0
	github.com/grafana/loki v1.6.2-0.20201127162223-8c1fe88409fe
	github.com/pierrec/lz4 v2.5.3-0.20200429092203-e876bbd321b3+incompatible // indirect
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/collector v0.15.0
	go.uber.org/zap v1.16.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v36.2.0+incompatible

replace github.com/prometheus/prometheus v2.5.0+incompatible => github.com/prometheus/prometheus v1.8.2-0.20201014093524-73e2ce1bd643

// Keeping this same as Cortex to avoid dependency issues.
replace k8s.io/client-go => k8s.io/client-go v0.19.2

// Same as Cortex, we can't upgrade to grpc 1.30.0 until go.etcd.io/etcd will support it.
replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

// Same as Cortex
// Using a 3rd-party branch for custom dialer - see https://github.com/bradfitz/gomemcache/pull/86
replace github.com/bradfitz/gomemcache => github.com/themihai/gomemcache v0.0.0-20180902122335-24332e2d58ab
