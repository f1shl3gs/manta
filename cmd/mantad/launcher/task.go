package launcher

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
)

type tcs struct {
}

func (t *tcs) CreateRun(ctx context.Context, taskID manta.ID, scheduledFor time.Time, runAt time.Time) (*manta.Run, error) {
	panic("implement me")
}

func (t *tcs) CurrentlyRunning(ctx context.Context, taskID manta.ID) ([]*manta.Run, error) {
	panic("implement me")
}

func (t *tcs) ManualRuns(ctx context.Context, taskID manta.ID) ([]*manta.Run, error) {
	panic("implement me")
}

func (t *tcs) StartManualRun(ctx context.Context, taskID, runID manta.ID) (*manta.Run, error) {
	panic("implement me")
}

func (t *tcs) FinishRun(ctx context.Context, taskID, runID manta.ID) (*manta.Run, error) {
	panic("implement me")
}

func (t *tcs) UpdateRunState(ctx context.Context, taskID, runID manta.ID, when time.Time, state manta.RunStatus) error {
	panic("implement me")
}

func (t *tcs) AddRunLog(ctx context.Context, taskID, runID manta.ID, when time.Time, log string) error {
	panic("implement me")
}
