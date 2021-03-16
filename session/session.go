package session

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/kv"
)

type Service struct {
	kv  kv.Store
	now func() time.Time
	ttl time.Duration
}

func New(store kv.Store, ttl time.Duration) *Service {
	return &Service{
		kv:  store,
		now: time.Now,
		ttl: ttl,
	}
}

func (s *Service) CreateSession(ctx context.Context, uid manta.ID) (*manta.Session, error) {
	panic("implement me")
}

func (s *Service) FindSession(ctx context.Context, id manta.ID) (*manta.Session, error) {
	panic("implement me")
}

func (s *Service) RevokeSession(ctx context.Context, id manta.ID) error {
	panic("implement me")
}

func (s *Service) RenewSession(ctx context.Context, id manta.ID, expiration time.Time) error {
	panic("implement me")
}
