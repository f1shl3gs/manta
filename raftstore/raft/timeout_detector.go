package raft

import (
	"sync"
	"time"
)

// TimeoutDetector detects routine starvations by
// observing the actual time duration to finish an action
// or between two events that should happen in a fixed
// interval. If the observed duration is longer than
// the expectation, the detector will report the result.
type TimeoutDetector struct {
	mu          sync.Mutex // protects all
	maxDuration time.Duration
	// map from event to time
	// time is the last seen time of the event.
	records map[uint64]time.Time
}

// NewTimeoutDetector creates the TimeoutDetector.
func NewTimeoutDetector(maxDuration time.Duration) *TimeoutDetector {
	return &TimeoutDetector{
		maxDuration: maxDuration,
		records:     make(map[uint64]time.Time),
	}
}

// Reset resets the NewTimeoutDetector.
func (td *TimeoutDetector) Reset() {
	td.mu.Lock()
	defer td.mu.Unlock()

	td.records = make(map[uint64]time.Time)
}

// Observe observes an event for given id. It returns false and exceeded duration
// if the interval is longer than the expectation.
func (td *TimeoutDetector) Observe(which uint64) (bool, time.Duration) {
	td.mu.Lock()
	defer td.mu.Unlock()

	ok := true
	now := time.Now()
	exceed := time.Duration(0)

	if pt, found := td.records[which]; found {
		exceed = now.Sub(pt) - td.maxDuration
		if exceed > 0 {
			ok = false
		}
	}
	td.records[which] = now
	return ok, exceed
}
