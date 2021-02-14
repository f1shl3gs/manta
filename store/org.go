package store

import (
	"context"
	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/cache"
	"github.com/prometheus/prometheus/storage"
)

type org struct {
	ts TenantStorage

	queryables cache.Cache
	appendable cache.Cache
}

func (o *org) Queryable(ctx context.Context, id manta.ID) (storage.Queryable, error) {
	panic("implement me")
}

func (o *org) Appendable(ctx context.Context, id manta.ID) (storage.Appendable, error) {
	panic("implement me")
}
