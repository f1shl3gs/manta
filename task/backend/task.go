package backend

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
)

// TaskControlService is a low-level controller interface, intended to be passed to
// task executors and schedulers, which allows creation, completion, and status updates of runs.
type TaskControlService interface {

	// CreateRun creates a run with a scheduled for time.
	CreateRun(ctx context.Context, taskID manta.ID, scheduledFor time.Time, runAt time.Time) (*manta.Run, error)

	CurrentlyRunning(ctx context.Context, taskID manta.ID) ([]*manta.Run, error)

	ManualRuns(ctx context.Context, taskID manta.ID) ([]*manta.Run, error)

	// StartManualRun pulls a manual run from the list and moves it to currently running.
	StartManualRun(ctx context.Context, taskID, runID manta.ID) (*manta.Run, error)

	// FinishRun removes runID from the list of running tasks and if its `ScheduledFor` is later then last completed update it.
	FinishRun(ctx context.Context, taskID, runID manta.ID) (*manta.Run, error)

	// UpdateRunState sets the run state at the respective time.
	UpdateRunState(ctx context.Context, taskID, runID manta.ID, when time.Time, state manta.RunStatus) error

	// AddRunLog adds a file line to the run.
	AddRunLog(ctx context.Context, taskID, runID manta.ID, when time.Time, log string) error
}
