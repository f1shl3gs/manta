package multitsdb

import (
	"context"

	"github.com/f1shl3gs/manta"
	"github.com/prometheus/prometheus/storage"
)

type TenantStorage interface {
	Queryable(ctx context.Context, id manta.ID) (storage.Queryable, error)

	Appendable(ctx context.Context, id manta.ID) (storage.Appendable, error)
}

type Noop struct{}

func (n *Noop) Queryable(ctx context.Context, id manta.ID) (storage.Queryable, error) {
	return nil, ErrNotReady
}

func (n *Noop) Appendable(ctx context.Context, id manta.ID) (storage.Appendable, error) {
	return nil, ErrNotReady
}
