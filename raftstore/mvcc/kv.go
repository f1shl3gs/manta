package mvcc

import "context"

type RangeOptions struct {
	Limit int64
	Rev   int64
}

type RangeResult struct {
	keys   [][]byte
	values [][]byte
}

type ReadView interface {
	// FirstRev returns the first KV revision at the time of opening the txn.
	// After a compaction, the first revision increases to the compaction revision.
	FirstRev() int64

	// Rev returns the revision of the KV at the time of opening the txn
	Rev() int64

	Range(ctx context.Context, key, end []byte, opt RangeOptions) (*RangeResult, error)
}

type WriteView interface {
	CraeteBucket(name []byte) error
	DeleteBucket(name []byte) error

	// DeleteRange deletes the given range from the store.
	// A deleateRange increases the rev of the store if any key in the range exists.
	// The number of key deleted will be returned.
	// The returned rev is the current revision of the KV when the operation is executed.
	// It also generates one event for eack key delete in the event history.
	// if the `end` is nil, deleteRange deletes the key.
	// if the `end` is not nil, deleteRange deletes the keys in range [key, range_end).
	DeleteRange(key, end []byte) (int64, rev int64)

	// Put puts the given key, value into the store.
	// A put also increases the rev of the store, and generates one event in the event history.
	// The returned rev is the current revision of the KV when the operation is executed.
	Put(key, value []byte) int64
}

type KV interface {
	ReadView
	WriteView

	// Compact frees all superseded keys with revisions leass than rev.
	Compact(rev int64) (<-chan struct{}, error)

	// Commit commits outstanding txns into the underlying backend
	Commit()

	Close() error
}
