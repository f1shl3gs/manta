package vertex

import (
	"context"
	"sync"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/broadcast"

	"go.uber.org/zap"
)

type CoordinatingVertexService struct {
	manta.ConfigurationService

	logger *zap.Logger

	mtx          sync.RWMutex
	broadcasters map[manta.ID]*broadcast.Broadcaster[*manta.Configuration]
}

func NewCoordinatingVertexService(
	configurationService manta.ConfigurationService,
	logger *zap.Logger,
) *CoordinatingVertexService {
	cs := &CoordinatingVertexService{
		ConfigurationService: configurationService,
		logger:               logger,
		broadcasters:         make(map[manta.ID]*broadcast.Broadcaster[*manta.Configuration]),
	}

	return cs
}

func (s *CoordinatingVertexService) Sub(id manta.ID) *broadcast.Queue[*manta.Configuration] {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	b, exist := s.broadcasters[id]
	if !exist {
		b = broadcast.New[*manta.Configuration]()
		s.broadcasters[id] = b
	}

	return b.Sub()
}

func (s *CoordinatingVertexService) CreateConfiguration(ctx context.Context, cf *manta.Configuration) error {
	return s.ConfigurationService.CreateConfiguration(ctx, cf)
}

func (s *CoordinatingVertexService) GetConfiguration(ctx context.Context, id manta.ID) (*manta.Configuration, error) {
	return s.ConfigurationService.GetConfiguration(ctx, id)
}

func (s *CoordinatingVertexService) FindConfigurations(
	ctx context.Context,
	filter manta.ConfigurationFilter,
) ([]*manta.Configuration, error) {
	return s.ConfigurationService.FindConfigurations(ctx, filter)
}

func (s *CoordinatingVertexService) UpdateConfiguration(
	ctx context.Context,
	id manta.ID,
	upd manta.ConfigurationUpdate,
) (*manta.Configuration, error) {
	cf, err := s.ConfigurationService.UpdateConfiguration(ctx, id, upd)
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

func (s *CoordinatingVertexService) DeleteConfiguration(ctx context.Context, id manta.ID) error {
	err := s.ConfigurationService.DeleteConfiguration(ctx, id)
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
