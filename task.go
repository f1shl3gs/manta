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
	OrgID *ID
}

type TaskUpdate struct {
	Status *string

	LatestCompleted *time.Time
	LatestScheduled *time.Time
	LatestSuccess   *time.Time
	LatestFailure   *time.Time
	LastRunError    *string
}

type TaskService interface {
	// FindTaskByID returns a single task by id
	FindTaskByID(ctx context.Context, id ID) (*Task, error)

	// FindTasks returns
	FindTasks(ctx context.Context, filter TaskFilter) ([]*Task, error)

	// CreateTask
	CreateTask(ctx context.Context, task *Task) error

	// UpdateTask
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
