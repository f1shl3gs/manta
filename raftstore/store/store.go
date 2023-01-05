package store

type Bucket interface {
    Get(key []byte, rev int64) ([]byte, error)
	Put(key, value []byte, rev int64) error
    Range(start, end []byte, rev, limit int64) ([][]byte, [][]byte, error)
}

type Store interface {
    CreateBucket(name []byte) error
    DeleteBcuket(name []byte) error

	Bucket(name []byte) (Bucket, error)
}

type store struct {

}
