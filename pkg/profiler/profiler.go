package profiler

import (
	"context"
	"fmt"
	"os"
	"time"
)

const timestampFormat = "2006-01-02T15_04_05.000"

// resetHighWaterMarkInterval specifies how often the high-water mark value will
// be reset. Immediately after it is reset, a new profile will be taken.
//
// If the value is 0, the collection of profiles gets disabled.
var resetHighWaterMarkInterval = func() time.Duration {
	text := os.Getenv("MEMPROF_INTERVAL")
	interval, err := time.ParseDuration(text)
	if err != nil {
		return time.Hour
	}

	if interval <= 0 {
		// Instruction to disable.
		return 0
	}
	return interval
}()

type profiler struct {
	dumpStore

	// lastProfileTime marks the time when we took the last profile
	lastProfileTime time.Time

	// highWaterMarkBytes represents the maximum heap size that we've
	// seen since resetting the filed(which happens periodically)
	highWaterMarkBytes uint64

	now func() time.Time
}

func (p *profiler) maybeTakeProfile(
	ctx context.Context,
	thresholdValue uint64,
	takeProfileFn func(ctx context.Context, path string) error,
) (bool, error) {
	if resetHighWaterMarkInterval == 0 {
		return false, nil
	}

	now := p.now()
	// If it's been too long since we took a profile, make sure we'll take one now
	if now.Sub(p.lastProfileTime) >= resetHighWaterMarkInterval {
		p.highWaterMarkBytes = 0
	}

	takeProfile := thresholdValue > p.highWaterMarkBytes
	if !takeProfile {
		return false, nil
	}

	p.highWaterMarkBytes = thresholdValue
	p.lastProfileTime = now

	dumpfile := fmt.Sprintf("%s.%s.%d%s",
		p.prefix, now.Format(timestampFormat), thresholdValue, p.suffix)

	return true, takeProfileFn(ctx, dumpfile)
}
