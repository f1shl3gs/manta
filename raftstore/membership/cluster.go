package membership

import (
	"crypto/sha1"
	"encoding/binary"
	"sync"
	"time"
)

type Cluster struct {
	cid, localID uint64

	mtx     sync.RWMutex
	members map[uint64]*Member

	// removed contains the list of removed members
	// those ids connot be reused.
	//
	// TODO: what if host_1 is down for some reason, we remove it from cluster, then
	// host_1 is ready to serve, and we cannot add it to cluster?
	removed map[uint64]bool
}

func NewWithPeers(addrs []string) *Cluster {
	cl := &Cluster{
		members: make(map[uint64]*Member),
		removed: make(map[uint64]bool),
	}

	for _, addr := range addrs {
		id := generateID(addr)
		cl.members[id] = &Member{
			ID:      id,
			Address: addr,
			Learner: false,
		}
	}

	return cl
}

func New() *Cluster {
	return &Cluster{
		members: make(map[uint64]*Member),
		removed: make(map[uint64]bool),
	}
}

func GenerateID(input string) uint64 {
	b := []byte(input)
	b = append(b, []byte(time.Now().String())...)

	hash := sha1.Sum(b)

	return binary.BigEndian.Uint64(hash[:8])
}

func (c *Cluster) SetID(localID, cid uint64) {
	c.cid = cid
	c.localID = localID
}
