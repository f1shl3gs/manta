package manta

import (
	"context"
	"encoding/json"
	"time"
)

var (
	ErrUserAlreadyExist = &Error{
		Code: EInvalid,
		Msg:  "user already exist",
	}

	ErrUserNotFound = &Error{
		Code: ENotFound,
		Msg:  "user not found",
	}
)

type User struct {
	ID      ID        `json:"id,omitempty"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Name    string    `json:"name,omitempty"`
	Status  string    `json:"status,omitempty"`
}

func (u *User) Unmarshal(data []byte) error {
	return json.Unmarshal(data, u)
}

func (u *User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

// UserFilter represents a set of filter that restrict the returned results.
type UserFilter struct {
	ID   *ID
	Name *string
}

type UserUpdate struct {
	Name *string
}

func (upd *UserUpdate) Apply(user *User) {
	if upd.Name != nil {
		user.Name = *upd.Name
	}
}

type UserService interface {
	FindUserByID(ctx context.Context, id ID) (*User, error)

	// return the first user that match the filter
	FindUser(ctx context.Context, filter UserFilter) (*User, error)

	// return a list of users and the total count of the matching user
	FindUsers(ctx context.Context, filter UserFilter, opts ...FindOptions) ([]*User, error)

	// CreateUser create a new user and set user.id with identifier
	CreateUser(ctx context.Context, user *User) error

	// Update a single user with changeset
	// Return the new User after update
	UpdateUser(ctx context.Context, id ID, udp UserUpdate) (*User, error)

	// Remove a user by ID
	DeleteUser(ctx context.Context, id ID) error
}
