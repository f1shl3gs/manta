package bolt_test

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"

    "github.com/f1shl3gs/manta/bolt"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/tests"
	"go.uber.org/zap/zaptest"
)

func NewTestKVStore(t *testing.T) (*bolt.KVStore, func(), error) {
	f, err := ioutil.TempFile("", "manta-bolt")
	if err != nil {
		return nil, nil, errors.New("unable to open temporary boltdb file")
	}
	f.Close()

	path := f.Name()
	s := bolt.NewKVStore(zaptest.NewLogger(t), path, bolt.WithNoSync)
	if err := s.Open(context.TODO()); err != nil {
		return nil, nil, err
	}

	close := func() {
		s.Close()
		os.Remove(path)
	}

	return s, close, nil
}

func initKVStore(f tests.KVStoreFields, t *testing.T) (kv.Store, func()) {
	s, closeFn, err := NewTestKVStore(t)
	if err != nil {
		t.Fatalf("failed to create new kv store: %v", err)
	}

	mustCreateBucket(t, s, f.Bucket)

	err = s.Update(context.Background(), func(tx kv.Tx) error {
		b, err := tx.Bucket(f.Bucket)
		if err != nil {
			return err
		}

		for _, p := range f.Pairs {
			if err := b.Put(p.Key, p.Value); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		t.Fatalf("failed to put keys: %v", err)
	}
	return s, func() {
		closeFn()
	}
}

func TestKVStore(t *testing.T) {
	tests.KVStore(initKVStore, t)
}

func mustCreateBucket(t testing.TB, store kv.SchemaStore, bucket []byte) {
	t.Helper()

	if err := store.CreateBucket(context.Background(), bucket); err != nil {
		t.Fatal(err)
	}
}
