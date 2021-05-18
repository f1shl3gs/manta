package membership

import (
	"github.com/f1shl3gs/manta"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

type Attributes struct {
	Name string `json:"name,omitempty"`

	Addresses []string `json:"addresses,omitempty"`

	// Learner indicates if the member is raft learner.
	Learner bool `json:"learner,omitempty"`
}

type Member struct {
	ID manta.ID `json:"id"`
	Attributes
}

// Cluster is an interface representing a collection of members in one cluster
type Cluster interface {
	// ID returns the cluster ID
	ID() manta.ID

	// Members returns a slice of members sorted by their ID
	Members() []*Member

	// Member retrieves a particular member based on ID, or nil if
	// the member does not exist in the cluster
	Member(id manta.ID) *Member

	// UpdateAttributes
	UpdateAttributes(id manta.ID, attr Attributes)

	// ValidateConfigurationChange takes a proposed ConfChange and
	// ensures that it is still valid
	ValidateConfigurationChange(cc raftpb.ConfChange) error

	// AddMember adds a new Member into the cluster, and saves the
	// given member's raftAttributes into the store. The given member
	// should have empty attributes.
	// A Member with a matching id must not exist
	AddMember(m *Member)

	// PromoteMember marks the member's IsLearner RaftAttributes to false
	PromoteMember(id uint64)

	// RemoveMember removes a member from the store
	// The given id MUST exist, or the function panics
	RemoveMember(id uint64)
}

// ConfigChangeContext represents a context for confChange.
type ConfigChangeContext struct {
	Member
	// IsPromote indicates if the config change is for promoting a learner member.
	// This flag is needed because both adding a new member and promoting a learner member
	// uses the same config change type 'ConfChangeAddNode'.
	IsPromote bool `json:"isPromote"`
}
