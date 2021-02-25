package backend

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

// Coordinator is a type with a single method which
// is called when a task has been created
type Coordinator interface {
	TaskCreated(context.Context, *manta.Task) error
}

// NotifyCoordinatorOfExisting lists all tasks by the provided task service and for
// each task it calls the provided coordinators task created method
func NotifyCoordinatorOfExisting(ctx context.Context, log *zap.Logger, ts manta.TaskService, coord Coordinator) error {
	// If we missed a Create Action
	tasks, err := ts.FindTasks(ctx, manta.TaskFilter{})
	if err != nil {
		return err
	}

	latestCompleted := time.Now()
	for _, task := range tasks {
		if task.Status != manta.TaskActive {
			continue
		}

		task, err := ts.UpdateTask(context.Background(), task.ID, manta.TaskUpdate{
			LatestCompleted: &latestCompleted,
			LatestScheduled: &latestCompleted,
		})
		if err != nil {
			log.Error("Failed to set latestCompleted", zap.Error(err))
			continue
		}

		if err = coord.TaskCreated(ctx, task); err != nil {
			log.Error("task create failed",
				zap.Error(err))
		}
	}

	return nil
}
