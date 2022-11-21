package kv

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/f1shl3gs/manta"
)

var (
	RunsBucket         = []byte("runs")
	RunTaskIndexBucket = []byte("runtaskindex")
)

// CreateRun creates a run with a scheduled for time.
func (s *Service) CreateRun(ctx context.Context, taskID manta.ID, scheduledFor time.Time, runAt time.Time) (*manta.Run, error) {
	run := &manta.Run{
		ID:           s.idGen.ID(),
		TaskID:       taskID,
		ScheduledFor: scheduledFor,
		RunAt:        runAt,
		Status:       manta.RunScheduled,
		Logs:         []manta.RunLog{},
	}

	err := s.kv.Update(ctx, func(tx Tx) error {
		// store run
		b, err := tx.Bucket(RunsBucket)
		if err != nil {
			return err
		}

		data, err := json.Marshal(run)
		if err != nil {
			return err
		}

		pk, _ := run.ID.Encode()
		return b.Put(pk, data)
	})

	if err != nil {
		return nil, err
	}

	return run, nil
}

func putRun(tx Tx, run *manta.Run) error {
	pk, err := run.ID.Encode()
	if err != nil {
		return err
	}

	// task id index
	fk, err := run.TaskID.Encode()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(RunTaskIndexBucket)
	if err != nil {
		return err
	}

	index := IndexKey(fk, pk)
	if err = b.Put(index, pk); err != nil {
		return err
	}

	// store run
	b, err = tx.Bucket(RunsBucket)
	if err != nil {
		return err
	}

	data, err := json.Marshal(run)
	if err != nil {
		return err
	}

	return b.Put(pk, data)
}

func (s *Service) CurrentlyRunning(ctx context.Context, taskID manta.ID) ([]*manta.Run, error) {
	panic("implement me")
}

func (s *Service) ManualRuns(ctx context.Context, taskID manta.ID) ([]*manta.Run, error) {
	panic("implement me")
}

// StartManualRun pulls a manual run from the list and moves it to currently running.
func (s *Service) StartManualRun(ctx context.Context, taskID, runID manta.ID) (*manta.Run, error) {
	panic("implement me")
}

// FinishRun removes runID from the list of running tasks and if its `ScheduledFor` is later then last completed update it.
func (s *Service) FinishRun(ctx context.Context, taskID, runID manta.ID) (*manta.Run, error) {
	var (
		run *manta.Run
		err error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		run, err := findByID[manta.Run](tx, runID, RunsBucket)
		if err != nil {
			return err
		}

		task, err := findByID[manta.Task](tx, taskID, TasksBucket)
		if err != nil {
			return err
		}

		if run.Status == manta.RunFail {
			task.LatestFailure = run.ScheduledFor
		} else {
			task.LatestSuccess = run.ScheduledFor
		}

		task.LatestCompleted = run.ScheduledFor
		task.LastRunStatus = string(run.Status)
		task.LastRunError = func() string {
			if run.Status == manta.RunFail {
				if len(run.Logs) > 1 {
					return run.Logs[len(run.Logs)-2].Message
				} else if len(run.Logs) > 0 {
					return run.Logs[len(run.Logs)-1].Message
				}
			}

			return ""
		}()

		if err = putTask(tx, task); err != nil {
			return err
		}

		return removeRun(tx, taskID, runID)
	})

	if err != nil {
		return nil, err
	}

	return run, nil
}

func removeRun(tx Tx, taskID, runID manta.ID) error {
	pk, err := runID.Encode()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(RunsBucket)
	if err != nil {
		return err
	}

	if err = b.Delete(pk); err != nil {
		return err
	}

	// remove task index
	fk, err := taskID.Encode()
	if err != nil {
		return err
	}

	index := IndexKey(fk, pk)
	b, err = tx.Bucket(RunTaskIndexBucket)
	if err != nil {
		return err
	}

	return b.Delete(index)
}

// UpdateRunState sets the run state at the respective time.
func (s *Service) UpdateRunState(ctx context.Context, taskID, runID manta.ID, when time.Time, state manta.RunStatus) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		run, err := findByID[manta.Run](tx, runID, RunsBucket)
		if err != nil {
			return err
		}

		run.Status = state
		switch state {
		case manta.RunStarted:
			run.StartedAt = when
		case manta.RunSuccess, manta.RunFail, manta.RunCanceled:
			run.FinishedAt = when
		default:
			return errors.New("unknown run state")
		}

		return putRun(tx, run)
	})
}

// AddRunLog adds a file line to the run.
func (s *Service) AddRunLog(ctx context.Context, taskID, runID manta.ID, when time.Time, log string) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		run, err := findByID[manta.Run](tx, runID, RunsBucket)
		if err != nil {
			return err
		}

		run.Logs = append(run.Logs, manta.RunLog{
			RunID:   runID,
			Message: log,
			Time:    when.Format(time.RFC3339Nano),
		})

		return putRun(tx, run)
	})
}
