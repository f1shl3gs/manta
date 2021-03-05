package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	taskBucket           = []byte("tasks")
	taskOrgIndexBucket   = []byte("taskorgindex")
	taskOwnerIndexBucket = []byte("taskownerindex")
)

func (s *Service) FindTaskByID(ctx context.Context, id manta.ID) (*manta.Task, error) {
	var (
		task *manta.Task
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		task, err = s.findTaskByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) findTaskByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Task, error) {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	pk, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(taskBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	task := &manta.Task{}
	if err := task.Unmarshal(data); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) FindTasks(ctx context.Context, filter manta.TaskFilter) ([]*manta.Task, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		tasks []*manta.Task
		err   error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.OrgID != nil {
			tasks, err = s.findTasksByOrgID(ctx, tx, *filter.OrgID)
			if err != nil {
				return err
			}

			return nil
		}

		tasks, err = s.findAllTasks(ctx, tx)
		return err
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) findTasksByOrgID(ctx context.Context, tx Tx, orgID manta.ID) ([]*manta.Task, error) {
	prefix, err := orgID.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(taskOrgIndexBucket)
	if err != nil {
		return nil, err
	}

	cur, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix))
	if err != nil {
		return nil, err
	}

	keys := make([][]byte, 0, 16)
	err = WalkCursor(ctx, cur, func(k, v []byte) error {
		keys = append(keys, v)
		return nil
	})

	if err != nil {
		return nil, err
	}

	b, err = tx.Bucket(taskBucket)
	if err != nil {
		return nil, err
	}

	values, err := b.GetBatch(keys...)
	if err != nil {
		return nil, err
	}

	tasks := make([]*manta.Task, 0, len(values))
	for i := 0; i < len(values); i++ {
		task := &manta.Task{}
		if err = task.Unmarshal(values[i]); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Service) findAllTasks(ctx context.Context, tx Tx) ([]*manta.Task, error) {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	b, err := tx.Bucket(taskBucket)
	if err != nil {
		return nil, err
	}

	c, err := b.Cursor()
	if err != nil {
		return nil, err
	}

	tasks := make([]*manta.Task, 0)
	for k, v := c.First(); k != nil; k, v = c.Next() {
		task := &manta.Task{}
		if err := task.Unmarshal(v); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Service) CreateTask(ctx context.Context, task *manta.Task) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.createTask(ctx, tx, task)
	})
}

func (s *Service) createTask(ctx context.Context, tx Tx, task *manta.Task) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	task.ID = s.idGen.ID()
	task.Created = time.Now()
	task.Updated = time.Now()

	return s.putTask(ctx, tx, task)
}

func (s *Service) putTask(ctx context.Context, tx Tx, task *manta.Task) error {
	pk, err := task.ID.Encode()
	if err != nil {
		return err
	}

	fk, err := task.OrgID.Encode()
	if err != nil {
		return err
	}

	// organization index
	indexKey := IndexKey(fk, pk)
	b, err := tx.Bucket(taskOrgIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Put(indexKey, pk); err != nil {
		return err
	}

	// save task
	b, err = tx.Bucket(taskBucket)
	if err != nil {
		return err
	}

	data, err := task.Marshal()
	if err != nil {
		return err
	}

	return b.Put(pk, data)
}

func (s *Service) UpdateTask(ctx context.Context, id manta.ID, udp manta.TaskUpdate) (*manta.Task, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		task *manta.Task
		err  error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		task, err = s.deleteTask(ctx, tx, id)
		if err != nil {
			return err
		}

		if udp.Status != nil {
			task.Status = *udp.Status
		}

		if udp.LatestScheduled != nil {
			task.LatestScheduled = *udp.LatestScheduled
		}

		if udp.LatestCompleted != nil {
			task.LatestCompleted = *udp.LatestCompleted
		}

		if udp.LatestSuccess != nil {
			task.LatestSuccess = *udp.LatestSuccess
		}

		if udp.LatestFailure != nil {
			task.LatestFailure = *udp.LatestFailure
		}

		if udp.LastRunError != nil {
			task.LastRunError = *udp.LastRunError
		}

		task.Updated = time.Now()

		return s.putTask(ctx, tx, task)
	})

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) DeleteTask(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		_, err := s.deleteTask(ctx, tx, id)
		return err
	})
}

func (s *Service) deleteTask(ctx context.Context, tx Tx, id manta.ID) (*manta.Task, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	task, err := s.findTaskByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	pk, err := id.Encode()
	if err != nil {
		return nil, err
	}

	// delete organization index
	fk, err := task.OrgID.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(taskOrgIndexBucket)
	if err != nil {
		return nil, err
	}

	index := IndexKey(fk, pk)
	if err := b.Delete(index); err != nil {
		return nil, err
	}

	// delete task
	b, err = tx.Bucket(taskBucket)
	if err != nil {
		return nil, err
	}

	if err := b.Delete(pk); err != nil {
		return nil, err
	}

	return task, nil
}
