package manta

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrAuthorizerNotSupported notes that the provided authorizer is not supported for the action you are trying to perform.
	ErrAuthorizerNotSupported = errors.New("your authorizer is not supported, please use *platform.Authorization as authorizer")
	// ErrInvalidResourceType notes that the provided resource is invalid
	ErrInvalidResourceType = errors.New("unknown resource type for permission")
	// ErrInvalidAction notes that the provided action is invalid
	ErrInvalidAction = errors.New("unknown action for permission")
)

type UpdateAuthorization struct {
	Token  *string
	Status *string
}

type AuthorizationService interface {
	// FindAuthorizationByID returns a single authorization by id
	FindAuthorizationByID(ctx context.Context, id ID) (*Authorization, error)

	// FindAuthorizationByToken returns a single authorization by token
	FindAuthorizationByToken(ctx context.Context, token string) (*Authorization, error)

	// CreateAuthorization create a new Authorization and sets a.UserID and a.Token
	CreateAuthorization(ctx context.Context, a *Authorization) error

	// UpdateAuthorization updates the status and token if available
	UpdateAuthorization(ctx context.Context, id ID, u UpdateAuthorization) (*Authorization, error)

	// DeleteAuthorization delete a authorization by ID
	DeleteAuthorization(ctx context.Context, id ID) error
}

// ResourceType is an enum defining all resource types that have a permission model in platform
type ResourceType string

const (
	AuthorizationsResourceType = ResourceType("authorizations")
	DashboardsResourceType     = ResourceType("dashboards")
	OrganizationsResourceType  = ResourceType("organizations")
	OtclResourceType           = ResourceType("otcl")
	UsersResourceType          = ResourceType("users")
)

var AllResourceTypes = []ResourceType{
	AuthorizationsResourceType,
	DashboardsResourceType,
	OrganizationsResourceType,
	OtclResourceType,
	UsersResourceType,
}

// Action is an enum defining all possible resource operations
type Action string

const (
	// ReadAction is the action for reading.
	ReadAction Action = "read" // 1
	// WriteAction is the action for writing.
	WriteAction Action = "write" // 2
)

var actions = []Action{
	ReadAction,  // 1
	WriteAction, // 2
}

// Valid checks if the action is a member of the Action enum
func (a Action) Valid() (err error) {
	switch a {
	case ReadAction: // 1
	case WriteAction: // 2
	default:
		err = ErrInvalidAction
	}

	return err
}

type PermissionSet []Permission

func (ps PermissionSet) Allowed(p Permission) bool {
	return PermissionAllowed(p, ps)
}

// Matches returns whether or not one permission matches the other.
func (p Permission) Matches(perm Permission) bool {
	return p.matches(perm)
}

func (p Permission) matches(perm Permission) bool {
	if p.Action != perm.Action {
		return false
	}

	if p.Resource.Type != perm.Resource.Type {
		return false
	}

	if p.Resource.OrgID == nil && p.Resource.ID == nil {
		return true
	}

	if p.Resource.OrgID != nil && perm.Resource.OrgID != nil && p.Resource.ID != nil && perm.Resource.ID != nil {
		if p.Resource.OrgID != perm.Resource.OrgID && p.Resource.ID == perm.Resource.ID {
			fmt.Printf("match used: p.Resource.OrgID=%s perm.Resource.OrgID=%s p.Resource.ID=%s",
				p.Resource.OrgID, perm.Resource.OrgID, p.Resource.ID)
		}
	}

	if p.Resource.OrgID != nil {
		if perm.Resource.OrgID != nil {
			if p.Resource.OrgID == perm.Resource.OrgID {
				if p.Resource.ID == nil {
					return true
				}
				if perm.Resource.ID != nil {
					return p.Resource.ID == perm.Resource.ID
				}
			}
			return false
		}
	}

	if p.Resource.ID != nil {
		pID := p.Resource.ID
		if perm.Resource.ID != nil {
			permID := perm.Resource.ID
			if pID == permID {
				return true
			}
		}
	}

	return false
}

// PermissionAllowed determines if a permission is allowed.
func PermissionAllowed(perm Permission, ps []Permission) bool {
	for _, p := range ps {
		if p.Matches(perm) {
			return true
		}
	}
	return false
}

type Authorizer interface {
	Identifier() ID

	GetUserID() ID

	Kind() string

	PermissionSet() PermissionSet
}

func (a *Authorization) Identifier() ID {
	return a.ID
}

func (a *Authorization) GetUserID() ID {
	return a.UID
}

func (a *Authorization) Kind() string {
	return "auth"
}

func (a *Authorization) PermissionSet() PermissionSet {
	return a.Permissions
}

// OwnerPermissions are the default permissions for those who own a resource
func OwnerPermissions(orgID ID) []Permission {
	var ps []Permission

	for _, r := range AllResourceTypes {
		for _, a := range actions {
			if r == OrganizationsResourceType {
				ps = append(ps, Permission{
					Action: a,
					Resource: Resource{
						Type: r,
						ID:   &orgID,
					},
				})

				continue
			}

			ps = append(ps, Permission{
				Action: a,
				Resource: Resource{
					Type:  r,
					OrgID: &orgID,
				},
			})
		}
	}

	return ps
}

// MemberPermissions are the default permissions for those who can see a resource
func MemberPermissions(orgID ID) []Permission {
	var ps []Permission

	for _, r := range AllResourceTypes {
		if r == OrganizationsResourceType {
			ps = append(ps, Permission{
				Action: ReadAction,
				Resource: Resource{
					Type: r,
					ID:   &orgID,
				},
			})
			continue
		}

		ps = append(ps, Permission{
			Action: ReadAction,
			Resource: Resource{
				Type:  r,
				OrgID: &orgID,
			},
		})
	}

	return ps
}

// MePermissions is the permission to read/write user itself
func MePermissions(userID ID) []Permission {
	var ps []Permission

	for _, a := range actions {
		ps = append(ps, Permission{
			Action: a,
			Resource: Resource{
				Type: UsersResourceType,
				ID:   &userID,
			},
		})
	}

	return ps
}
