package member

import "context"

type Role string

const (
	Checker = "checker"
	Mantad  = "manta"
)

type Member struct {
	Addr string `json:"address"`
	Role Role   `json:"role"`
}

type Membership interface {
	Join(ctx context.Context, peer string) error

	Leave(ctx context.Context) error

	Members(ctx context.Context) ([]*Member, error)
}
