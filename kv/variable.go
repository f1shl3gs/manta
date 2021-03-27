package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
)

var (
	variableBucket         = []byte("variables")
	variableOrgIndexBucket = []byte("variableorgindex")
)

func (s *Service) FindVariableByID(ctx context.Context, id manta.ID) (*manta.Variable, error) {
	var (
		variable *manta.Variable
		err      error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		variable, err = s.findVariableByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return variable, nil
}

func (s *Service) findVariableByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Variable, error) {
	pk, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(variableBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	v := &manta.Variable{}
	err = v.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (s *Service) FindVariables(ctx context.Context, filter manta.VariableFilter) ([]*manta.Variable, error) {
	var (
		variables []*manta.Variable
		err       error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.OrgID != nil {
			variables, err = s.findVariablesByOrgID(ctx, tx, *filter.OrgID)
			return err
		}

		if filter.ID != nil {
			v, err := s.findVariableByID(ctx, tx, *filter.ID)
			if err != nil {
				return err
			}

			variables = append(variables, v)
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return variables, nil
}

func (s *Service) findVariablesByOrgID(ctx context.Context, tx Tx, orgID manta.ID) ([]*manta.Variable, error) {
	prefix, err := orgID.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(variableOrgIndexBucket)
	if err != nil {
		return nil, err
	}

	c, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix))
	if err != nil {
		return nil, err
	}

	keys := make([][]byte, 0)
	err = WalkCursor(ctx, c, func(k, v []byte) error {
		keys = append(keys, v)
		return nil
	})

	if err != nil {
		return nil, err
	}

	b, err = tx.Bucket(variableBucket)
	if err != nil {
		return nil, err
	}

	values, err := b.GetBatch(keys...)
	if err != nil {
		return nil, err
	}

	variables := make([]*manta.Variable, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}

		variable := &manta.Variable{}
		err = variable.Unmarshal(value)
		if err != nil {
			return nil, err
		}

		variables = append(variables, variable)
	}

	return variables, nil
}

func (s *Service) CreateVariable(ctx context.Context, v *manta.Variable) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.createVariable(ctx, tx, v)
	})
}

func (s *Service) createVariable(ctx context.Context, tx Tx, v *manta.Variable) error {
	now := time.Now()
	v.ID = s.idGen.ID()
	v.Created = now
	v.Updated = now

	return s.putVariable(ctx, tx, v)
}

func (s *Service) putVariable(ctx context.Context, tx Tx, v *manta.Variable) error {
	pk, err := v.ID.Encode()
	if err != nil {
		return err
	}

	// org index
	fk, err := v.OrgID.Encode()
	if err != nil {
		return err
	}

	indexKey := IndexKey(fk, pk)
	b, err := tx.Bucket(variableOrgIndexBucket)
	if err != nil {
		return err
	}

	err = b.Put(indexKey, pk)
	if err != nil {
		return err
	}

	// put variable itself
	val, err := v.Marshal()
	if err != nil {
		return err
	}

	b, err = tx.Bucket(variableBucket)
	if err != nil {
		return err
	}

	return b.Put(pk, val)
}

func (s *Service) PatchVariable(ctx context.Context, id manta.ID, udp *manta.VariableUpdate) (*manta.Variable, error) {
	var (
		variable *manta.Variable
		err      error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		variable, err = s.deleteVariable(ctx, tx, id)
		if err != nil {
			return err
		}

		variable.Updated = time.Now()
		udp.Apply(variable)

		return s.putVariable(ctx, tx, variable)
	})

	if err != nil {
		return nil, err
	}

	return variable, nil
}

func (s *Service) UpdateVariable(ctx context.Context, v *manta.Variable) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.updateVariable(ctx, tx, v)
	})
}

func (s *Service) updateVariable(ctx context.Context, tx Tx, v *manta.Variable) error {
	_, err := s.deleteVariable(ctx, tx, v.ID)
	if err != nil {
		return err
	}

	v.Updated = time.Now()
	return s.putVariable(ctx, tx, v)
}

func (s *Service) DeleteVariable(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		_, err := s.deleteVariable(ctx, tx, id)
		return err
	})
}

func (s *Service) deleteVariable(ctx context.Context, tx Tx, id manta.ID) (*manta.Variable, error) {
	variable, err := s.findVariableByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	pk, _ := variable.ID.Encode()

	// delete orgID index
	fk, err := variable.OrgID.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(variableOrgIndexBucket)
	if err != nil {
		return nil, err
	}

	if err = b.Delete(IndexKey(fk, pk)); err != nil {
		return nil, err
	}

	// delete variable itself
	b, err = tx.Bucket(variableBucket)
	if err != nil {
		return nil, err
	}

	if err = b.Delete(pk); err != nil {
		return nil, err
	}

	return variable, nil
}
