package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	notificationEndpointBucket          = []byte("notificationendpoints")
	notificationEndpointOrgIndexBucket  = []byte("notificationendpointorgindex")
	notificationEndpointNameIndexBucket = []byte("notificationendpointnameindex")
)

func (s *Service) FindNotificationEndpointByID(ctx context.Context, id manta.ID) (*manta.NotificationEndpoint, error) {
	var (
		ne  *manta.NotificationEndpoint
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		ne, err = s.findNotificationEndpointByID(ctx, tx, id)
		return err
	})

	return ne, err
}

func (s *Service) findNotificationEndpointByID(ctx context.Context, tx Tx, id manta.ID) (*manta.NotificationEndpoint, error) {
	b, err := tx.Bucket(notificationEndpointBucket)
	if err != nil {
		return nil, err
	}

	key, err := id.Encode()
	if err != nil {
		return nil, err
	}

	data, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	ne := &manta.NotificationEndpoint{}
	if err = ne.Unmarshal(data); err != nil {
		return nil, err
	}

	return ne, nil
}

func (s *Service) FindNotificationEndpoints(ctx context.Context, filter manta.NotificationEndpointFilter, opts ...manta.FindOptions) ([]*manta.NotificationEndpoint, int, error) {
	var (
		edps  []*manta.NotificationEndpoint
		total int
		err   error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.OrgID != nil {
			edps, err = s.findNotificationEndpointsByOrgID(ctx, tx, *filter.OrgID)
			total = len(edps)
			return err
		}

		edps, total, err = s.findNotificationEndpoints(ctx, tx, filter, opts...)
		return err
	})

	if err != nil {
		return nil, 0, err
	}

	return edps, total, nil
}

func (s *Service) findNotificationEndpointsByOrgID(ctx context.Context, tx Tx, orgID manta.ID) ([]*manta.NotificationEndpoint, error) {
	prefix, err := orgID.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(notificationEndpointOrgIndexBucket)
	if err != nil {
		return nil, err
	}

	c, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix))
	if err != nil {
		return nil, err
	}

	defer c.Close()

	keys := make([][]byte, 0)
	for k, v := c.Next(); k != nil; k, v = c.Next() {
		keys = append(keys, v)
	}

	if err = c.Err(); err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return make([]*manta.NotificationEndpoint, 0), nil
	}

	b, err = tx.Bucket(notificationEndpointBucket)
	if err != nil {
		return nil, err
	}

	values, err := b.GetBatch(keys...)
	if err != nil {
		return nil, err
	}

	eps := make([]*manta.NotificationEndpoint, 0, len(values))
	for _, v := range values {
		if v == nil {
			continue
		}

		ep := &manta.NotificationEndpoint{}
		err = ep.Unmarshal(v)
		if err != nil {
			return nil, err
		}

		eps = append(eps, ep)
	}

	return eps, nil
}

func (s *Service) findNotificationEndpoints(ctx context.Context, tx Tx, filter manta.NotificationEndpointFilter, opts ...manta.FindOptions) ([]*manta.NotificationEndpoint, int, error) {
	if filter.ID != nil {
		ep, err := s.findNotificationEndpointByID(ctx, tx, *filter.ID)
		if err != nil {
			return nil, 0, err
		}

		return []*manta.NotificationEndpoint{ep}, 1, nil
	}

	if filter.Name != nil {
		ep, err := s.findNotificationEndpointByName(ctx, tx, *filter.Name)
		if err != nil {
			return nil, 0, err
		}

		return []*manta.NotificationEndpoint{ep}, 1, nil
	}

	// todo: support options
	edps := make([]*manta.NotificationEndpoint, 0)
	b, err := tx.Bucket(notificationEndpointBucket)
	if err != nil {
		return nil, 0, err
	}

	cur, err := b.Cursor()
	if err != nil {
		return nil, 0, err
	}

	iter := iterator{
		cursor: cur,
		decodeFn: func(key, val []byte) (k []byte, decodedVal interface{}, err error) {
			ne := &manta.NotificationEndpoint{}
			if err = ne.Unmarshal(val); err != nil {
				return nil, nil, err
			}

			return key, ne, nil
		},
	}

	for k, v, err := iter.Next(ctx); k != nil; k, v, err = iter.Next(ctx) {
		if err != nil {
			return nil, 0, err
		}

		edps = append(edps, v.(*manta.NotificationEndpoint))
	}

	return edps, len(edps), nil
}

func (s *Service) findNotificationEndpointByName(ctx context.Context, tx Tx, name string) (*manta.NotificationEndpoint, error) {
	b, err := tx.Bucket(notificationEndpointNameIndexBucket)
	if err != nil {
		return nil, err
	}

	val, err := b.Get([]byte(name))
	if err != nil {
		return nil, err
	}

	var id manta.ID
	if err = id.Decode(val); err != nil {
		return nil, err
	}

	return s.findNotificationEndpointByID(ctx, tx, id)
}

func (s *Service) CreateNotificationEndpoint(ctx context.Context, ne *manta.NotificationEndpoint) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.createNotificationEndpoint(ctx, tx, ne)
	})
}

func (s *Service) createNotificationEndpoint(ctx context.Context, tx Tx, ne *manta.NotificationEndpoint) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	ne.ID = s.idGen.ID()
	ne.Created = time.Now()
	ne.Updated = time.Now()

	return s.putNotificationEndpoint(ctx, tx, ne)
}

func (s *Service) putNotificationEndpoint(ctx context.Context, tx Tx, ne *manta.NotificationEndpoint) error {
	pk, err := ne.ID.Encode()
	if err != nil {
		return err
	}

	// name index
	indexKey := IndexKey([]byte(ne.Name), pk)
	b, err := tx.Bucket(notificationEndpointNameIndexBucket)
	if err != nil {
		return err
	}

	if err := b.Put(indexKey, pk); err != nil {
		return err
	}

	// org index
	fk, err := ne.OrgID.Encode()
	if err != nil {
		return err
	}

	indexKey = IndexKey(fk, pk)
	b, err = tx.Bucket(notificationEndpointOrgIndexBucket)
	if err != nil {
		return err
	}

	err = b.Put(indexKey, pk)
	if err != nil {
		return err
	}

	// save notification endpoint
	b, err = tx.Bucket(notificationEndpointBucket)
	if err != nil {
		return err
	}

	data, err := ne.Marshal()
	if err != nil {
		return err
	}

	return b.Put(pk, data)
}

func (s *Service) UpdateNotificationEndpoint(ctx context.Context, id manta.ID, u manta.NotificationEndpointUpdate) (*manta.NotificationEndpoint, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Context()

	var (
		after *manta.NotificationEndpoint
		err   error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		after, err = s.updateNotificationEndpoint(ctx, tx, id, u)
		return err
	})

	if err != nil {
		return nil, err
	}

	return after, err
}

func (s *Service) updateNotificationEndpoint(ctx context.Context, tx Tx, id manta.ID, u manta.NotificationEndpointUpdate) (*manta.NotificationEndpoint, error) {
	// check if it is exist
	ne, err := s.findNotificationEndpointByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	err = s.deleteNotificationEndpoint(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	// apply update
	u.Apply(ne)
	ne.Updated = time.Now()

	err = s.putNotificationEndpoint(ctx, tx, ne)
	if err != nil {
		return nil, err
	}

	return ne, nil
}

func (s *Service) DeleteNotificationEndpoint(ctx context.Context, id manta.ID) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteNotificationEndpoint(ctx, tx, id)
	})
}

func (s *Service) deleteNotificationEndpoint(ctx context.Context, tx Tx, id manta.ID) error {
	ne, err := s.findNotificationEndpointByID(ctx, tx, id)
	if err != nil {
		return err
	}

	idKey, _ := id.Encode()
	nameKey := []byte(ne.Name)

	// delete index
	nameIndexKey := IndexKey(nameKey, idKey)
	b, err := tx.Bucket(notificationEndpointNameIndexBucket)
	if err != nil {
		return err
	}

	if err := b.Delete(nameIndexKey); err != nil {
		return err
	}

	// delete notification endpoint
	b, err = tx.Bucket(notificationEndpointBucket)
	if err != nil {
		return err
	}

	return b.Delete(idKey)
}
