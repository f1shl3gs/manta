package checks

import (
	"context"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

type Service struct {
	logger       *zap.Logger
	checkService manta.CheckService
	taskService  manta.TaskService
}

func (s *Service) FindCheckByID(ctx context.Context, id manta.ID) (*manta.Check, error) {
	return s.checkService.FindCheckByID(ctx, id)
}

func (s *Service) FindChecks(ctx context.Context, filter manta.CheckFilter, opts ...manta.FindOptions) ([]*manta.Check, int, error) {
	return s.checkService.FindChecks(ctx, filter, opts...)
}

func (s *Service) CreateCheck(ctx context.Context, c *manta.Check) error {
	panic("implement me")
}

func (s *Service) UpdateCheck(ctx context.Context, id manta.ID, c *manta.Check) (*manta.Check, error) {
	panic("implement me")
}

func (s *Service) PatchCheck(ctx context.Context, id manta.ID, u manta.CheckUpdate) (*manta.Check, error) {
	panic("implement me")
}

func (s *Service) DeleteCheck(ctx context.Context, id manta.ID) error {
	panic("implement me")
}
