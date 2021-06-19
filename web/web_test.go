package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/bolt"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/kv/migration"

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

	migrator := migration.New(logger, store, migration.All...)
	err = migrator.Up(context.Background())
	require.NoError(t, err)

	service := kv.NewService(logger, store)
	backend := &Backend{
		OrganizationService: service,
		SecretService:       service,
	}

	return backend, logger, func() {
		err := store.Close()
		require.NoError(t, err)

		err = os.Remove(t.Name())
		require.NoError(t, err)
	}
}

func NewTestBackendWithOrg(t *testing.T) (*Backend, *zap.Logger, manta.ID, func()) {
	backend, logger, closer := NewTestBackend(t)
	org := &manta.Organization{
		Name: "test",
		Desc: "test",
	}

	err := backend.OrganizationService.CreateOrganization(context.Background(), org)
	require.NoError(t, err)
	return backend, logger, org.ID, closer
}

func TestDemoServer(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(data))

		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":8080", nil)
}
