package mock

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/snowflake"
	"github.com/f1shl3gs/manta/task/backend"
)

var idgen = snowflake.NewIDGenerator()

// TaskControlService is a mock implementation of TaskControlService (used by NewScheduler).
type TaskControlService struct {
	mu sync.Mutex
	// Map of stringified task ID to last ID used for run.
	runs map[manta.ID]map[manta.ID]*manta.Run

	// Map of stringified, concatenated task and platform ID, to runs that have been created.
	created map[string]*manta.Run

	// Map of stringified task ID to task meta.
	tasks      map[manta.ID]*manta.Task
	manualRuns []*manta.Run
	// Map of task ID to total number of runs created for that task.
	totalRunsCreated map[manta.ID]int
	finishedRuns     map[manta.ID]*manta.Run
}

var _ backend.TaskControlService = (*TaskControlService)(nil)

func NewTaskControlService() *TaskControlService {
	return &TaskControlService{
		runs:             make(map[manta.ID]map[manta.ID]*manta.Run),
		finishedRuns:     make(map[manta.ID]*manta.Run),
		tasks:            make(map[manta.ID]*manta.Task),
		created:          make(map[string]*manta.Run),
		totalRunsCreated: make(map[manta.ID]int),
	}
}

// SetTask sets the task.
// SetTask must be called before CreateNextRun, for a given task ID.
func (tcs *TaskControlService) SetTask(task *manta.Task) {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()

	tcs.tasks[task.ID] = task
}

func (tcs *TaskControlService) SetManualRuns(runs []*manta.Run) {
	tcs.manualRuns = runs
}

func (tcs *TaskControlService) CreateRun(_ context.Context, taskID manta.ID, scheduledFor time.Time, runAt time.Time) (*manta.Run, error) {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()

	runID := idgen.ID()
	runs, ok := tcs.runs[taskID]
	if !ok {
		runs = make(map[manta.ID]*manta.Run)
	}
	runs[runID] = &manta.Run{
		ID:           runID,
		ScheduledFor: scheduledFor,
	}
	tcs.runs[taskID] = runs
	return runs[runID], nil
}

func (tcs *TaskControlService) StartManualRun(_ context.Context, taskID, runID manta.ID) (*manta.Run, error) {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()

	var run *manta.Run
	for i, r := range tcs.manualRuns {
		if r.ID == runID {
			run = r
			tcs.manualRuns = append(tcs.manualRuns[:i], tcs.manualRuns[i+1:]...)
		}
	}
	if run == nil {
		return nil, manta.ErrRunNotFound
	}
	return run, nil
}

func (tcs *TaskControlService) FinishRun(_ context.Context, taskID, runID manta.ID) (*manta.Run, error) {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()

	tid := taskID
	rid := runID
	r := tcs.runs[tid][rid]
	delete(tcs.runs[tid], rid)
	t := tcs.tasks[tid]

	if r.ScheduledFor.After(t.LatestCompleted) {
		t.LatestCompleted = r.ScheduledFor
	}

	tcs.finishedRuns[rid] = r
	delete(tcs.created, tid.String()+rid.String())
	return r, nil
}

func (tcs *TaskControlService) CurrentlyRunning(ctx context.Context, taskID manta.ID) ([]*manta.Run, error) {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()
	rtn := []*manta.Run{}
	for _, run := range tcs.runs[taskID] {
		rtn = append(rtn, run)
	}
	return rtn, nil
}

func (tcs *TaskControlService) ManualRuns(ctx context.Context, taskID manta.ID) ([]*manta.Run, error) {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()

	if tcs.manualRuns != nil {
		return tcs.manualRuns, nil
	}
	return []*manta.Run{}, nil
}

// UpdateRunState sets the run state at the respective time.
func (tcs *TaskControlService) UpdateRunState(ctx context.Context, taskID, runID manta.ID, when time.Time, state manta.RunStatus) error {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()

	run, ok := tcs.runs[taskID][runID]
	if !ok {
		panic("run state called without a run")
	}
	switch state {
	case manta.RunStarted:
		run.StartedAt = when
	case manta.RunSuccess, manta.RunFail, manta.RunCanceled:
		run.FinishedAt = when
	case manta.RunScheduled:
		// nothing
	default:
		panic("invalid status")
	}
	run.Status = state.String()
	return nil
}

// AddRunLog adds a file line to the run.
func (tcs *TaskControlService) AddRunLog(ctx context.Context, taskID, runID manta.ID, when time.Time, log string) error {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()

	run := tcs.runs[taskID][runID]
	if run == nil {
		panic("cannot add a file to a non existent run")
	}
	run.Logs = append(run.Logs, manta.RunLog{RunID: runID, Time: when.Format(time.RFC3339Nano), Message: log})
	return nil
}

func (tcs *TaskControlService) CreatedFor(taskID manta.ID) []*manta.Run {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()

	var qrs []*manta.Run
	for _, qr := range tcs.created {
		if qr.TaskID == taskID {
			qrs = append(qrs, qr)
		}
	}

	return qrs
}

// TotalRunsCreatedForTask returns the number of runs created for taskID.
func (tcs *TaskControlService) TotalRunsCreatedForTask(taskID manta.ID) int {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()

	return tcs.totalRunsCreated[taskID]
}

// PollForNumberCreated blocks for a small amount of time waiting for exactly the given count of created and unfinished runs for the given task ID.
// If the expected number isn't found in time, it returns an error.
//
// Because the scheduler and executor do a lot of state changes asynchronously, this is useful in test.
func (tcs *TaskControlService) PollForNumberCreated(taskID manta.ID, count int) ([]*manta.Run, error) {
	const numAttempts = 50
	actualCount := 0
	var created []*manta.Run
	for i := 0; i < numAttempts; i++ {
		time.Sleep(2 * time.Millisecond) // we sleep even on first so it becomes more likely that we catch when too many are produced.
		created = tcs.CreatedFor(taskID)
		actualCount = len(created)
		if actualCount == count {
			return created, nil
		}
	}
	return created, fmt.Errorf("did not see count of %dcs created run(s) for task with ID %s in time, instead saw %dcs", count, taskID, actualCount) // we return created anyways, to make it easier to debug
}

func (tcs *TaskControlService) FinishedRun(runID manta.ID) *manta.Run {
	tcs.mu.Lock()
	defer tcs.mu.Unlock()

	return tcs.finishedRuns[runID]
}

func (tcs *TaskControlService) FinishedRuns() []*manta.Run {
	rtn := []*manta.Run{}
	for _, run := range tcs.finishedRuns {
		rtn = append(rtn, run)
	}

	sort.Slice(rtn, func(i, j int) bool { return rtn[i].ScheduledFor.Before(rtn[j].ScheduledFor) })
	return rtn
}
