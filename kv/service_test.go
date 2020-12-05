package kv_test

import (
	"context"
	"github.com/f1shl3gs/manta/bolt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"os"
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
	dbName := t.Name() + ".bolt"
	logger := zaptest.NewLogger(t)
	var opts []bolt.KVOption
	if noSync {
		opts = append(opts, bolt.WithNoSync)
	}

	s := bolt.NewKVStore(logger, dbName, opts...)
	err := s.Open(context.TODO())
	require.NoError(t, err)

	return s, func() {
		s.Close()
		os.Remove(dbName)
	}
}
