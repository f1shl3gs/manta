package manta

import (
	"context"
)

type NodeFilter struct {
	Address *string
	Env     *string
	Search  *string
	OrgID   *ID
}

type NodeUpdate struct {
	Env      *string
	Hostname *string

	AnnotationOverwrite bool
	Annotations         map[string]string
}

type NodeService interface {
	// FindNodeByID returns a single Node by id
	FindNodeByID(ctx context.Context, id ID) (*Node, error)

	// FindNodes returns a list of Nodes that match filter and the total count of matching Nodes
	// additional options provide pagination & sorting
	FindNodes(ctx context.Context, filter NodeFilter, opt ...FindOptions) ([]*Node, int, error)

	// CreateNode create a Node and set it's id with identifier
	CreateNode(ctx context.Context, svr *Node) error

	// UpdateNode update a single Node with changeset
	// returns the new Node after update
	UpdateNode(ctx context.Context, id ID, u NodeUpdate) (*Node, error)

	// DeleteNode delete a single Node by ID
	DeleteNode(ctx context.Context, id ID) error
}
