package kv

import (
	"bytes"
	"context"

	"github.com/f1shl3gs/manta"
	"github.com/pkg/errors"
)

func IndexKey(foreignKey, primaryKey []byte) []byte {
	key := make([]byte, len(foreignKey)+len(primaryKey)+1)

	copy(key, foreignKey)
	key[len(foreignKey)] = '/'
	copy(key[len(foreignKey)+1:], primaryKey)

	return key
}

func indexKeyParts(indexKey []byte) (fk, pk []byte, err error) {
	// this function is called with items missing in index
	parts := bytes.SplitN(indexKey, []byte("/"), 2)
	if len(parts) < 2 {
		return nil, nil, errors.New("malformed index key")
	}

	// parts are fk/pk
	fk, pk = parts[0], parts[1]

	return
}

// FindCaptureFn is the mechanism for closing over the key and decoded value pair
// for adding results to the call sites collection. This generic implementation allows
// it to be reused. The returned decodedVal should always satisfy whatever decoding
// of the bucket value was set on the store that calls Find.
type FindCaptureFn func(key, value []byte) error

type FindOpts struct {
	manta.FindOptions
	CaptureFn FindCaptureFn
	FilterFn  FilterFn
}

// IndexMapping is a type which configures and Index to map items
// from a source bucket to an index bucket via a mapping known as
// IndexSourceOn. This function is called on the values in the source
// to derive the foreign key on which to index each item.
type IndexMapping interface {
	SourceBucket() []byte
	IndexBucket() []byte
	IndexSourceOn(value []byte) (foreignKey []byte, err error)
}

// NewIndexMapping creates an implementation of IndexMapping for the provided source bucket
// to a destination index bucket.
func NewIndexMapping(sourceBucket, indexBucket []byte, fn IndexSourceOnFunc) IndexMapping {
	return indexMapping{
		source: sourceBucket,
		index:  indexBucket,
		fn:     fn,
	}
}

// IndexSourceOnFunc is a function which can be used to derive the foreign key
// of a value in a source bucket.
type IndexSourceOnFunc func([]byte) ([]byte, error)

type indexMapping struct {
	source []byte
	index  []byte
	fn     IndexSourceOnFunc
}

func (i indexMapping) SourceBucket() []byte { return i.source }

func (i indexMapping) IndexBucket() []byte { return i.index }

func (i indexMapping) IndexSourceOn(v []byte) ([]byte, error) {
	return i.fn(v)
}

type Index struct {
	IndexMapping

	// canRead configures whether or not Walk accesses the index at all
	// or skips the index altogether and returns nothing.
	// This is used when you want to integrate only the write path before
	// releasing the read path.
	canRead bool
}

// IndexOption is a function which configures an index
type IndexOption func(*Index)

// WithIndexReadPathEnabled enables the read paths of the index (Walk)
// This should be enabled once the index has been fully populated and
// the Insert and Delete paths are correctly integrated.
func WithIndexReadPathEnabled(i *Index) {
	i.canRead = true
}

// NewIndex configures and returns a new *Index for a given index mapping.
// By default the read path (Walk) is disabled. This is because the index needs to
// be fully populated before depending upon the read path.
// The read path can be enabled using WithIndexReadPathEnabled option.
func NewIndex(mapping IndexMapping, opts ...IndexOption) *Index {
	index := &Index{IndexMapping: mapping}

	for _, opt := range opts {
		opt(index)
	}

	return index
}

func (i *Index) indexBucket(tx Tx) (Bucket, error) {
	return tx.Bucket(i.IndexBucket())
}

func (i *Index) sourceBucket(tx Tx) (Bucket, error) {
	return tx.Bucket(i.SourceBucket())
}

// Insert creates a single index entry for the provided primary key on the foreign key.
func (i *Index) Insert(tx Tx, foreignKey, primaryKey []byte) error {
	b, err := i.indexBucket(tx)
	if err != nil {
		return err
	}

	return b.Put(IndexKey(foreignKey, primaryKey), primaryKey)
}

// Delete removes the foreignKey and primaryKey mapping from the underlying index.
func (i *Index) Delete(tx Tx, foreignKey, primaryKey []byte) error {
	b, err := i.indexBucket(tx)
	if err != nil {
		return err
	}

	return b.Delete(IndexKey(foreignKey, primaryKey))
}

// Walk walks the source bucket using keys found in the index using the provided foreign key
// given the index has been fully populated.
func (i *Index) Walk(ctx context.Context, tx Tx, foreignKey []byte, visitFn VisitFunc) error {
	// skip walking if configured to do so as the index
	// is currently being used purely to write the index
	if !i.canRead {
		return nil
	}

	sourceBucket, err := i.sourceBucket(tx)
	if err != nil {
		return err
	}

	indexBucket, err := i.indexBucket(tx)
	if err != nil {
		return err
	}

	cursor, err := indexBucket.ForwardCursor(foreignKey,
		WithCursorPrefix(foreignKey))
	if err != nil {
		return err
	}

	return indexWalk(ctx, cursor, sourceBucket, visitFn)
}

// indexWalk consumes the IndexKey and primaryKey pairs in the index bucket and looks up their
// associated primaryKey's value in the provided source bucket.
// When an item is located in the source, the provided visit function is called with primary key and associated value.
func indexWalk(ctx context.Context, indexCursor ForwardCursor, sourceBucket Bucket, visit VisitFunc) (err error) {
	var keys [][]byte
	for ik, pk := indexCursor.Next(); ik != nil; ik, pk = indexCursor.Next() {
		keys = append(keys, pk)
	}

	if err := indexCursor.Err(); err != nil {
		return err
	}

	if err := indexCursor.Close(); err != nil {
		return err
	}

	values, err := sourceBucket.GetBatch(keys...)
	if err != nil {
		return err
	}

	for i, value := range values {
		if value != nil {
			if err := visit(keys[i], value); err != nil {
				return err
			}
		}
	}

	return nil
}
