package tsdb

import (
	"context"
	"errors"

	"github.com/f1shl3gs/manta"
	"github.com/prometheus/prometheus/storage"
)

var (
	ErrUnknownTenantStorage = errors.New("unknown tenant storage")
)

type TenantStorage interface {
	Queryable(ctx context.Context, id manta.ID) (storage.Queryable, error)

	Appendable(ctx context.Context, id manta.ID) (storage.Appendable, error)
}
