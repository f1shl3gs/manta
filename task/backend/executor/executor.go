package executor

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/f1shl3gs/manta/task/backend"
	"github.com/f1shl3gs/manta/task/backend/scheduler"
	"github.com/f1shl3gs/manta/task/mock"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type TaskHandler func(ctx context.Context, task *manta.Task) error

type Executor struct {
	logger *zap.Logger
	ts     manta.TaskService
	tcs    backend.TaskControlService

	checkFunc TaskHandler
}

func NewExecutor(logger *zap.Logger, ts manta.TaskService, tcs backend.TaskControlService, cf TaskHandler) *Executor {
	return &Executor{
		logger:    logger,
		ts:        ts,
		tcs:       tcs,
		checkFunc: cf,
	}
}

func (e *Executor) Execute(ctx context.Context, id scheduler.ID, scheduledFor time.Time, runAt time.Time) error {
	span, ctx := tracing.StartSpanFromContextWithOperationName(ctx, "execute")
	defer span.Finish()

	task, err := e.ts.FindTaskByID(ctx, manta.ID(id))
	if err != nil {
		return err
	}

	mtcs, ok := e.tcs.(*mock.TaskControlService)
	if ok {
		mtcs.SetTask(task)
	}

	run, err := e.tcs.CreateRun(ctx, task.ID, scheduledFor, runAt)
	if err != nil {
		return err
	}

	defer func() {
		if _, err = e.tcs.FinishRun(ctx, task.ID, run.ID); err != nil {
			e.logger.Error("finish run failed",
				zap.String("task", task.ID.String()),
				zap.Error(err))
		}
	}()

	err = e.checkFunc(ctx, task)
	if err != nil {
		now := time.Now()
		errMsg := err.Error()

		_, err = e.ts.UpdateTask(ctx, task.ID, manta.TaskUpdate{
			LatestCompleted: &now,
			LatestFailure:   &now,
			LastRunError:    &errMsg,
		})

		if err != nil {
			return err
		}

		return nil
	}

	// success
	now := time.Now()
	_, err = e.ts.UpdateTask(ctx, task.ID, manta.TaskUpdate{
		LatestCompleted: &now,
		LatestSuccess:   &now,
	})

	if err != nil {
		return errors.Wrap(err, "update task status failed")
	}

	return nil
}
