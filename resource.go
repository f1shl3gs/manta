package manta

import "context"

type ResourceLog interface {
	Log(ctx context.Context, c Change) error
}

func (m *Change) Validate() error {
	// todo
	return nil
}
