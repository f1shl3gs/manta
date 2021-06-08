module github.com/f1shl3gs/manta

go 1.16

require (
	github.com/OneOfOne/xxhash v1.2.6 // indirect
	github.com/armon/go-metrics v0.3.3 // indirect
	github.com/benbjohnson/clock v1.0.3
	github.com/cespare/xxhash v1.1.0
	github.com/cockroachdb/pebble v0.0.0-20210526183633-dd2a545f5d75
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dustin/go-humanize v1.0.0
	github.com/go-kit/kit v0.10.0
	github.com/gogo/protobuf v1.3.2
	github.com/google/btree v1.0.0
	github.com/google/go-cmp v0.5.5
	github.com/gopherjs/gopherjs v0.0.0-20191106031601-ce3c9ade29de // indirect
	github.com/hashicorp/go-hclog v0.12.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.2.0 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/influxdata/cron v0.0.0-20200427044617-e2c242fdf59f
	github.com/jsternberg/zap-logfmt v1.2.0
	github.com/julienschmidt/httprouter v1.3.1-0.20200114094804-8c9f31f047a3
	github.com/mattn/go-isatty v0.0.12
	github.com/mileusna/useragent v1.0.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.20.0
	github.com/prometheus/prometheus v1.8.2-0.20210327162702-0f74bea24ec8
	github.com/smartystreets/assertions v1.0.1 // indirect
	github.com/soheilhy/cmux v0.1.5-0.20210205191134-5ec6847320e5
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible
	go.etcd.io/bbolt v1.3.5
	go.etcd.io/etcd/pkg/v3 v3.5.0-alpha.0
	go.etcd.io/etcd/raft/v3 v3.5.0-alpha.0
	go.etcd.io/etcd/server/v3 v3.5.0-alpha.0
	go.uber.org/zap v1.16.1-0.20210329175301-c23abee72d19
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/tools v0.1.0
	google.golang.org/grpc v1.36.1
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	k8s.io/client-go v12.0.0+incompatible // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.20.0

// At the time of writing (i.e. as of this version below) the `etcd` repo is in the process of properly introducing
// modules, and as part of that uses an unsatisfiable version for this dependency (v3.0.0-00010101000000-000000000000).
// We just force it to the same SHA as the `go.etcd.io/etcd/raft/v3` module (they live in the same VCS root).
//
// While this is necessary, make sure that the require block above does not diverge.
replace go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20201109164711-01844fd28560

// offical httprouter is not support for wildcard paths with other children
// more info: https://github.com/julienschmidt/httprouter/pull/329 and https://github.com/gin-gonic/gin/pull/2706
replace github.com/julienschmidt/httprouter v1.3.1-0.20200114094804-8c9f31f047a3 => ./third_party/httprouter
