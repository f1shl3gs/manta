package manta

import (
	"context"
	"encoding/json"
	"errors"
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

type Task struct {
	ID      ID        `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	// store some data for this task
	Annotations map[string]string `json:"annotations,omitempty"`
	Type        string            `json:"type,omitempty"`
	Status      string            `json:"status,omitempty"`
	OrgID       ID                `json:"orgID,omitempty"`
	// OwnerID store the creater's id
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

func (t *Task) Marshal() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Task) Unmarshal(data []byte) error {
	return json.Unmarshal(data, t)
}

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

type RunStatus string

const (
	RunStarted   RunStatus = "started"
	RunSuccess             = "success"
	RunFail                = "fail"
	RunCanceled            = "canceled"
	RunScheduled           = "scheduled"
)
