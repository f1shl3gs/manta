package kv

import (
	"context"

	"github.com/f1shl3gs/manta"
)

var (
	bucketBucket         = []byte("buckets")
	bucketOrgIndexBucket = []byte("bucketorgindex")
)

func (s *Service) FindBucketByID(ctx context.Context, id manta.ID) (*manta.Bucket, error) {
	panic("implement me")
}

func (s *Service) findBucketByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Bucket, error) {
	pk, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(bucketBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	bkt := &manta.Bucket{}
	err = bkt.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	return bkt, nil
}

func (s *Service) FindBucket(ctx context.Context, filter manta.BucketFilter) (*manta.Bucket, error) {
	panic("implement me")
}

func (s *Service) FindBuckets(ctx context.Context, filter manta.BucketFilter) ([]*manta.Bucket, error) {
	panic("implement me")
}

func (s *Service) CreateBucket(ctx context.Context, b *manta.Bucket) error {
	panic("implement me")
}

func (s *Service) UpdateBucket(ctx context.Context, id manta.ID, upd manta.BucketUpdate) (*manta.Bucket, error) {
	panic("implement me")
}

func (s *Service) DeleteBucket(ctx context.Context, id manta.ID) {
	panic("implement me")
}
