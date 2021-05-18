package state

import "github.com/prometheus/client_golang/prometheus"

type collectors []prometheus.Collector

func (cs collectors) Describe(descs chan<- *prometheus.Desc) {
	for _, c := range cs {
		c.Describe(descs)
	}
}

func (cs collectors) Collect(metrics chan<- prometheus.Metric) {
	for _, c := range cs {
		c.Collect(metrics)
	}
}

func (s *Server) PromCollector() prometheus.Collector {
	return collectors{
		s.applySnapshotInProgress,
		s.proposalsApplied,
		s.isLearner,
	}
}
