package manta

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	TaskActive   = "active"
	TaskInactive = "inactive"
)

var (
	ErrTaskNotClaimed = errors.New("task not claimed")

	ErrRunNotFound = &Error{
		Code: ENotFound,
		Msg:  "run not found",
	}
)

type TaskFilter struct {
	OrgID   *ID
	OwnerID *ID
}

type TaskUpdate struct {
	Status *string

	LatestCompleted *time.Time
	LatestScheduled *time.Time
	LatestSuccess   *time.Time
	LatestFailure   *time.Time
	LastRunError    *string
}

func (udp TaskUpdate) Apply(task *Task) {
	if udp.Status != nil {
		task.Status = *udp.Status
	}
}

type TaskService interface {
	// FindTaskByID returns a single task by id
	FindTaskByID(ctx context.Context, id ID) (*Task, error)

	// FindTasks returns all tasks which match the filter
	FindTasks(ctx context.Context, filter TaskFilter) ([]*Task, error)

	// CreateTask creates a task
	CreateTask(ctx context.Context, task *Task) error

	// UpdateTask updates a single task with a patch
	UpdateTask(ctx context.Context, id ID, udp TaskUpdate) (*Task, error)

	// DeleteTask delete a single task by ID
	DeleteTask(ctx context.Context, id ID) error
}

type RunStatus int

const (
	RunStarted RunStatus = iota
	RunSuccess
	RunFail
	RunCanceled
	RunScheduled
)

func (r RunStatus) String() string {
	switch r {
	case RunStarted:
		return "started"
	case RunSuccess:
		return "success"
	case RunFail:
		return "failed"
	case RunCanceled:
		return "canceled"
	case RunScheduled:
		return "scheduled"
	}
	panic(fmt.Sprintf("unknown RunStatus: %d", r))
}
