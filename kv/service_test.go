package kv_test

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/bolt"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/kv/migration"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestingT is a subset of the API provided by all *testing.T and *testing.B
// objects.
type TestingT interface {
	// Logs the given message without failing the test.
	Logf(string, ...interface{})

	// Logs the given message and marks the test as failed.
	Errorf(string, ...interface{})

	// Marks the test as failed.
	Fail()

	// Returns true if the test has been marked as failed.
	Failed() bool

	// Returns the name of the test.
	Name() string

	// Marks the test as failed and stops execution of that test.
	FailNow()
}

func NewTestBolt(t TestingT, noSync bool) (*bolt.KVStore, func()) {
	f, err := ioutil.TempFile("", "manta-test")
	require.NoError(t, err)
	f.Close()

	dbName := f.Name()
	logger := zaptest.NewLogger(t)
	var opts []bolt.KVOption
	if noSync {
		opts = append(opts, bolt.WithNoSync)
	}

	s := bolt.NewKVStore(logger, dbName, opts...)
	err = s.Open(context.TODO())
	require.NoError(t, err)

	return s, func() {
		s.Close()
		os.RemoveAll(dbName)
	}
}

func NewTestService(t TestingT, opts ...kv.Option) (*kv.Service, func()) {
	store, closer := NewTestBolt(t, true)

	svc := kv.NewService(zaptest.NewLogger(t), store, opts...)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	migrator := migration.New(zaptest.NewLogger(t), store, migration.All...)
	err := migrator.Up(ctx)
	require.NoError(t, err)

	return svc, closer
}

func CreateTestOrg(t TestingT, svc *kv.Service, name string) manta.ID {
	org := &manta.Organization{
		Name: name,
		Desc: name + " desc",
	}

	err := svc.CreateOrganization(context.Background(), org)
	require.NoError(t, err, "create test org %q failed", name)

	return org.ID
}

func CreateDefaultOrg(t TestingT, svc *kv.Service) manta.ID {
	return CreateTestOrg(t, svc, "test")
}
