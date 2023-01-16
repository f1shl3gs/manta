package membership

import (
	"crypto/sha1"
	"encoding/binary"
    "encoding/json"
    "sync"
	"time"
)

type Member struct {
	ID uint64 `json:"id"`

	// Addresses is the list of peers in the raft cluster
	Address string `json:"address"`

	// Learner indicates if the member is raft learner
	Learner bool `json:"learner,omitempty"`
}

type Cluster struct {
	cid, localID uint64

	mtx     sync.RWMutex
	members map[uint64]*Member
	// removed contains the list of removed members
	// those ids connot be reused.
	removed map[uint64]bool
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

func (c *Cluster) Members() []Member {
    c.mtx.RLock()
    defer c.mtx.RUnlock()

    var ms []Member
    for _, m := range c.members {
        ms = append(ms, *m)
    }

    return ms
}

func (c *Cluster) Add(m *Member) {
    c.mtx.Lock()
    c.members[m.ID] = m
    c.mtx.Unlock()
}

func (c *Cluster) Remove(id uint64) {
    c.mtx.Lock()
    delete(c.members, id)
    c.removed[id] = true
    c.mtx.Unlock()
}

type state struct {
    Members map[uint64]*Member
    // removed contains the list of removed members
    // those ids connot be reused.
    Removed map[uint64]bool
}

func (c *Cluster) Marshal() []byte {
    c.mtx.RLock()
    defer c.mtx.RUnlock()

    data, err := json.Marshal(&state{
        Members: c.members,
        Removed: c.removed,
    })
    if err != nil {
        panic(err)
    }

    return data
}
