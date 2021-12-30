module github.com/f1shl3gs/manta

go 1.17

require (
	github.com/armon/go-metrics v0.3.3 // indirect
	github.com/benbjohnson/clock v1.0.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-kit/kit v0.10.0
	github.com/gogo/protobuf v1.3.2
	github.com/google/btree v1.0.1
	github.com/google/go-cmp v0.5.6
	github.com/gorilla/websocket v1.4.2
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
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.32.1
	github.com/prometheus/prometheus v1.8.2-0.20211119115433-692a54649ed7
	github.com/soheilhy/cmux v0.1.5
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible
	go.etcd.io/bbolt v1.3.6
	go.etcd.io/etcd/pkg/v3 v3.5.1
	go.etcd.io/etcd/raft/v3 v3.5.1
	go.etcd.io/etcd/server/v3 v3.5.1
	go.uber.org/zap v1.17.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/tools v0.1.7
	google.golang.org/grpc v1.40.0
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/client-go v12.0.0+incompatible // indirect
)

require (
	github.com/alecthomas/units v0.0.0-20210927113745-59d0afb8317a // indirect
	github.com/aws/aws-sdk-go v1.42.8 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1 // indirect
	golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20211021150943-2b146023228c // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	// At the time of writing (i.e. as of this version below) the `etcd` repo is in the process of properly introducing
	// modules, and as part of that uses an unsatisfiable version for this dependency (v3.0.0-00010101000000-000000000000).
	// We just force it to the same SHA as the `go.etcd.io/etcd/raft/v3` module (they live in the same VCS root).
	//
	// While this is necessary, make sure that the require block above does not diverge.
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20201109164711-01844fd28560
	k8s.io/client-go => k8s.io/client-go v0.20.0

	// From Prometheus.
	k8s.io/klog => github.com/simonpasquier/klog-gokit v0.3.0
	k8s.io/klog/v2 => github.com/simonpasquier/klog-gokit/v2 v2.0.1
)
