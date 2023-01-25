package kv

import (
	"context"
	"encoding/binary"
	"encoding/json"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	// ChangesBucket
	//   key:    ResourceID + Timestamp
	//   value:  Marshaled OperationLogEntry
	ChangesBucket = []byte("changes")

	// ChangesUserIndexBucket
	//   key:    UserID + Timestamp
	//   value:  ResourceID + Timestamp (key format of Changes Bucket)
	ChangesUserIndexBucket = []byte("changesuserindex")
)

func changeKey(c *manta.OperationLogEntry) ([]byte, error) {
	buf, err := c.ResourceID.Encode()
	if err != nil {
		return nil, err
	}

	buf = binary.BigEndian.AppendUint64(buf, uint64(c.Time.UnixNano()))
	return buf, nil
}

func changeUserIndexKey(c *manta.OperationLogEntry) ([]byte, error) {
	buf, err := c.UserID.Encode()
	if err != nil {
		return nil, err
	}

	buf = binary.BigEndian.AppendUint64(buf, uint64(c.Time.UnixNano()))
	return buf, nil
}

// AddLogEntry add an operation log entry.
func (s *Service) AddLogEntry(ctx context.Context, ent manta.OperationLogEntry) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	if err := ent.Valid(); err != nil {
		return err
	}

	key, err := changeKey(&ent)
	if err != nil {
		return err
	}

	indexKey, err := changeUserIndexKey(&ent)
	if err != nil {
		return err
	}

	value, err := json.Marshal(ent)
	if err != nil {
		return err
	}

	return s.kv.Update(ctx, func(tx Tx) error {
		b, err := tx.Bucket(ChangesBucket)
		if err != nil {
			return err
		}

		err = b.Put(key, value)
		if err != nil {
			return err
		}

		b, err = tx.Bucket(ChangesUserIndexBucket)
		if err != nil {
			return err
		}

		return b.Put(indexKey, key)
	})
}

// FindOperationLogsByID return operation logs of a resource.
func (s *Service) FindOperationLogsByID(ctx context.Context, id manta.ID, opts manta.FindOptions) ([]*manta.OperationLogEntry, int, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	prefix, err := id.Encode()
	if err != nil {
		return nil, 0, err
	}

	changes := []*manta.OperationLogEntry{}
	err = s.kv.View(ctx, func(tx Tx) error {
		b, err := tx.Bucket(ChangesBucket)
		if err != nil {
			return err
		}

		cursor, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix), WithCursorLimit(opts.Limit))
		if err != nil {
			return err
		}

		for k, v := cursor.Next(); k != nil; k, v = cursor.Next() {
			c := &manta.OperationLogEntry{}
			err = json.Unmarshal(v, c)
			if err != nil {
				return err
			}

			changes = append(changes, c)
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return changes, len(changes), nil
}

// FindOperationLogsByUser returns operation logs made by a user.
func (s *Service) FindOperationLogsByUser(ctx context.Context, userID manta.ID, opts manta.FindOptions) ([]*manta.OperationLogEntry, int, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	prefix, err := userID.Encode()
	if err != nil {
		return nil, 0, err
	}

	var changes []*manta.OperationLogEntry
	err = s.kv.View(ctx, func(tx Tx) error {
		b, err := tx.Bucket(ChangesUserIndexBucket)
		if err != nil {
			return err
		}

		cursor, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix), WithCursorLimit(opts.Limit))
		if err != nil {
			return err
		}

		b, err = tx.Bucket(ChangesBucket)
		if err != nil {
			return err
		}

		for k, v := cursor.Next(); k != nil; k, v = cursor.Next() {
			v, err = b.Get(v)
			if err != nil {
				return err
			}

			c := &manta.OperationLogEntry{}
			err = json.Unmarshal(v, c)
			if err != nil {
				return err
			}

			changes = append(changes, c)
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return changes, len(changes), nil
}
