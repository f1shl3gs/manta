package manta

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	// ErrInvalidUserType notes that the provided UserType is invalid
	ErrInvalidUserType = errors.New("unknown user type")
	// ErrInvalidMappingType notes that the provided MappingType is invalid
	ErrInvalidMappingType = errors.New("unknown mapping type")
)

// UserType can either be owner or member.
type UserType string

const (
	// Owner can read and write to a resource
	Owner UserType = "owner" // 1
	// Member can read from a resource.
	Member UserType = "member" // 2
)

// Valid checks if the UserType is a member of the UserType enum
func (ut UserType) Valid() (err error) {
	switch ut {
	case Owner: // 1
	case Member: // 2
	default:
		err = ErrInvalidUserType
	}

	return err
}

//
type MappingType int32

const (
	UserMappingType = 0
	OrgMappingType  = 1
)

func (mt MappingType) Valid() error {
	switch mt {
	case UserMappingType, OrgMappingType:
		return nil
	}

	return ErrInvalidMappingType
}

func (mt MappingType) String() string {
	switch mt {
	case UserMappingType:
		return "user"
	case OrgMappingType:
		return "org"
	}

	return "unknown"
}

func (mt MappingType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.String())
}

func (mt *MappingType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case "user":
		*mt = UserMappingType
		return nil
	case "org":
		*mt = OrgMappingType
		return nil
	}

	return ErrInvalidMappingType
}

type UserResourceMappingFilter struct {
	ResourceID   ID
	ResourceType ResourceType
	UserID       ID
	UserType     UserType
}

// UserResourceMappingService maps the relationships between users and resources
type UserResourceMappingService interface {
	// FindUserResourceMappings returns a list of UserResourceMappings that match filter
	FindUserResourceMappings(ctx context.Context, filter UserResourceMappingFilter) ([]*UserResourceMapping, error)

	// CreateUserResourceMapping creates a user resource mapping
	CreateUserResourceMapping(ctx context.Context, m *UserResourceMapping) error

	// DeleteUserResourceMapping deletes a user resource mapping
	DeleteUserResourceMapping(ctx context.Context, userID, resourceID ID) error
}
