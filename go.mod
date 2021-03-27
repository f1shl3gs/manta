module github.com/f1shl3gs/manta

go 1.16

require (
	github.com/benbjohnson/clock v1.0.3
	github.com/cespare/xxhash v1.1.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-kit/kit v0.10.0
	github.com/gogo/protobuf v1.3.2
	github.com/google/btree v1.0.0
	github.com/google/go-cmp v0.5.4
	github.com/influxdata/cron v0.0.0-20200427044617-e2c242fdf59f
	github.com/json-iterator/go v1.1.10
	github.com/jsternberg/zap-logfmt v1.2.0
	github.com/julienschmidt/httprouter v1.3.1-0.20200114094804-8c9f31f047a3
	github.com/mattn/go-isatty v0.0.12
	github.com/mileusna/useragent v1.0.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.15.0
	github.com/prometheus/prometheus v1.8.2-0.20201119181812-c8f810083d3f
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/thanos-io/thanos v0.18.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible
	go.etcd.io/bbolt v1.3.5
	go.uber.org/zap v1.16.0
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/tools v0.0.0-20210106214847-113979e3529a
	google.golang.org/grpc v1.33.1
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2
