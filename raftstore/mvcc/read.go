package mvcc

import (
	bolt "go.etcd.io/bbolt"

	"github.com/f1shl3gs/manta/kv"
)

type readTx struct {
	tx *bolt.Tx
}

func (tx *readTx) Get(bucket, key []byte, ver int64) ([]byte, int64, error) {
	b := tx.tx.Bucket(bucket)
	if b == nil {
		return nil, 0, kv.ErrBucketNotFound
	}

	c := b.Cursor()

}

func (tx *readTx) Range(bucket, start, end []byte, rev, limit int64) ([][]byte, [][]byte, error) {

}
