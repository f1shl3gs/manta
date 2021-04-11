package kv

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	runBucket          = []byte("runs")
	runTaskIndexBucket = []byte("runtaskindex")
)

func (s *Service) CreateRun(ctx context.Context, taskID manta.ID, scheduledFor time.Time, runAt time.Time) (*manta.Run, error) {
	var (
		run *manta.Run
		err error
	)

	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	err = s.kv.Update(ctx, func(tx Tx) error {
		run, err = s.createRun(ctx, tx, taskID, scheduledFor, runAt)
		return err
	})

	if err != nil {
		return nil, err
	}

	return run, nil
}

func (s *Service) createRun(ctx context.Context, tx Tx, taskID manta.ID, scheduledFor time.Time, runAt time.Time) (*manta.Run, error) {
	run := &manta.Run{
		ID:           s.idGen.ID(),
		TaskID:       taskID,
		ScheduledFor: scheduledFor,
		RunAt:        runAt,
		Status:       manta.RunScheduled.String(),
		Logs:         []manta.RunLog{},
	}

	data, err := run.Marshal()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(runBucket)
	if err != nil {
		return nil, err
	}

	pk, err := run.ID.Encode()
	if err != nil {
		return nil, err
	}

	err = b.Put(pk, data)
	if err != nil {
		return nil, err
	}

	return run, nil
}

func (s *Service) putTaskRun(ctx context.Context, tx Tx, run *manta.Run) error {
	data, err := run.Marshal()
	if err != nil {
		return err
	}

	pk, err := run.ID.Encode()
	if err != nil {
		return err
	}

	// task id index
	fk, err := run.TaskID.Encode()
	if err != nil {
		return err
	}
	indexKey := IndexKey(fk, pk)

	b, err := tx.Bucket(runTaskIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Put(indexKey, pk); err != nil {
		return err
	}

	// put itself
	b, err = tx.Bucket(runBucket)
	if err != nil {
		return err
	}

	if err = b.Put(pk, data); err != nil {
		return err
	}

	return nil
}

func (s *Service) CurrentlyRunning(ctx context.Context, taskID manta.ID) ([]*manta.Run, error) {
	panic("implement me")
}

func (s *Service) ManualRuns(ctx context.Context, taskID manta.ID) ([]*manta.Run, error) {
	panic("implement me")
}

func (s *Service) StartManualRun(ctx context.Context, taskID, runID manta.ID) (*manta.Run, error) {
	panic("implement me")
}

func (s *Service) FinishRun(ctx context.Context, taskID, runID manta.ID) (*manta.Run, error) {
	var (
		run *manta.Run
		err error
	)

	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	err = s.kv.Update(ctx, func(tx Tx) error {
		run, err = s.findRunByID(ctx, tx, runID)
		if err != nil {
			return err
		}

		task, err := s.findTaskByID(ctx, tx, taskID)
		if err != nil {
			return err
		}

		if run.Status == "failed" {
			task.LatestFailure = run.ScheduledFor
		} else {
			task.LatestSuccess = run.ScheduledFor
		}

		task.LatestCompleted = run.ScheduledFor
		task.LastRunStatus = run.Status

		task.LastRunError = func() string {
			if run.Status != "failed" {
				return ""
			}

			var b strings.Builder
			b.WriteString(task.LastRunError)
			b.WriteByte('\n')

			for _, l := range run.Logs {
				b.WriteString(l.Time)
				b.WriteString(": ")
				b.WriteString(l.Message)
				b.WriteByte('\n')
			}

			return b.String()

			/*
				if len(run.Logs) > 1 {
					return run.Logs[len(run.Logs)-2].Message
				} else {
					return run.Logs[len(run.Logs)-1].Message
				}*/

			// todo: handle logs
			// return task.LastRunError
		}()

		if err = s.putTask(ctx, tx, task); err != nil {
			return err
		}

		return s.removeRun(ctx, tx, taskID, runID)
	})

	if err != nil {
		return nil, err
	}

	return run, nil
}

func (s *Service) findRunByID(ctx context.Context, tx Tx, runID manta.ID) (*manta.Run, error) {
	pk, err := runID.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(runBucket)
	if err != nil {
		return nil, err
	}

	val, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	run := &manta.Run{}
	err = run.Unmarshal(val)
	if err != nil {
		return nil, err
	}

	return run, nil
}

func (s *Service) UpdateRunState(
	ctx context.Context,
	taskID, runID manta.ID,
	when time.Time,
	state manta.RunStatus,
) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		run, err := s.findRunByID(ctx, tx, runID)
		if err != nil {
			return err
		}

		run.Status = state.String()
		switch state {
		case manta.RunStarted:
			run.StartedAt = when
		case manta.RunSuccess, manta.RunFail, manta.RunCanceled:
			run.FinishedAt = when
		default:
			return errors.New("unknown run state")
		}

		return s.putTaskRun(ctx, tx, run)
	})
}

func (s *Service) AddRunLog(ctx context.Context, taskID, runID manta.ID, when time.Time, log string) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.kv.Update(ctx, func(tx Tx) error {
		run, err := s.findRunByID(ctx, tx, runID)
		if err != nil {
			return err
		}

		run.Logs = append(run.Logs, manta.RunLog{
			RunID:   runID,
			Time:    when.Format(time.RFC3339Nano),
			Message: log,
		})

		return s.putTaskRun(ctx, tx, run)
	})
}

func (s *Service) removeRun(ctx context.Context, tx Tx, taskID, runID manta.ID) error {
	// remove run
	pk, err := runID.Encode()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(runBucket)
	if err != nil {
		return err
	}

	if err := b.Delete(pk); err != nil {
		return err
	}

	// remove task index
	fk, _ := taskID.Encode()
	indexKey := IndexKey(fk, pk)
	b, err = tx.Bucket(runTaskIndexBucket)
	if err != nil {
		return err
	}

	return b.Delete(indexKey)
}
