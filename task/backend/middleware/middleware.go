package middleware

import (
	"context"

	"github.com/f1shl3gs/manta"
)

// Coordinator is a type which is used to react to
// task related actions
type Coordinator interface {
	TaskCreated(context.Context, *manta.Task) error
	TaskUpdated(ctx context.Context, from, to *manta.Task) error
	TaskDeleted(context.Context, manta.ID) error

	// RunCancelled(ctx context.Context, runID manta.ID) error
	// RunRetried(ctx context.Context, task *proto.Task, run *influxdb.Run) error
	// RunForced(ctx context.Context, task *proto.Task, run *influxdb.Run) error
}
