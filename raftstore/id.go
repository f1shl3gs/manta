package raftstore

import (
	"math"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	tsLen     = 5 * 8
	cntLen    = 8
	suffixLen = tsLen + cntLen
)

// idGenerator generates unique identifiers based on counters, timestamps, and
// a node member ID.
//
// The initial id is in this format:
// High order 2 bytes are from memberID, next 5 bytes are from timestamp,
// and low order one byte is a counter.
// | prefix   | suffix              |
// | 2 bytes  | 5 bytes   | 1 byte  |
// | memberID | timestamp | cnt     |
//
// The timestamp 5 bytes is different when the machine is restart
// after 1 ms and before 35 years.
//
// It increases suffix to generate the next id.
// The count field may overflow to timestamp field, which is intentional.
// It helps to extend the event window to 2^56. This doesn't break that
// id generated after restart is unique because etcd throughput is <<
// 256req/ms(250k reqs/second).
type idGenerator struct {
	// high order 2 bytes
	prefix uint64
	// low order 6 bytes
	suffix uint64
}

func newGenerator(memberID uint16, now time.Time) *idGenerator {
	prefix := uint64(memberID) << suffixLen
	unixMilli := uint64(now.UnixNano()) / uint64(time.Millisecond/time.Nanosecond)
	suffix := lowbit(unixMilli, tsLen) << cntLen

	return &idGenerator{
		prefix: prefix,
		suffix: suffix,
	}
}

// Next generates an id that is unique
func (g *idGenerator) Next() uint64 {
	suffix := atomic.AddUint64(&g.suffix, 1)
	return g.prefix | lowbit(suffix, suffixLen)
}

func lowbit(x uint64, n uint) uint64 {
	return x & (math.MaxUint64 >> (64 - n))
}
