package backend

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/task/backend/scheduler"
)

type UpdateTaskService interface {
	UpdateTask(ctx context.Context, id manta.ID, udp manta.TaskUpdate) (*manta.Task, error)
}

type SchedulableTaskService struct {
	UpdateTaskService
}

func NewSchedulableTaskService(ts UpdateTaskService) SchedulableTaskService {
	return SchedulableTaskService{ts}
}

func (s SchedulableTaskService) UpdateLastScheduled(ctx context.Context, id scheduler.ID, t time.Time) error {
	_, err := s.UpdateTask(ctx, manta.ID(id), manta.TaskUpdate{
		LatestScheduled: &t,
	})

	if err != nil {
		return errors.Wrap(err, "could not update last scheduled for task")
	}

	return nil
}
