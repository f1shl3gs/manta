package manta

import (
	"context"
)

var (
	ErrOtclNotFound = &Error{
		Code: ENotFound,
		Msg:  "otcls not found",
	}
)

type OtclFilter struct {
	OrgID *ID
	Type  *string
}

type OtclPatch struct {
	Name        *string
	Description *string
	Content     *string
}

type OtclService interface {
	FindOtclByID(ctx context.Context, id ID) (*Otcl, error)

	FindOtcls(ctx context.Context, filter OtclFilter) ([]*Otcl, error)

	CreateOtcl(ctx context.Context, o *Otcl) error

	PatchOtcl(ctx context.Context, id ID, u OtclPatch) (*Otcl, error)

	DeleteOtcl(ctx context.Context, id ID) error
}
