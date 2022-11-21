package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
)

var (
	TasksBucket          = []byte("tasks")
	TaskOrgIndexBucket   = []byte("taskorgindex")
	TaskOwnerIndexBucket = []byte("taskownerindex")
)

// FindTaskByID returns a single task by id
func (s *Service) FindTaskByID(ctx context.Context, id manta.ID) (*manta.Task, error) {
	var (
		task = &manta.Task{}
		err  error
	)

	key, err := id.Encode()
	if err != nil {
		return nil, err
	}

	err = s.kv.View(ctx, func(tx Tx) error {
		b, err := tx.Bucket(TasksBucket)
		if err != nil {
			return err
		}

		val, err := b.Get(key)
		if err != nil {
			return err
		}

		return task.Unmarshal(val)
	})

	if err != nil {
		return nil, err
	}

	return task, nil
}

// FindTasks returns all tasks which match the filter
func (s *Service) FindTasks(ctx context.Context, filter manta.TaskFilter) ([]*manta.Task, error) {
	var (
		ts  []*manta.Task
		err error
	)

	if filter.OrgID != nil {
		err = s.kv.View(ctx, func(tx Tx) error {
			ts, err = findOrgIndexed[manta.Task](ctx, tx, *filter.OrgID, TasksBucket, TaskOrgIndexBucket)
			return err
		})

		return ts, err
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
			err = task.Unmarshal(v)
			if err != nil {
				return err
			}

			ts = append(ts, task)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ts, nil
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

	data, err := task.Marshal()
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
		task, err = getOrgIndexed[manta.Task](tx, id, TasksBucket)
		if err != nil {
			return err
		}

		udp.Apply(task)
		task.Updated = time.Now()

		return putTask(tx, task)
	})

	return task, err
}

// DeleteTask delete a single task by ID
func (s *Service) DeleteTask(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		pk, err := id.Encode()
		if err != nil {
			return err
		}

		// retrieve task
		b, err := tx.Bucket(TasksBucket)
		if err != nil {
			return err
		}

		val, err := b.Get(pk)
		if err != nil {
			return err
		}

		var task manta.Task
		if err = task.Unmarshal(val); err != nil {
			return err
		}

		// delete org index
		fk, err := task.OrgID.Encode()
		if err != nil {
			return err
		}

		index := IndexKey(fk, pk)
		b, err = tx.Bucket(TaskOrgIndexBucket)
		if err != nil {
			return err
		}

		if err = b.Delete(index); err != nil {
			return err
		}

		// delete owner index
		fk, err = task.OwnerID.Encode()
		if err != nil {
			return err
		}

		index = IndexKey(fk, pk)
		b, err = tx.Bucket(TaskOwnerIndexBucket)
		if err != nil {
			return err
		}

		if err = b.Delete(index); err != nil {
			return err
		}

		// delete the task
		b, err = tx.Bucket(TasksBucket)
		if err != nil {
			return err
		}

		return b.Delete(pk)
	})
}
