module github.com/f1shl3gs/manta/collector/receiver/kmsgreceiver

go 1.15

require go.opentelemetry.io/collector v0.15.0

// Same as Cortex, we can't upgrade to grpc 1.30.0 until go.etcd.io/etcd will support it.
replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

// Keeping this same as Cortex to avoid dependency issues.
replace k8s.io/client-go => k8s.io/client-go v0.19.2

replace github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v36.2.0+incompatible

replace github.com/prometheus/prometheus v2.5.0+incompatible => github.com/prometheus/prometheus v1.8.2-0.20201014093524-73e2ce1bd643

// github.com/f1shl3gs/manta/collector/receiver/kmsgreceiver
// github.com/f1shl3gs/manta/collector/receiver/kmsgreceiver
