package membership

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.etcd.io/etcd/raft/v3/raftpb"
)

const maxLearners = 1

var (
	ErrIDRemoved        = errors.New("membership: ID removed")
	ErrIDExists         = errors.New("membership: ID exists")
	ErrIDNotFound       = errors.New("membership: ID not found")
	ErrPeerURLexists    = errors.New("membership: peerURL exists")
	ErrMemberNotLearner = errors.New("membership: can only promote a learner member")
	ErrTooManyLearners  = errors.New("membership: too many learner members in cluster")
)

// Cluster is an interface representing a collection of members in cluster
type Cluster interface {
	// ID returns the cluster ID
	ID() uint64

	// SetID set the cluster ID
	SetID(id uint64)

	// Members returns a slice of members sorted by their ID
	Members() []*Member

	// Member retrieves a particular member based on ID, or nil if the
	// member does not exist in the cluster
	Member(id uint64) *Member

	RemoveMember(id uint64)

	UpdateMember(id uint64, addr string)

	AddMember(m *Member)

	PromoteMember(id uint64)

	IDRemoved(id uint64) bool

	ValidateConfigurationChange(cc raftpb.ConfChange) error
}

type cluster struct {
	id      uint64
	mtx     sync.RWMutex
	members map[uint64]*Member

	// removed contains the list of removed Members
	// those ids cannot be reused
	removed map[uint64]bool

	// TODO: do we need a storage for members and removed too?
}

// ID returns the raft cluster ID, mtx is not needed
// cluster ID will not changed on the fly
func (c *cluster) ID() uint64 {
	return c.id
}

// SetID sets the raft cluster ID, mtx is not needed
// cluster ID will not changed on the fly
func (c *cluster) SetID(id uint64) {
	c.id = id
}

func (c *cluster) Members() []*Member {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	ms := make([]*Member, 0, len(c.members))
	for _, m := range c.members {
		ms = append(ms, m)
	}

	return ms
}

func (c *cluster) Member(id uint64) *Member {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.members[id]
}

// RemoveMember removes a member from the cluster and store?
// The given id MUST exist
func (c *cluster) RemoveMember(id uint64) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.members, id)
}

// UpdateMember updates a member's addr
func (c *cluster) UpdateMember(id uint64, addr string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	m, exists := c.members[id]
	if !exists {
		return
	}

	m.Address = addr
}

func (c *cluster) AddMember(m *Member) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.members[m.ID] = m
}

func (c *cluster) PromoteMember(id uint64) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	m := c.members[id]
	if m == nil {
		return
	}

	m.Learner = false
}

func (c *cluster) IDRemoved(id uint64) bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.removed[id]
}

func NewClusterFromAddresses(addrs []string) Cluster {
	cls := &cluster{
		id:      generateID("manta"),
		members: make(map[uint64]*Member, len(addrs)),
		removed: map[uint64]bool{},
	}

	for _, addr := range addrs {
		id := generateID(addr)
		m := &Member{ID: id}
		m.Address = addr

		cls.members[id] = m
	}

	return cls
}

func NewCluster() Cluster {
	return &cluster{
		members: make(map[uint64]*Member),
	}
}

func generateID(input string) uint64 {
	b := []byte(input)
	b = append(b, []byte(time.Now().String())...)

	hash := sha1.Sum(b)

	return binary.BigEndian.Uint64(hash[:8])
}

// ConfigChangeContext represents a context for confChange.
type ConfigChangeContext struct {
	Member
	// IsPromote indicates if the config change is for promoting a learner member.
	// This flag is needed because both adding a new member and promoting a learner member
	// uses the same config change type 'ConfChangeAddNode'.
	IsPromote bool `json:"isPromote"`
}

// ValidateConfigurationChange takes a proposed ConfChange and
// ensures that it is still valied
func (c *cluster) ValidateConfigurationChange(cc raftpb.ConfChange) error {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	id := cc.NodeID

	if c.removed[id] {
		return ErrIDRemoved
	}

	switch cc.Type {
	case raftpb.ConfChangeAddNode, raftpb.ConfChangeAddLearnerNode:
		confChangeContext := &ConfigChangeContext{}
		if err := json.Unmarshal(cc.Context, confChangeContext); err != nil {
			fmt.Println(string(cc.Context))
			panic("failed to unmarshal confChangeContext, err: " + err.Error())
		}

		if confChangeContext.IsPromote {
			// promoting a learner member to voting member
			if c.members[id] == nil {
				return ErrIDNotFound
			}

			if !c.members[id].Learner {
				return ErrMemberNotLearner
			}
		} else {
			// adding a new member
			if c.members[id] != nil {
				return ErrIDExists
			}

			for _, m := range c.members {
				if m.Address == confChangeContext.Address {
					return ErrPeerURLexists
				}
			}

			if confChangeContext.Member.Learner {
				// the new member is a learner
				numLearners := 0
				for _, m := range c.members {
					if m.Learner {
						numLearners += 1
					}
				}

				if numLearners+1 > maxLearners {
					return ErrTooManyLearners
				}
			}
		}

	case raftpb.ConfChangeRemoveNode:
		if c.members[id] == nil {
			return ErrIDNotFound
		}

	case raftpb.ConfChangeUpdateNode:
		if c.members[id] == nil {
			return ErrIDNotFound
		}

		addrs := make(map[string]bool)
		for _, m := range c.members {
			if m.ID == id {
				continue
			}

			addrs[m.Address] = true
		}

		m := &Member{}
		if err := json.Unmarshal(cc.Context, m); err != nil {
			panic("failed to unmarshal member, err: " + err.Error())
		}

		if addrs[m.Address] {
			return ErrPeerURLexists
		}

	default:
		panic("unknown ConfChange type " + cc.Type.String())
	}

	return nil
}
