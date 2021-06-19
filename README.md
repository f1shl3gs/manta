# Manta

### TODO

- ~~ui: implement timeMachine~~
- ui: check editor
- ui: notification endpoint page
- ui: implement 'bottomContents' for TimeMachine, so we can reuse it
- dep: replace testify with custom package in `pkg`, which should save

### Shrunk binary

MantaD is build by command `CGO_ENABLED=0 go build --ldflags "-s -w"`

#### TODO

- remove viper, we don't need it that much, and it will reduce 808KB
- specify gzip level to `-9` only shrunk 24k, which is not much, and the decompression will take more time

| Size | Delta | Changes
|---- | ---- | ----|
| 37322752 | - | - |
| 37322752 | 0 | remove direct dependency `github.com/cespare/xxhash from scheduler` |
| 37298176 | -24k | remove `go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp` |
| 37298176 | 0 | remove direct dependency `github.com/dustin/go-humanize`  |
| 37064704 | -228k | remove direct dependency `github.com/prometheus/client_golang/prometheus/promhttp`
| 27619328 | -9224 K | embed assets.tgz instead of raw files |
