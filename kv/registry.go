package kv

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/f1shl3gs/manta"

	"go.uber.org/zap"
)

var (
	RegistryBucket = []byte("registry")
)

func (s *Service) Register(ctx context.Context, ins *manta.Instance) error {
	now := time.Now()
	ins.Created = now

	data, err := json.Marshal(ins)
	if err != nil {
		return err
	}

	return s.kv.Update(ctx, func(tx Tx) error {
		b, err := tx.Bucket(RegistryBucket)
		if err != nil {
			return err
		}

		return b.Put([]byte(ins.Uuid), data)
	})
}

var (
	once sync.Once
)

func (s *Service) Catalog(ctx context.Context) ([]*manta.Instance, error) {
	// TODO: remove this
	once.Do(func() {
		_ = s.kv.Update(ctx, func(tx Tx) error {
			b, err := tx.Bucket(RegistryBucket)
			if err != nil {
				return err
			}

			cursor, err := b.Cursor()
			if err != nil {
				return err
			}

			for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
				ins := &manta.Instance{}
				if err := json.Unmarshal(v, ins); err != nil {
					if err = b.Delete(k); err != nil {
						return err
					}

					continue
				}

			}

			return nil
		})
	})

	var (
		list []*manta.Instance
	)

	err := s.kv.View(ctx, func(tx Tx) error {
		b, err := tx.Bucket(RegistryBucket)
		if err != nil {
			return err
		}

		cursor, err := b.Cursor()
		if err != nil {
			return err
		}

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			ins := &manta.Instance{}
			if err := json.Unmarshal(v, ins); err != nil {
				s.logger.Warn("unmarshal instance failed", zap.Error(err), zap.ByteString("value", v))
				continue
			}

			list = append(list, ins)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}
