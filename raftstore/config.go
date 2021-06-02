package raftstore

import "time"

const (
	ClusterStateFlagNew      = "new"
	ClusterStateFlagExisting = "existing"

	DefaultName                  = "default"
	DefaultMaxSnapshots          = 5
	DefaultMaxWALs               = 5
	DefaultMaxTxnOps             = uint(128)
	DefaultWarningApplyDuration  = 100 * time.Millisecond
	DefaultMaxRequestBytes       = 1.5 * 1024 * 1024
	DefaultGRPCKeepAliveMinTime  = 5 * time.Second
	DefaultGRPCKeepAliveInterval = 2 * time.Hour
	DefaultGRPCKeepAliveTimeout  = 20 * time.Second
	DefaultDowngradeCheckTime    = 5 * time.Second

	DefaultListenPeerURLs   = "http://localhost:2380"
	DefaultListenClientURLs = "http://localhost:2379"

	DefaultLogOutput = "default"
	JournalLogOutput = "systemd/journal"
	StdErrLogOutput  = "stderr"
	StdOutLogOutput  = "stdout"

	// DefaultLogRotationConfig is the default configuration used for log rotation.
	// Log rotation is disabled by default.
	// MaxSize    = 100 // MB
	// MaxAge     = 0 // days (no limit)
	// MaxBackups = 0 // no limit
	// LocalTime  = false // use computers local time, UTC by default
	// Compress   = false // compress the rotated log in gzip format
	DefaultLogRotationConfig = `{"maxsize": 100, "maxage": 0, "maxbackups": 0, "localtime": false, "compress": false}`

	// ExperimentalDistributedTracingAddress is the default collector address.
	ExperimentalDistributedTracingAddress = "localhost:4317"
	// ExperimentalDistributedTracingServiceName is the default etcd service name.
	ExperimentalDistributedTracingServiceName = "etcd"

	// DefaultStrictReconfigCheck is the default value for "--strict-reconfig-check" flag.
	// It's enabled by default.
	DefaultStrictReconfigCheck = true
	// DefaultEnableV2 is the default value for "--enable-v2" flag.
	// v2 API is disabled by default.
	DefaultEnableV2 = false

	// maxElectionMs specifies the maximum value of election timeout.
	// More details are listed in ../Documentation/tuning.md#time-parameters.
	maxElectionMs = 50000
	// backend freelist map type
	freelistArrayType = "array"
)

type Config struct {
	BindAddr  string
	DataDir   string
	WALDir    string
	SnapDir   string
	MemberDir string

	TickMs          uint
	ElectionTicks   int
	PreVote         bool
	ForceNewCluster bool
	NewCluster      bool

	MaxRequestBytes uint64
	InitialPeers    []string

	SnapshotCount uint64
	// SnapshotCatchUpEntries is the number of entries for a
	// slow follower to catch-up after compacting the raft storage
	// entries. We expect the follower has a millisecond level latency
	// with the leader. The max throughput is around 10K, Keep a 5K
	// entries is enough for helping follower to catch up.
	// WARNING: only change this for tests. Always use "DefaultSnapshotCatchUpEntries"
	SnapshotCatchUpEntries uint64

	// InitialElectionTickAdvance is true, then local member fast-forwards
	// election ticks to speed up "initial" leader election trigger. This
	// benefits the case of larger election ticks. For instance, cross
	// datacenter deployment may require longer election timeout of 10-second.
	// If true, local node does not need wait up to 10-second. Instead,
	// forwards its election ticks to 8-second, and have only 2-second left
	// before leader election.
	//
	// Major assumptions are that:
	//  - cluster has no active leader thus advancing ticks enables faster
	//    leader election, or
	//  - cluster already has an established leader, and rejoining follower
	//    is likely to receive heartbeats from the leader after tick advance
	//    and before election timeout.
	//
	// However, when network from leader to rejoining follower is congested,
	// and the follower does not receive leader heartbeat within left election
	// ticks, disruptive election has to happen thus affecting cluster
	// availabilities.
	//
	// Disabling this would slow down initial bootstrap process for cross
	// datacenter deployments. Make your own tradeoffs by configuring
	// --initial-election-tick-advance at the cost of slow initial bootstrap.
	//
	// If single-node, it advances ticks regardless.
	//
	// See https://github.com/etcd-io/etcd/issues/9333 for more detail.
	InitialElectionTickAdvance bool
	MaxTxnOps                  uint
}

// RequestTimeout returns timeout for request to finish.
func (cf *Config) RequestTimeout() time.Duration {
	// 5s for queue waiting, computation and disk IO delay
	// + 2 * election timeout for possible leader election
	return 5*time.Second + 2*time.Duration(cf.ElectionTicks*int(cf.TickMs))*time.Millisecond
}

func NewConfig() *Config {
	cfg := &Config{
		SnapshotCount:              DefaultSnapshotCount,
		TickMs:                     100,
		ElectionTicks:              1000,
		InitialElectionTickAdvance: true,
		PreVote:                    true,
		MaxRequestBytes:            DefaultMaxRequestBytes,
		MaxTxnOps:                  DefaultMaxTxnOps,

		DataDir: "data",
		WALDir:  "wal",
		SnapDir: "snapshot",
	}

	return cfg
}
