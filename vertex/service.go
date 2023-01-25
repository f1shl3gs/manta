package vertex

import (
	"context"
	"sync"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/broadcast"

	"go.uber.org/zap"
)

type CoordinatingVertexService struct {
	manta.ConfigService

	logger *zap.Logger

	mtx          sync.RWMutex
	broadcasters map[manta.ID]*broadcast.Broadcaster[*manta.Config]
}

func NewCoordinatingVertexService(
	configService manta.ConfigService,
	logger *zap.Logger,
) *CoordinatingVertexService {
	cs := &CoordinatingVertexService{
		ConfigService: configService,
		logger:        logger,
		broadcasters:  make(map[manta.ID]*broadcast.Broadcaster[*manta.Config]),
	}

	return cs
}

func (s *CoordinatingVertexService) Sub(id manta.ID) *broadcast.Queue[*manta.Config] {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	b, exist := s.broadcasters[id]
	if !exist {
		b = broadcast.New[*manta.Config]()
		s.broadcasters[id] = b
	}

	return b.Sub()
}

func (s *CoordinatingVertexService) CreateConfign(ctx context.Context, cf *manta.Config) error {
	return s.ConfigService.CreateConfig(ctx, cf)
}

func (s *CoordinatingVertexService) FindConfigByID(ctx context.Context, id manta.ID) (*manta.Config, error) {
	return s.ConfigService.FindConfigByID(ctx, id)
}

func (s *CoordinatingVertexService) FindConfigs(
	ctx context.Context,
	filter manta.ConfigFilter,
) ([]*manta.Config, error) {
	return s.ConfigService.FindConfigs(ctx, filter)
}

func (s *CoordinatingVertexService) UpdateConfig(
	ctx context.Context,
	id manta.ID,
	upd manta.ConfigUpdate,
) (*manta.Config, error) {
	cf, err := s.ConfigService.UpdateConfig(ctx, id, upd)
	if err != nil {
		return nil, err
	}

	s.mtx.RLock()
	defer s.mtx.RUnlock()

	b := s.broadcasters[id]
	if b != nil {
		b.Pub(cf)
	}

	return cf, nil
}

func (s *CoordinatingVertexService) DeleteConfig(ctx context.Context, id manta.ID) error {
	err := s.ConfigService.DeleteConfig(ctx, id)
	if err != nil {
		return err
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	b := s.broadcasters[id]
	if b == nil {
		return nil
	}

	delete(s.broadcasters, id)
	b.Close()

	return nil
}
