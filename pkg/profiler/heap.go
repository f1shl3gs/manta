package profiler

import (
	"context"
	"os"
	"runtime/pprof"
	"time"
)

// HeapFileNamePrefix is the prefix of files containing pprof data.
const HeapFileNamePrefix = "memprof"

// HeapFileNameSuffix is the suffix of files containing pprof data.
const HeapFileNameSuffix = ".pprof"

// HeapProfiler is used to take Go heap profiles
type HeapProfiler struct {
	*profiler
}

func NewHeapProfiler(dir string, maxSize uint64) *HeapProfiler {
	if maxSize == 0 {
		maxSize = defaultMaxSize
	}

	return &HeapProfiler{
		profiler: &profiler{
			dumpStore: dumpStore{
				maxSize: maxSize,
				dir:     dir,
				prefix:  HeapFileNamePrefix,
				suffix:  HeapFileNameSuffix,
			},
			lastProfileTime:    time.Time{},
			highWaterMarkBytes: 0,
			now:                time.Now,
		},
	}
}

func (p *HeapProfiler) MaybeTakeProfile(ctx context.Context, thresholdValue uint64) error {
	if resetHighWaterMarkInterval == 0 {
		// Instruction to disable
		return nil
	}

	now := p.now()
	// If it's been too long since we took a prefile, make sure we'll
	// take one now
	if now.Sub(p.lastProfileTime) >= resetHighWaterMarkInterval {
		p.highWaterMarkBytes = 0
	}

	takeProfile := thresholdValue > p.highWaterMarkBytes
	if !takeProfile {
		return nil
	}

	p.highWaterMarkBytes = thresholdValue
	p.lastProfileTime = now

	err := takeHeapProfile(ctx, p.format(now, thresholdValue))
	if err != nil {
		return err
	}

	if err = p.GC(); err != nil {
		return err
	}

	return nil
}

// takeHeapProfile returns true if and only if the prefile dump
// was taken successfully
func takeHeapProfile(ctx context.Context, path string) error {
	// Try writing a go heap profile.
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if err = pprof.WriteHeapProfile(f); err != nil {
		return err
	}

	return nil
}
