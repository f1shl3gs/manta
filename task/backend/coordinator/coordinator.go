package coordinator

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/f1shl3gs/manta/task/backend/scheduler"
	"go.uber.org/zap"
)

// Executor is an abstraction of the task executor with only the functions needed by the coordinator
type Executor interface {
	ManualRun(ctx context.Context, id manta.ID, runID manta.ID) error
	Cancel(ctx context.Context, runID manta.ID) error
}

// SchedulableTask is a wrapper around the Task struct, giving it methods to make it compatible with the Scheduler
type SchedulableTask struct {
	*manta.Task

	sch scheduler.Schedule

	// last scheduled
	lsc time.Time
}

func (s SchedulableTask) ID() scheduler.ID {
	return scheduler.ID(s.Task.ID)
}

func (s SchedulableTask) Schedule() scheduler.Schedule {
	return s.sch
}

func (s SchedulableTask) Offset() time.Duration {
	// todo
	return 0
}

func (s SchedulableTask) LastScheduled() time.Time {
	return s.lsc
}

type Coordinator struct {
	logger *zap.Logger

	scheduler scheduler.Scheduler
	executor  Executor
}

func NewCoordinator(logger *zap.Logger, sch scheduler.Scheduler, executor Executor) *Coordinator {
	c := &Coordinator{
		logger:    logger,
		scheduler: sch,
		executor:  executor,
	}

	return c
}

func (c *Coordinator) TaskCreated(ctx context.Context, task *manta.Task) error {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	t, err := NewSchedulableTask(task)
	if err != nil {
		return err
	}

	return c.scheduler.Schedule(t)
}

func (c *Coordinator) TaskUpdated(ctx context.Context, from, to *manta.Task) error {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	sid := scheduler.ID(to.ID)
	t, err := NewSchedulableTask(to)
	if err != nil {
		return err
	}

	// if disabling the task, release it before schedule update
	if to.Status != from.Status && to.Status == string(manta.TaskInactive) {
		if err := c.scheduler.Release(sid); err != nil && err != manta.ErrTaskNotClaimed {
			return err
		}
	} else {
		if err := c.scheduler.Schedule(t); err != nil {
			return err
		}
	}

	return nil
}

func (c *Coordinator) TaskDeleted(ctx context.Context, id manta.ID) error {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	tid := scheduler.ID(id)
	if err := c.scheduler.Release(tid); err != nil && err != manta.ErrTaskNotClaimed {
		return err
	}

	return nil
}

func NewSchedulableTask(task *manta.Task) (SchedulableTask, error) {
	// todo: handle the last scheduled
	// for now always set it for now
	ts := time.Now()
	sch, ts, err := scheduler.NewSchedule(task.Cron, ts)
	if err != nil {
		return SchedulableTask{}, err
	}

	return SchedulableTask{
		Task: task,
		sch:  sch,
		lsc:  ts,
	}, nil
}
