package mvcc

import bolt "go.etcd.io/bbolt"

type ReadTx interface {
	Get(bucket, key []byte, ver int64) ([]byte, int64, error)

	Range(bucket, start, end []byte, rev, limit int64) ([][]byte, [][]byte, error)
}

type WriteTxn interface {
	ReadTx

	CreateBucket(name []byte) error
	DeleteBcuket(name []byte) error

	Put(key, value []byte, rev int64) error

	// Delete delete the given key.
	// The returned rev is the current revision of the KV when the operation is excuted.
	Delete(key []byte) (int64, error)
}

type Store interface {
	Read() ReadTx
	Write() WriteTx

	Restore(r io.Reader) error
}

type store struct {
	db *bolt.DB
}

func (s *store) Read() ReadTx {

}

func (s *store) Write() WriteTx {

}

func (s *store) Restore(r io.Reader) error {

}
