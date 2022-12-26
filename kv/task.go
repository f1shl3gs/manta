package kv

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/f1shl3gs/manta"
)

const TaskDefaultPageSize = 100

var (
	TasksBucket          = []byte("tasks")
	TaskOrgIndexBucket   = []byte("taskorgindex")
	TaskOwnerIndexBucket = []byte("taskownerindex")

	// Errors
	ErrInvalidTaskID = &manta.Error{
		Code: manta.EInvalid,
		Msg:  "invalid task id",
	}
)

func ErrInternalTaskService(err error) *manta.Error {
	return &manta.Error{
		Code: manta.EInternal,
		Msg:  "unexpected error in tasks",
		Err:  err,
	}
}

// FindTaskByID returns a single task by id
func (s *Service) FindTaskByID(ctx context.Context, id manta.ID) (*manta.Task, error) {
	var (
		task = &manta.Task{}
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		task, err = findByID[manta.Task](tx, id, TasksBucket)
		return err
	})

	if err != nil {
		return nil, err
	}

	return task, nil
}

// FindTasks returns all tasks which match the filter
func (s *Service) FindTasks(ctx context.Context, filter manta.TaskFilter) ([]*manta.Task, error) {
	var (
		tasks []*manta.Task
		err   error
	)

	if filter.OrgID != nil {
		err = s.kv.View(ctx, func(tx Tx) error {
			tasks, err = findOrgIndexed[manta.Task](ctx, tx, *filter.OrgID, TasksBucket, TaskOrgIndexBucket)
			return err
		})

		return tasks, err
	}

	if filter.OwnerID != nil {
		fk, err := filter.OwnerID.Encode()
		if err != nil {
			return nil, err
		}

		err = s.kv.View(ctx, func(tx Tx) error {
			b, err := tx.Bucket(TaskOwnerIndexBucket)
			if err != nil {
				return err
			}

			cursor, err := b.Cursor()
			if err != nil {
				return err
			}

			var keys [][]byte
			for k, v := cursor.First(); k != nil && bytes.HasPrefix(k, fk); k, v = cursor.Next() {
				keys = append(keys, v)
			}

			b, err = tx.Bucket(TasksBucket)
			if err != nil {
				return err
			}

			values, err := b.GetBatch(keys...)
			if err != nil {
				return err
			}

			for _, value := range values {
				task := &manta.Task{}
				err = json.Unmarshal(value, task)
				if err != nil {
					return err
				}

				tasks = append(tasks, task)
			}

			return nil
		})

		return tasks, err
	}

	// list all tasks
	err = s.kv.View(ctx, func(tx Tx) error {
		b, err := tx.Bucket(TasksBucket)
		if err != nil {
			return err
		}

		cursor, err := b.Cursor()
		if err != nil {
			return err
		}

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			task := new(manta.Task)
			err = json.Unmarshal(v, task)
			if err != nil {
				return err
			}

			tasks = append(tasks, task)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) FindRuns(ctx context.Context, filter manta.RunFilter) ([]*manta.Run, int, error) {
	if filter.Limit == 0 {
		filter.Limit = TaskDefaultPageSize
	}

	if filter.Limit < 0 || filter.Limit > TaskDefaultPageSize {
		return nil, 0, manta.ErrOutOfBoundsLimit
	}

	var (
	    list = make([]*manta.Run, 0)
        err error
	)
	err = s.kv.View(ctx, func(tx Tx) error {
        list, _, err = findRuns(tx, filter)
        return err
	})

	if err != nil {
		return nil, 0, err
	}

	return list, len(list), nil
}

func findRuns(tx Tx, filter manta.RunFilter) ([]*manta.Run, int, error) {
    var runs []*manta.Run

    taskKey, err := filter.Task.Encode()
    if err != nil {
        return nil, 0, ErrInvalidTaskID
    }

    dataBucket, err := tx.Bucket(RunsBucket)
    if err != nil {
        return nil, 0, err
    }

    indexBucket, err := tx.Bucket(RunTaskIndexBucket)
    if err != nil {
        return nil, 0, err
    }

    cursor, err := indexBucket.Cursor(WithCursorHintPrefix(string(taskKey)))
    if err != nil {
        return nil, 0, err
    }

    k, v := cursor.Seek(taskKey)
    for {
        if k == nil || !bytes.HasPrefix(k, taskKey) {
            break
        }

        // TODO: filter with after and before

        value, err := dataBucket.Get(v)
        if err != nil {
            return nil, 0, err
        }

        run := &manta.Run{}
        if err := json.Unmarshal(value, run); err != nil {
            return nil, 0, ErrInternalTaskService(err)
        }

        runs = append(runs, run)
        k, v = cursor.Next()
    }

    return runs, len(runs), nil
}

// CreateTask creates a task
func (s *Service) CreateTask(ctx context.Context, task *manta.Task) error {
	var (
		now = time.Now()
	)

	task.ID = s.idGen.ID()
	task.Created = now
	task.Updated = now

	return s.kv.Update(ctx, func(tx Tx) error {
		return putTask(tx, task)
	})
}

func putTask(tx Tx, task *manta.Task) error {
	pk, _ := task.ID.Encode()

	// store org index
	index, err := indexIDKey(task.ID, task.OrgID)
	if err != nil {
		return err
	}

	b, err := tx.Bucket(TaskOrgIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Put(index, pk); err != nil {
		return err
	}

	// store owner index
	index, err = indexIDKey(task.ID, task.OwnerID)
	if err != nil {
		return err
	}

	b, err = tx.Bucket(TaskOwnerIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Put(index, pk); err != nil {
		return err
	}

	// store task
	b, err = tx.Bucket(TasksBucket)
	if err != nil {
		return err
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return b.Put(pk, data)
}

// UpdateTask updates a single task with a patch
func (s *Service) UpdateTask(ctx context.Context, id manta.ID, udp manta.TaskUpdate) (*manta.Task, error) {
	var (
		task *manta.Task
		err  error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		task, err = findByID[manta.Task](tx, id, TasksBucket)
		if err != nil {
			return err
		}

		udp.Apply(task)
		task.Updated = time.Now()

		return putTask(tx, task)
	})

	return task, err
}

func updateTask(tx Tx, id manta.ID, upd manta.TaskUpdate) (*manta.Task, error) {
	task, err := findByID[manta.Task](tx, id, TasksBucket)
	if err != nil {
		return nil, err
	}

	upd.Apply(task)
	task.Updated = time.Now()

	err = putTask(tx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask delete a single task by ID
func (s *Service) DeleteTask(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return deleteTask(tx, id)
	})
}

func deleteTask(tx Tx, id manta.ID) error {
    return deleteOrgIndexed[manta.Task](tx, id, TasksBucket, TaskOrgIndexBucket)
}
