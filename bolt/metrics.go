package bolt

import (
	"github.com/prometheus/client_golang/prometheus"
	bolt "go.etcd.io/bbolt"
)

var _ prometheus.Collector = (*KVStore)(nil)

var (
	boltWritesDesc = prometheus.NewDesc(
		"boltdb_writes_total",
		"Total number of boltdb writes",
		nil, nil)

	boltReadsDesc = prometheus.NewDesc(
		"boltdb_reads_total",
		"Total number of boltdb reads",
		nil, nil)

	boltBucketKeysDesc = prometheus.NewDesc(
		"boltdb_keys_totoal",
		"Total number of keys of the bucket",
		[]string{"bucket"},
		nil,
	)
)

// Describe returns all descriptions of the collector.
func (c *KVStore) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("mantad_boltdb", "", nil, nil)
}

// Collect returns the current state of all metrics of the collector.
func (c *KVStore) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()
	writes := stats.TxStats.Write
	reads := stats.TxN

	ch <- prometheus.MustNewConstMetric(
		boltReadsDesc,
		prometheus.CounterValue,
		float64(reads),
	)

	ch <- prometheus.MustNewConstMetric(
		boltWritesDesc,
		prometheus.CounterValue,
		float64(writes),
	)

	keys := make(map[string]int)
	_ = c.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			keys[string(name)] = b.Stats().KeyN
			return nil
		})
	})

	for key, n := range keys {
		ch <- prometheus.MustNewConstMetric(
			boltBucketKeysDesc,
			prometheus.GaugeValue,
			float64(n),
			key,
		)
	}

	c.commitSec.Collect(ch)
	c.writeSec.Collect(ch)
	c.updateSec.Collect(ch)
	c.viewSec.Collect(ch)
}
