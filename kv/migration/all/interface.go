package all

import (
	"context"

	"github.com/f1shl3gs/manta/kv"
)

// Spec is a specification for a particular migration.
type Spec interface {
	Name() string

	Up(ctx context.Context, store kv.SchemaStore) error
	Down(ctx context.Context, store kv.SchemaStore) error
}

type spec struct {
	name string
	up   func(ctx context.Context, store kv.SchemaStore) error
	down func(ctx context.Context, store kv.SchemaStore) error
}

func (s *spec) Name() string {
	return s.name
}

func (s *spec) Up(ctx context.Context, store kv.SchemaStore) error {
	return s.up(ctx, store)
}

func (s *spec) Down(ctx context.Context, store kv.SchemaStore) error {
	return s.down(ctx, store)
}
