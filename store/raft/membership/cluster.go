package membership

import (
	"sync"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

var (
	// ErrIDExists is thrown when a node wants to join the existing cluster but its ID already exists
	ErrIDExists = errors.New("membership: can't add node to cluster, node id is a duplicate")
	// ErrIDRemoved is thrown when a node tries to perform an operation on an existing cluster but was removed
	ErrIDRemoved = errors.New("membership: node was removed during cluster lifetime")
	// ErrIDNotFound is thrown when we try an operation on a member that does not exist in the cluster list
	ErrIDNotFound = errors.New("membership: member not found in cluster list")
	// ErrConfigChangeInvalid is thrown when a configuration change we received looks invalid in form
	ErrConfigChangeInvalid = errors.New("membership: ConfChange type should be either AddNode, RemoveNode or UpdateNode")
	// ErrCannotUnmarshalConfig is thrown when a node cannot unmarshal a configuration change
	ErrCannotUnmarshalConfig = errors.New("membership: cannot unmarshal configuration change")
	// ErrMemberRemoved is thrown when a node was removed from the cluster
	ErrMemberRemoved = errors.New("raft: member was removed from the cluster")
)

type Attributes struct {
	Name string `json:"name,omitempty"`

	Addresses []string `json:"addresses,omitempty"`

	// Learner indicates if the member is raft learner.
	Learner bool `json:"learner,omitempty"`
}

type Member struct {
	ID uint64
	Attributes
}

type Cluster struct {
	mtx sync.RWMutex

	members map[uint64]*Member

	// removed contains the list of removed members,
	// those ids cannot be reused
	removed map[uint64]bool
}

func New() *Cluster {
	return &Cluster{
		members: make(map[uint64]*Member),
		removed: make(map[uint64]bool),
	}
}

func (c *Cluster) Members() map[uint64]*Member {
	members := make(map[uint64]*Member)

	c.mtx.RLock()
	defer c.mtx.RUnlock()

	for k, v := range c.members {
		members[k] = v
	}

	return members
}

// Removed returns the list of raft Members removed from the Cluster
func (c *Cluster) Removed() []uint64 {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	removed := make([]uint64, 0, len(c.removed))
	for k := range c.removed {
		removed = append(removed, k)
	}

	return removed
}

func (c *Cluster) Member(id uint64) *Member {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.members[id]
}

func (c *Cluster) AddMember(m *Member) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.removed[m.ID] {
		return ErrIDRemoved
	}

	c.members[m.ID] = m

	return nil
}

// RemoveMember removes a node from the Cluster's member list,
// and adds it to the removed list
func (c *Cluster) RemoveMember(id uint64) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.removed[id] = true

	delete(c.members, id)

	return nil
}

func (c *Cluster) UpdateMember(id uint64, m *Member) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.removed[id] {
		return ErrIDRemoved
	}

	_, ok := c.members[id]
	if !ok {
		return ErrIDNotFound
	}

	c.members[id] = m

	return nil
}

func (c *Cluster) ValidateConfigurationChange(cc raftpb.ConfChange) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.removed[cc.NodeID] {
		return ErrIDRemoved
	}

	switch cc.Type {
	case raftpb.ConfChangeAddNode:
		if c.members[cc.NodeID] != nil {
			return ErrIDExists
		}

	case raftpb.ConfChangeRemoveNode:
		if c.members[cc.NodeID] == nil {
			return ErrIDNotFound
		}

	case raftpb.ConfChangeUpdateNode:
		if c.members[cc.NodeID] == nil {
			return ErrIDNotFound
		}

	// todo: handle learner

	default:
		return ErrConfigChangeInvalid
	}

	return nil
}
