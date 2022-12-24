package executor

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/f1shl3gs/manta/task/backend"
	"github.com/f1shl3gs/manta/task/backend/scheduler"
	"github.com/f1shl3gs/manta/task/mock"
)

type TaskHandler func(ctx context.Context, task *manta.Task, ts time.Time) error

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
	span, ctx := tracing.StartSpanFromContextWithOperationName(ctx, "Execute")
	defer span.Finish()

	span.LogKV("task_id", manta.ID(id).String())

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

	span.LogKV("run_id", run.ID.String())

	defer func() {
		if _, err = e.tcs.FinishRun(ctx, task.ID, run.ID); err != nil {
			e.logger.Error("Finish run failed",
				zap.String("task", task.ID.String()),
				zap.Error(err))
		}
	}()

	err = e.tcs.AddRunLog(ctx, task.ID, run.ID, time.Now(), "Start running")
	if err != nil {
		return err
	}

	err = e.tcs.UpdateRunState(ctx, task.ID, run.ID, time.Now(), manta.RunStarted)
	if err != nil {
		e.logger.Warn("Update started status failed",
			zap.String("task", task.ID.String()),
			zap.Error(err))
		return err
	}

	perr := e.checkFunc(ctx, task, scheduledFor)

	// success
	if perr == nil {
		now := time.Now()

		// add to run log
		err = e.tcs.AddRunLog(ctx, task.ID, run.ID, now, "Success")
		if err != nil {
			return err
		}

		// update run status
		err = e.tcs.UpdateRunState(ctx, task.ID, run.ID, now, manta.RunSuccess)
		if err != nil {
			return err
		}

		_, err = e.ts.UpdateTask(ctx, task.ID, manta.TaskUpdate{
			LatestCompleted: &now,
			LatestSuccess:   &now,
			LatestScheduled: &scheduledFor,
		})

		return err
	}

	// ctx is canceled
	if ctx.Err() != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = e.tcs.AddRunLog(ctx, task.ID, run.ID, time.Now(), "Context canceled")
		if err != nil {
			return err
		}

		// update run status
		err = e.tcs.UpdateRunState(ctx, task.ID, run.ID, time.Now(), manta.RunCanceled)
		if err != nil {
			e.logger.Warn("Update cancel status failed",
				zap.String("task", task.ID.String()),
				zap.Error(err))
		}

		return err
	}

	// exec error
	now := time.Now()
	errMsg := perr.Error()

	err = e.tcs.AddRunLog(ctx, task.ID, run.ID, now, fmt.Sprintf("Fail: %s", errMsg))
	if err != nil {
		return err
	}

	if err := e.tcs.UpdateRunState(ctx, task.ID, run.ID, now, manta.RunFail); err != nil {
		e.logger.Warn("Update fail status failed",
			zap.String("task", task.ID.String()),
			zap.Error(err))
	}

	_, err = e.ts.UpdateTask(ctx, task.ID, manta.TaskUpdate{
		LatestCompleted: &now,
		LatestFailure:   &now,
		LastRunError:    &errMsg,
	})

	return err
}
