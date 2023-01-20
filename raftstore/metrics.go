package raftstore

import (
	"github.com/f1shl3gs/manta/bolt"

	"github.com/prometheus/client_golang/prometheus"
)

type boltCollector func(chan<- prometheus.Metric)

func (c boltCollector) Describe(chan<- *prometheus.Desc) {}

func (c boltCollector) Collect(ch chan<- prometheus.Metric) {
	c(ch)
}

func (s *Store) Collectors() []prometheus.Collector {
	// trick of closure
	var bc boltCollector = func(ch chan<- prometheus.Metric) {
		db := s.db.Load()
		if db == nil {
			return
		}

		c := bolt.NewCollector(db)
		c.Collect(ch)
	}

	return []prometheus.Collector{
		bc,

		s.leaderChanges,
		s.hasLeader,
		s.isLeader,
		s.slowReadInex,
		s.readIndexFailed,
	}
}
