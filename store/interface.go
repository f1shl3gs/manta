package store

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
	TenantStorage(ctx context.Context, id manta.ID) (storage.Storage, error)
}
