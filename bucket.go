package manta

import (
	"context"
	"time"
)

type BucketFilter struct {
	OrgID *ID
}

type BucketUpdate struct {
	Retention *time.Duration
}

func (upd *BucketUpdate) Apply(b *Bucket) {
	if upd.Retention != nil {
		b.Retention = *upd.Retention
	}
}

// BucketService represents a service for managing bucket data
type BucketService interface {
	// FindBucketByID returns a single bucket by ID
	FindBucketByID(ctx context.Context, id ID) (*Bucket, error)

	// FindBucket returns the first bucket that matches filter
	FindBucket(ctx context.Context, filter BucketFilter) (*Bucket, error)

	FindBuckets(ctx context.Context, filter BucketFilter) ([]*Bucket, error)

	// CreateBucket creates a new bucket and sets b.ID with the new identifier
	CreateBucket(ctx context.Context, b *Bucket) error

	// UpdateBucket updates a single bucket with changeset
	UpdateBucket(ctx context.Context, id ID, upd BucketUpdate) (*Bucket, error)

	// DeleteBucket removes a bucket by ID
	DeleteBucket(ctx context.Context, id ID)
}
