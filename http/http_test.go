package http

import (
	"context"
    "github.com/f1shl3gs/manta/telemetry/prom"
    "os"
	"testing"

	"github.com/f1shl3gs/manta/bolt"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/kv/migration"

	"github.com/f1shl3gs/manta/http/router"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func NewTestHTTPService(t *testing.T) (*Service, *Backend) {
	logger := zaptest.NewLogger(t)
	tmpFile := t.TempDir() + "/" + t.Name()
	store := bolt.NewKVStore(logger, tmpFile, bolt.WithNoSync)
	err := store.Open(context.Background())
	if err != nil {
		panic(err)
	}

	migrator := migration.New(logger, store, migration.All...)
	err = migrator.Up(context.Background())
	require.NoError(t, err)

	service := kv.NewService(logger, store)
	backend := &Backend{
		router:              router.New(),
		OrganizationService: service,
		UserService:         service,
		OnBoardingService:   service,
		PasswordService:     service,
		SessionService:      service,
        PromRegistry: prom.NewRegistry(logger),
	}

	// This is very trick, this will deletedashboard the data file, and
	// you will not find it by something like `ls` or `find`.
	// While the bolt will be still working, and the dta file
	// will be deleted when this test process exit.
	_ = os.RemoveAll(tmpFile)

	return New(logger, backend), backend
}
