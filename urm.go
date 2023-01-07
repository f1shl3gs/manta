package manta

import "context"

// UserType can either be owner or member.
type UserType string

const (
	// Owner can read and write to a resource
	Owner UserType = "owner"
	// Member can read from a resource
	Member UserType = "member"
)

type UserResourceMapping struct {
	UserID       ID           `json:"userID"`
	UserType     UserType     `json:"userType"`
	ResourceType ResourceType `json:"resourceType"`
}

type UserResourceMappingFilter struct {
	ResourceID   ID
	ResourceType ResourceType
	UserID       ID
	UserType     UserType
}

// UserResourceMappingService maps the relationships between users and resources
type UserResourceMappingService interface {
	// FindUserResourceMappings returns a list of UserResourceMappings that match filter and the total count of matching mappings.
	FindUserResourceMappings(ctx context.Context, filter UserResourceMappingFilter, opt ...FindOptions) ([]*UserResourceMapping, int, error)

	// CreateUserResourceMapping creates a user resource mapping.
	CreateUserResourceMapping(ctx context.Context, m *UserResourceMapping) error

	// DeleteUserResourceMapping deletes a user resource mapping.
	DeleteUserResourceMapping(ctx context.Context, resourceID, userID ID) error
}
