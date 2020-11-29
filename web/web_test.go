package web

import (
	"context"
	"os"
	"testing"

	"github.com/f1shl3gs/manta/bolt"
	"github.com/f1shl3gs/manta/kv"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func NewTestBackend(t *testing.T) (*Backend, *zap.Logger, func()) {
	logger := zaptest.NewLogger(t)

	store := bolt.NewKVStore(logger, t.Name(), bolt.WithNoSync)
	err := store.Open(context.Background())
	if err != nil {
		panic(err)
	}

	err = kv.Initial(context.Background(), store)
	require.NoError(t, err)

	service := kv.NewService(logger, store)
	backend := &Backend{
		OrganizationService: service,
	}

	return backend, logger, func() {
		err := store.Close()
		require.NoError(t, err)

		err = os.Remove(t.Name())
		require.NoError(t, err)
	}
}
