package kv

import (
	"context"
	"encoding/json"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/notification"
)

var (
	NotificationEndpointsBucket        = []byte("notificationendpoints")
	NotificationENdpointOrgIndexBucket = []byte("notificationendpointorgindex")
)

// FindNotificationByID returns a single notification endpoint by ID
func (s *Service) FindNotificationEndpointByID(ctx context.Context, id manta.ID) (manta.NotificationEndpoint, error) {
	var (
		ne  manta.NotificationEndpoint
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		ne, err = findNotificationEndpointByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return ne, nil
}

func findNotificationEndpointByID(tx Tx, id manta.ID) (manta.NotificationEndpoint, error) {
	key, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(NotificationEndpointsBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	return notification.UnmarshalJSON(data)
}

// FindNotificationEndpoints returns a list of notication endpoints that match filter.
func (s *Service) FindNotificationEndpoints(ctx context.Context, filter manta.NotificationEndpointFilter) ([]manta.NotificationEndpoint, error) {
	var (
		list []manta.NotificationEndpoint
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		prefix := filter.OrgID.String()

		b, err := tx.Bucket(NotificationENdpointOrgIndexBucket)
		if err != nil {
			return err
		}

		cursor, err := b.Cursor(WithCursorHintPrefix(prefix))
		if err != nil {
			return err
		}

		keys := [][]byte{}
		for k, v := cursor.Next(); k != nil; k, v = cursor.Next() {
			keys = append(keys, v)
		}

		b, err = tx.Bucket(NotificationEndpointsBucket)
		if err != nil {
			return err
		}

		values, err := b.GetBatch(keys...)
		if err != nil {
			return err
		}

		for _, value := range values {
			ne, err := notification.UnmarshalJSON(value)
			if err != nil {
				return err
			}

			list = append(list, ne)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

// CreateNotificationEndpoint creates a new notification endpoint and sets b.ID with the new identifier
func (s *Service) CreateNotificationEndpoint(ctx context.Context, ne manta.NotificationEndpoint) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.createNotificationEndpoint(ctx, tx, ne)
	})
}

func (s *Service) createNotificationEndpoint(ctx context.Context, tx Tx, ne manta.NotificationEndpoint) error {
	ne.SetID(s.idGen.ID())
	now := time.Now()
	ne.SetCreated(now)
	ne.SetUpdated(now)
	ne.BackfillSecretKeys()

	if err := ne.Valid(); err != nil {
		return err
	}

	return putNotificationEndpoint(tx, ne)
}

func putNotificationEndpoint(tx Tx, ne manta.NotificationEndpoint) error {
	pk, err := ne.GetID().Encode()
	if err != nil {
		return err
	}

	fk, err := ne.GetOrgID().Encode()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(NotificationEndpointsBucket)
	if err != nil {
		return err
	}

	value, err := json.Marshal(ne)
	if err != nil {
		return err
	}

	err = b.Put(pk, value)
	if err != nil {
		return err
	}

	// index
	index := IndexKey(fk, pk)
	b, err = tx.Bucket(NotificationENdpointOrgIndexBucket)
	if err != nil {
		return err
	}

	return b.Put(index, pk)
}

// UpdateNotificationEndpoint updates a single notification endpoint.
// Returns the new notification endpoint after update.
func (s *Service) UpdateNotificationEndpoint(ctx context.Context, id manta.ID, ne manta.NotificationEndpoint) (manta.NotificationEndpoint, error) {
	var (
		updated manta.NotificationEndpoint
		err     error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		current, err := findNotificationEndpointByID(tx, id)
		if err != nil {
			return err
		}

		ne.SetCreated(current.GetCreated())
		ne.SetUpdated(time.Now())

		if err := ne.Valid(); err != nil {
			return err
		}

		updated = ne

		return putNotificationEndpoint(tx, ne)
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// PatchNotificationENdpoint patch a single notification endpoint.
// Returns the new notification endpoint after patch.
func (s *Service) PatchNotificationEndpoint(ctx context.Context, id manta.ID, upd manta.NotificationEndpointUpdate) (manta.NotificationEndpoint, error) {
	var (
		patched manta.NotificationEndpoint
		err     error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		ne, err := findNotificationEndpointByID(tx, id)
		if err != nil {
			return err
		}

		upd.Apply(ne)
		ne.SetUpdated(time.Now())

		if err = ne.Valid(); err != nil {
			return err
		}

		patched = ne

		return putNotificationEndpoint(tx, ne)
	})
	if err != nil {
		return nil, err
	}

	return patched, nil
}

// DeleteNotificationEndpoint remove a notification endpoint by ID, return it's secret fields, orgID for further deletion
func (s *Service) DeleteNotificationEndpoint(ctx context.Context, id manta.ID) ([]manta.SecretField, manta.ID, error) {
	var (
		orgID   manta.ID
		secrets []manta.SecretField
	)

	err := s.kv.Update(ctx, func(tx Tx) error {
		ne, err := findNotificationEndpointByID(tx, id)
		if err != nil {
			return err
		}

		orgID = ne.GetOrgID()
		secrets = ne.SecretFields()

		pk, err := ne.GetID().Encode()
		if err != nil {
			return err
		}

		fk, err := ne.GetOrgID().Encode()
		if err != nil {
			return err
		}

		// delete index
		index := IndexKey(fk, pk)
		b, err := tx.Bucket(NotificationENdpointOrgIndexBucket)
		if err != nil {
			return err
		}

		if err = b.Delete(index); err != nil {
			return err
		}

		// delete entity
		b, err = tx.Bucket(NotificationEndpointsBucket)
		if err != nil {
			return err
		}

		return b.Delete(pk)
	})

	if err != nil {
		return nil, 0, err
	}

	return secrets, orgID, nil
}
