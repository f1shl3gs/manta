package mock

import "github.com/f1shl3gs/manta"

var _ manta.Authorizer = &Authorizer{}

// Authorizer is an Authorizer for testing that can allow everything or use specific permissions
type Authorizer struct {
	UserID      manta.ID
	Permissions []manta.Permission
	AllowAll    bool
}

func NewAuthorizer(allowAll bool, permissions []manta.Permission) *Authorizer {
	if allowAll {
		return &Authorizer{
			AllowAll: true,
		}
	}

	return &Authorizer{
		AllowAll:    false,
		Permissions: permissions,
	}
}

func (a *Authorizer) Identifier() manta.ID {
	return 1
}

func (a *Authorizer) GetUserID() manta.ID {
	if a.UserID.Valid() {
		return a.UserID
	}

	return 2
}

func (a *Authorizer) Kind() string {
	return "mock"
}

func (a *Authorizer) PermissionSet() (manta.PermissionSet, error) {
	if a.AllowAll {
		return manta.OperPermissions(), nil
	}

	return a.Permissions, nil
}
