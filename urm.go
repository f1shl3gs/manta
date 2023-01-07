package manta

import (
	"context"
	"errors"
)

var (
	ErrInvalidUserType = errors.New("unknown user type")

	ErrInvalidMappingType = errors.New("unknown mapping type")
)

// UserType can either be owner or member.
type UserType string

const (
	// Owner can read and write to a resource
	Owner UserType = "owner"
	// Member can read from a resource
	Member UserType = "member"
)

func (ut UserType) Valid() error {
	switch ut {
	case Owner, Member:
	default:
		return ErrInvalidUserType
	}

	return nil
}

type MappingType string

const (
	UserMappingType = "user"
	OrgMappingType  = "org"
)

func (mt MappingType) Valid() error {
	switch mt {
	case UserMappingType, OrgMappingType:
		return nil
	default:
		return ErrInvalidMappingType
	}
}

type UserResourceMapping struct {
	UserID       ID           `json:"userID"`
	UserType     UserType     `json:"userType"`
	MappingType  MappingType  `json:"mappingType"`
	ResourceType ResourceType `json:"resourceType"`
	ResourceID   ID           `json:"resourceID"`
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

func (m *UserResourceMapping) ownerPerms() ([]Permission, error) {
	if m.ResourceType == OrgsResourceType {
		return OwnerPermissions(m.ResourceID), nil
	}

	ps := []Permission{
		{
			Action: ReadAction,
			Resource: Resource{
				Type: m.ResourceType,
				ID:   &m.ResourceID,
			},
		},
		{
			Action: WriteAction,
			Resource: Resource{
				Type: m.ResourceType,
				ID:   &m.ResourceID,
			},
		},
	}

	return ps, nil
}

func (m *UserResourceMapping) memberPerms() ([]Permission, error) {
	if m.ResourceType == OrgsResourceType {
		return MemberPermissions(m.ResourceID), nil
	}

	ps := []Permission{
		{
			Action: ReadAction,
			Resource: Resource{
				Type: m.ResourceType,
				ID:   &m.ResourceID,
			},
		},
	}

	return ps, nil
}

// ToPermissions converts a user resource mapping into a set of permissions
func (m *UserResourceMapping) ToPermissions() ([]Permission, error) {
	switch m.UserType {
	case Owner:
		return m.ownerPerms()
	case Member:
		return m.memberPerms()
	default:
		return nil, ErrInvalidUserType
	}
}
