package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

type CoordinatingCheckService struct {
	manta.CheckService

	logger      *zap.Logger
	coordinator Coordinator
	taskService manta.TaskService
}

func NewCheckService(cs manta.CheckService, ts manta.TaskService, coord Coordinator) *CoordinatingCheckService {
	return &CoordinatingCheckService{
		CheckService: cs,
		coordinator:  coord,
		taskService:  ts,
	}
}

func (cs *CoordinatingCheckService) CreateCheck(ctx context.Context, check *manta.Check) error {
	if err := cs.CheckService.CreateCheck(ctx, check); err != nil {
		return err
	}

	task, err := cs.taskService.FindTaskByID(ctx, check.TaskID)
	if err != nil {
		return err
	}

	err = cs.coordinator.TaskCreated(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CoordinatingCheckService) UpdateCheck(ctx context.Context, id manta.ID, check *manta.Check) (*manta.Check, error) {
	from, err := cs.CheckService.FindCheckByID(ctx, id)
	if err != nil {
		return nil, err
	}

	fromTask, err := cs.taskService.FindTaskByID(ctx, from.TaskID)
	if err != nil {
		return nil, err
	}

	to, err := cs.CheckService.UpdateCheck(ctx, id, check)
	if err != nil {
		return nil, err
	}

	toTask, err := cs.taskService.FindTaskByID(ctx, to.TaskID)
	if err != nil {
		return nil, err
	}

	// if the update is to active and the previous task was inactive we should add a "latest completed"
	// update this allows us to see not run the task for inactive time
	if fromTask.Status == manta.TaskInactive && toTask.Status == manta.TaskActive {
		toTask.LatestCompleted = time.Now()
	}

	return check, cs.coordinator.TaskUpdated(ctx, fromTask, toTask)
}

func (cs *CoordinatingCheckService) PatchCheck(ctx context.Context, id manta.ID, upd manta.CheckUpdate) (*manta.Check, error) {
	from, err := cs.CheckService.FindCheckByID(ctx, id)
	if err != nil {
		return nil, err
	}

	fromTask, err := cs.taskService.FindTaskByID(ctx, from.TaskID)
	if err != nil {
		return nil, err
	}

	to, err := cs.CheckService.PatchCheck(ctx, id, upd)
	if err != nil {
		return nil, err
	}

	toTask, err := cs.taskService.FindTaskByID(ctx, to.TaskID)
	if err != nil {
		return nil, err
	}

	// if the update is to activate and the previous task was inactive we should add a "latest completed" update
	// this allows us to see not run the task for inactive time
	if fromTask.Status == manta.TaskInactive && toTask.Status == manta.TaskActive {
		toTask.LatestCompleted = time.Now()
	}

	return to, cs.coordinator.TaskUpdated(ctx, fromTask, toTask)
}

func (cs *CoordinatingCheckService) DeleteCheck(ctx context.Context, id manta.ID) error {
	tasks, err := cs.taskService.FindTasks(ctx, manta.TaskFilter{OwnerID: &id})
	if err != nil {
		return err
	}

	if err = cs.CheckService.DeleteCheck(ctx, id); err != nil {
		return err
	}

	for _, task := range tasks {
		err = cs.coordinator.TaskDeleted(ctx, task.ID)
		if err != nil {
			cs.logger.Error("Delete task from coordinator failed",
				zap.String("check", id.String()),
				zap.String("task", task.ID.String()),
				zap.Error(err))
		}
	}

	return nil
}
