package kv

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	// UrmsBucket is the bucket to store UserResourceMapping
	// Key: resource_id + user_id
	// Value: manta.UserResourceMapping
	UrmsBucket = []byte("userresourcemappings")

	// UrmUserIndexBucket is the bucket to store user index
	// Key: user_id + resource_id
	// Value: nil
	UrmUserIndexBucket = []byte("userresourcemappinguserindex")
)

// FindUserResourceMappings returns a list of UserResourceMappings that match filter and the total count of matching mappings.
func (s *Service) FindUserResourceMappings(ctx context.Context, filter manta.UserResourceMappingFilter, opts ...manta.FindOptions) ([]*manta.UserResourceMapping, int, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		list  []*manta.UserResourceMapping
		total int
		err   error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.UserID.Valid() {
			list, total, err = findUserResourceMappingByUser(tx, filter, opts[0])
		} else if filter.ResourceID.Valid() {
			list, total, err = findUserResourceMappingByResource(tx, filter, opts[0])
		} else {
			// TODO
		}

		return err
	})

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func findUserResourceMappingByUser(tx Tx, filter manta.UserResourceMappingFilter, opts manta.FindOptions) ([]*manta.UserResourceMapping, int, error) {
	var (
		list []*manta.UserResourceMapping
		seen = 0
	)

	filterFn := func(m *manta.UserResourceMapping) bool {
		return (!filter.UserID.Valid() || (filter.UserID == m.UserID)) &&
			(!filter.ResourceID.Valid() || (filter.ResourceID == m.ResourceID)) &&
			(filter.UserType == "" || (filter.UserType == m.UserType)) &&
			(filter.ResourceType == "" || (filter.ResourceType == m.ResourceType))
	}

	var (
		keys      [][]byte
		prefix, _ = filter.UserID.Encode()
	)

	b, err := tx.Bucket(UrmUserIndexBucket)
	if err != nil {
		return nil, 0, err
	}

	c, err := b.Cursor()
	for k, _ := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = c.Next() {
		fk, pk, err := indexKeyParts(k)
		if err != nil {
			return nil, 0, err
		}

		keys = append(keys, IndexKey(pk, fk))
	}

	if len(keys) == 0 {
		return nil, 0, nil
	}

	b, err = tx.Bucket(UrmsBucket)
	if err != nil {
		return nil, 0, err
	}

	values, err := b.GetBatch(keys...)
	if err != nil {
		return nil, 0, err
	}

	for _, value := range values {
		if len(value) == 0 {
			continue
		}

		m := &manta.UserResourceMapping{}
		err = json.Unmarshal(value, m)
		if err != nil {
			return nil, 0, err
		}

		reachOffset := seen >= opts.Offset
		if reachOffset && filterFn(m) {
			list = append(list, m)
		}

		seen += 1
	}

	return list, len(list), nil
}

func findUserResourceMappingByResource(tx Tx, filter manta.UserResourceMappingFilter, opts manta.FindOptions) ([]*manta.UserResourceMapping, int, error) {
	var (
		list []*manta.UserResourceMapping
		seen = 0
	)

	filterFn := func(m *manta.UserResourceMapping) bool {
		return (!filter.UserID.Valid() || (filter.UserID == m.UserID)) &&
			(!filter.ResourceID.Valid() || (filter.ResourceID == m.ResourceID)) &&
			(filter.UserType == "" || (filter.UserType == m.UserType)) &&
			(filter.ResourceType == "" || (filter.ResourceType == m.ResourceType))
	}

	prefix, _ := filter.ResourceID.Encode()
	b, err := tx.Bucket(UrmsBucket)
	if err != nil {
		return nil, 0, err
	}

	c, err := b.Cursor()
	if err != nil {
		return nil, 0, err
	}

	for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
		m := &manta.UserResourceMapping{}
		if err = json.Unmarshal(v, m); err != nil {
			return nil, 0, err
		}

		reachOffset := seen >= opts.Offset
		if reachOffset && filterFn(m) {
			list = append(list, m)
		}

		seen += 1
	}

	return list, len(list), nil

}

// CreateUserResourceMapping creates a user resource mapping.
func (s *Service) CreateUserResourceMapping(ctx context.Context, m *manta.UserResourceMapping) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		data, err := json.Marshal(m)
		if err != nil {
			return err
		}

		b, err := tx.Bucket(UrmsBucket)
		if err != nil {
			return err
		}

		key, err := indexIDKey(m.ResourceID, m.UserID)
		if err != nil {
			return err
		}

		return b.Put(key, data)
	})
}

// DeleteUserResourceMapping deletes a user resource mapping.
func (s *Service) DeleteUserResourceMapping(ctx context.Context, resourceID, userID manta.ID) error {
	key, err := indexIDKey(resourceID, userID)
	if err != nil {
		return err
	}

	return s.kv.Update(ctx, func(tx Tx) error {
		b, err := tx.Bucket(UrmsBucket)
		if err != nil {
			return err
		}

		return b.Delete(key)
	})
}
