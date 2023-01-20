package kv

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricService struct {
	SchemaStore

	updateSec prometheus.Histogram
	viewSec   prometheus.Histogram
}

func NewMetricService(store SchemaStore, reg prometheus.Registerer) *MetricService {
	s := &MetricService{
		SchemaStore: store,
	}

	s.updateSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "manta",
		Subsystem: "kv",
		Name:      "update_duration_seconds",
		Help:      "The latency distributions of update called by kv",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	s.viewSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "manta",
		Subsystem: "kv",
		Name:      "view_duration_seconds",
		Help:      "The latency distributions of view called by kv",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	reg.MustRegister(s.updateSec, s.viewSec)

	return s
}

// View opens up a transaction that will not write to any data. Implementing interfaces
// should take care to ensure that all view transactions do not mutate any data.
func (s *MetricService) View(ctx context.Context, fn func(Tx) error) error {
	start := time.Now()
	err := s.SchemaStore.View(ctx, fn)
	s.viewSec.Observe(time.Since(start).Seconds())

	return err
}

// Update opens up a transaction that will mutate data.
func (s *MetricService) Update(ctx context.Context, fn func(Tx) error) error {
	start := time.Now()
	err := s.SchemaStore.Update(ctx, fn)
	s.updateSec.Observe(time.Since(start).Seconds())

	return err
}
