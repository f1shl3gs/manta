package raftstore

import (
	"runtime"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/bloom"
	"github.com/cockroachdb/pebble/vfs"
	"github.com/f1shl3gs/manta/pkg/env"
)

func defaultPebbleOptions() *pebble.Options {
	// In RocksDB, the concurrency setting corresponds to both
	// flushes and compactions. In Pebble, there is always a slot
	// for a flush, and compactions are counted separately
	maxConcurrentCompactions := func() int {
		const max = 4
		if n := runtime.GOMAXPROCS(0); n <= max {
			return n
		}

		return max
	}()

	if maxConcurrentCompactions < 1 {
		maxConcurrentCompactions = 1
	}

	opts := &pebble.Options{
		// TODO: test only
		DisableWAL: false,

		Comparer:                    pebble.DefaultComparer,
		L0CompactionThreshold:       2,
		L0StopWritesThreshold:       1000,
		LBaseMaxBytes:               64 << 20, // 64 MB
		Levels:                      make([]pebble.LevelOptions, 7),
		MaxConcurrentCompactions:    maxConcurrentCompactions,
		MemTableSize:                64 << 20, // 64 MB
		MemTableStopWritesThreshold: 4,
		Merger:                      pebble.DefaultMerger,
	}

	// Automatically flush 10s after the first range tombstone is added
	// to a memtable. This ensures that we can reclaim space even when
	// there's no activity on the database generating flushes.
	opts.Experimental.DeleteRangeFlushDelay = 10 * time.Second
	// Enable deletion pacing. This helps prevent disk slowness even
	// on some SSDs, that kick off an expensive GC if a lot of files
	// are deleted at onec
	opts.Experimental.MinDeletionRate = 128 << 20 // 128 MB

	for i := 0; i < len(opts.Levels); i++ {
		l := &opts.Levels[i]
		l.BlockSize = 32 << 10
		l.IndexBlockSize = 256 << 10
		l.FilterPolicy = bloom.FilterPolicy(10)
		l.FilterType = pebble.TableFilter
		if i > 0 {
			l.TargetFileSize = opts.Levels[i-1].TargetFileSize * 2
		}
		l.EnsureDefaults()
	}

	// Do not create bloom filters for the last level (i.e. the largest
	// level which contains data in the LSM store). This configuration
	// reduces the size of the bloom filters by 10x. This is significant
	// given that bloom filters require 1.25 bytes (10 bits) per key
	// which can translate into gigabytes of memory given typical key and
	// value sizes. The downside is that bloom filters will only be usable
	// on the higher levels, but that seems acceptable. We typically see
	// read amplification of 5-6x on clusters (i.e. there are 5-6 levels
	// of sstables) which means we'll achieve 80-90% of the benefit of
	// having bloom filters on every level for only 10% of the memory cast.
	opts.Levels[6].FilterPolicy = nil

	// Set disk health check interval to min(5s, maxSyncDurationDefault).
	// This is mostly to ease testing; the default of 5s is too infrequent
	// to test conveniently. See the disk-stalled roachtest for an example
	// of how this is used.
	diskHealthCheckInterval := 5 * time.Second
	maxSyncDuration := env.EnvOrDefaultDuration("MAX_SYNC_DURATION", 60*time.Second)
	if diskHealthCheckInterval.Seconds() > maxSyncDuration.Seconds() {
		diskHealthCheckInterval = maxSyncDuration
	}

	// Instantiate a file system with disk health checking enabled. The
	// FS wraps vfs.Default, and can be wrapped for encryption-at-rest
	opts.FS = vfs.WithDiskHealthChecks(vfs.Default, diskHealthCheckInterval, func(s string, d time.Duration) {
		opts.EventListener.DiskSlow(pebble.DiskSlowInfo{
			Path:     s,
			Duration: d,
		})
	})

	return opts
}
