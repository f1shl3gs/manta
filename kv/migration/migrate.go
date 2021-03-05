package migration

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/kv/migration/all"
)

type Migration struct {
	ID   manta.ID `json:"id"`
	Name string   `json:"name"`
	// State value should be 'up' or 'down'
	State      string     `json:"state"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

var (
	All = []all.Spec{
		all.Migration0000Initial(),
	}

	//
	migrationBucket = []byte("migrations")
)

type Migrator struct {
	logger *zap.Logger
	store  kv.SchemaStore
	specs  []all.Spec
}

func New(logger *zap.Logger, store kv.SchemaStore, specs ...all.Spec) *Migrator {
	return &Migrator{
		logger: logger,
		store:  store,
		specs:  specs,
	}
}

func (m *Migrator) Up(ctx context.Context) error {
	var (
		from int
	)

	err := m.store.CreateBucket(ctx, migrationBucket)
	if err != nil {
		return err
	}

	// find the last applied index
	for from = 0; from < len(m.specs); from++ {
		id := manta.ID(from + 1)
		mig, err := m.getMigration(ctx, id)
		if err != nil {
			if err == kv.ErrKeyNotFound {
				break
			}

			return err
		}

		if mig.FinishedAt == nil {
			break
		}
	}

	// apply the remain specs
	for i := from; i < len(m.specs); i++ {
		id := manta.ID(from + 1)
		mig, err := m.getMigration(ctx, id)
		if err != nil {
			if err == kv.ErrKeyNotFound {
				break
			}

			return err
		}

		if mig.FinishedAt == nil {
			break
		}

		spec := m.specs[i]
		started := time.Now()

		migration := &Migration{
			ID:        manta.ID(i + 1),
			Name:      spec.Name(),
			StartedAt: &started,
		}

		err = m.putMigration(ctx, migration)
		if err != nil {
			return err
		}

		err = spec.Up(ctx, m.store)
		if err != nil {
			return err
		}

		finished := time.Now()
		migration.FinishedAt = &finished
		err = m.putMigration(ctx, migration)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) getMigration(ctx context.Context, id manta.ID) (*Migration, error) {
	var (
		mig *Migration
		err error
	)

	err = m.store.View(ctx, func(tx kv.Tx) error {
		b, err := tx.Bucket(migrationBucket)
		if err != nil {
			return err
		}

		pk, err := id.Encode()
		if err != nil {
			return err
		}

		data, err := b.Get(pk)
		if err != nil {
			return err
		}

		mig = &Migration{}
		return json.Unmarshal(data, mig)
	})

	if err != nil {
		return nil, err
	}

	return mig, nil
}

func (m *Migrator) putMigration(ctx context.Context, mig *Migration) error {
	return m.store.Update(ctx, func(tx kv.Tx) error {
		b, err := tx.Bucket(migrationBucket)
		if err != nil {
			return err
		}

		pk, err := mig.ID.Encode()
		if err != nil {
			return err
		}

		data, err := json.Marshal(mig)
		if err != nil {
			return err
		}

		return b.Put(pk, data)
	})
}
