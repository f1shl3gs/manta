package manta

import (
	"context"
	"time"
)

var (
	ErrTaskNotClaimed = &Error{
		Code: EConflict,
		Msg:  "task not claimed",
	}

	ErrOutOfBoundsLimit = &Error{
		Code: EUnprocessableEntity,
		Msg:  "run limit is out of bounds, must be between",
	}

	ErrRunNotFound = &Error{
		Code: ENotFound,
		Msg:  "run not found",
	}
)

type TaskType string

type TaskStatus string

const (
	TaskActive   TaskStatus = "active"
	TaskInactive TaskStatus = "inactive"
)

type Task struct {
	ID      ID        `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Type   TaskType   `json:"type,omitempty"`
	Status TaskStatus `json:"status,omitempty"`
	OrgID  ID         `json:"orgID,omitempty"`
	// OwnerID store the creater's id, e.g. check
	OwnerID ID     `json:"ownerID,omitempty"`
	Cron    string `json:"cron,omitempty"`

	// status
	LatestCompleted time.Time `json:"latestCompleted"`
	LatestScheduled time.Time `json:"latestScheduled"`
	LatestSuccess   time.Time `json:"latestSuccess"`
	LatestFailure   time.Time `json:"latestFailure"`
	LastRunStatus   string    `json:"lastRunStatus,omitempty"`
	LastRunError    string    `json:"lastRunError,omitempty"`
}

func (t *Task) GetID() ID {
	return t.ID
}

func (t *Task) GetOrgID() ID {
	return t.OrgID
}

type TaskFilter struct {
	OrgID   *ID
	OwnerID *ID
}

type TaskUpdate struct {
	Desc   *string
	Cron   *string
	Status *TaskStatus

	LatestCompleted *time.Time
	LatestScheduled *time.Time
	LatestSuccess   *time.Time
	LatestFailure   *time.Time
	LastRunError    *string
}

func (upd TaskUpdate) Apply(task *Task) {
	if upd.Status != nil {
		task.Status = *upd.Status
	}

	if upd.LatestCompleted != nil {
		task.LatestCompleted = *upd.LatestCompleted
	}

	if upd.LatestScheduled != nil {
		task.LatestScheduled = *upd.LatestScheduled
	}

	if upd.LatestSuccess != nil {
		task.LatestSuccess = *upd.LatestSuccess
	}

	if upd.LatestFailure != nil {
		task.LatestFailure = *upd.LatestFailure
	}

	if upd.LastRunError != nil {
		task.LastRunError = *upd.LastRunError
	}
}

// RunFilter represents a set of filters that restrict the returned results
type RunFilter struct {
	// Task ID is required for listing runs
	Task ID `json:"task"`

	Limit  int `json:"limit"`
	After  *time.Time
	Before *time.Time
}

type TaskService interface {
	// FindTaskByID returns a single task by id
	FindTaskByID(ctx context.Context, id ID) (*Task, error)

	// FindTasks returns all tasks which match the filter
	FindTasks(ctx context.Context, filter TaskFilter) ([]*Task, error)

	// FindRuns returns a list of runs that match a filter and the total count of returned runs.
	FindRuns(ctx context.Context, filter RunFilter) ([]*Run, int, error)

	// CreateTask creates a task
	CreateTask(ctx context.Context, task *Task) error

	// UpdateTask updates a single task with a patch
	UpdateTask(ctx context.Context, id ID, upd TaskUpdate) (*Task, error)

	// DeleteTask delete a single task by ID
	DeleteTask(ctx context.Context, id ID) error
}

type RunStatus string

const (
	RunStarted   RunStatus = "started"
	RunSuccess   RunStatus = "success"
	RunFail      RunStatus = "fail"
	RunCanceled  RunStatus = "canceled"
	RunScheduled RunStatus = "scheduled"
)

type RunLog struct {
	RunID   ID     `json:"runID"`
	Time    string `json:"time"`
	Message string `json:"message"`
}

type Run struct {
	ID           ID        `json:"id"`
	TaskID       ID        `json:"taskID"`
	ScheduledFor time.Time `json:"scheduledFor"`
	RunAt        time.Time `json:"runAt"`
	StartedAt    time.Time `json:"startedAt"`
	FinishedAt   time.Time `json:"finishedAt"`
	Status       RunStatus `json:"status"`
	Logs         []RunLog  `json:"logs,omitempty"`
}
