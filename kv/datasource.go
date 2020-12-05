package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/pkg/errors"
)

var (
	datasourceBucket          = []byte("datasource")
	datasourceNameIndexBucket = []byte("datasourcenameindex")
	datasourceOrgIndexBucket  = []byte("datasourceorgindex")
)

func (s *Service) FindDatasourceByID(ctx context.Context, id manta.ID) (*manta.Datasource, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		ds  *manta.Datasource
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		ds, err = s.findDatasourceByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (s *Service) findDatasourceByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Datasource, error) {
	pk, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(datasourceBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	ds := &manta.Datasource{}
	if err := ds.Unmarshal(data); err != nil {
		return nil, err
	}

	return ds, nil
}

func (s *Service) FindDatasource(ctx context.Context, filter manta.DatasourceFilter) (*manta.Datasource, error) {
	var (
		ds  *manta.Datasource
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.Name != nil {
			ds, err = s.findDatasourceByName(ctx, tx, *filter.Name)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (s *Service) findDatasourceByName(ctx context.Context, tx Tx, name string) (*manta.Datasource, error) {
	b, err := tx.Bucket(datasourceNameIndexBucket)
	if err != nil {
		return nil, err
	}

	pk, err := b.Get([]byte(name))
	if err != nil {
		return nil, err
	}

	var id manta.ID
	if err = id.Decode(pk); err != nil {
		return nil, err
	}

	return s.findDatasourceByID(ctx, tx, id)
}

func (s *Service) FindDatasources(ctx context.Context, filter manta.DatasourceFilter) ([]*manta.Datasource, error) {
	var (
		datasources []*manta.Datasource
		err         error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.OrgID != nil {
			datasources, err = s.findDatasourcesByOrgID(ctx, tx, *filter.OrgID)
		} else {
			// todo: this should be remove in future, cause list all datasource is not allowed
			datasources, err = s.findDatasources(ctx, tx)
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return datasources, nil
}

func (s *Service) findDatasources(ctx context.Context, tx Tx) ([]*manta.Datasource, error) {
	b, err := tx.Bucket(datasourceBucket)
	if err != nil {
		return nil, err
	}

	c, err := b.Cursor()
	if err != nil {
		return nil, err
	}

	var list []*manta.Datasource
	for k, v := c.First(); k != nil; k, v = c.Next() {
		ds := &manta.Datasource{}
		if err = ds.Unmarshal(v); err != nil {
			return nil, err
		}

		list = append(list, ds)
	}

	return list, nil
}

func (s *Service) findDatasourcesByOrgID(ctx context.Context, tx Tx, id manta.ID) ([]*manta.Datasource, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.LogKV("orgID", id.String())

	b, err := tx.Bucket(datasourceOrgIndexBucket)
	if err != nil {
		return nil, err
	}

	prefix, err := id.Encode()
	if err != nil {
		return nil, err
	}

	c, err := b.ForwardCursor(prefix)
	if err != nil {
		return nil, err
	}

	keys := make([][]byte, 0, 4)
	err = WalkCursor(ctx, c, func(k, v []byte) error {
		keys = append(keys, v)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// read datasources
	b, err = tx.Bucket(datasourceBucket)
	if err != nil {
		return nil, err
	}

	vals, err := b.GetBatch(keys...)
	if err != nil {
		return nil, err
	}

	if len(vals) == 0 {
		return nil, nil
	}

	list := make([]*manta.Datasource, 0, len(vals))
	for i := 0; i < len(vals); i++ {
		ds := &manta.Datasource{}
		data := vals[i]
		if data == nil {
			continue
		}

		if err = ds.Unmarshal(data); err != nil {
			return nil, err
		}

		list = append(list, ds)
	}

	return list, nil
}

func (s *Service) CreateDatasource(ctx context.Context, ds *manta.Datasource) error {
	err := s.kv.Update(ctx, func(tx Tx) error {
		return s.createDatasource(ctx, tx, ds)
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) createDatasource(ctx context.Context, tx Tx, ds *manta.Datasource) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	ds.ID = s.idGen.ID()
	ds.Created = time.Now()
	ds.Modified = time.Now()

	b, err := tx.Bucket(datasourceNameIndexBucket)
	if err != nil {
		return err
	}

	_, err = b.Get([]byte(ds.Name))
	if err == ErrKeyNotFound {
		return s.putDatasource(ctx, tx, ds)
	}

	if err == nil {
		return ErrIndexConflict
	}

	return err
}

func (s *Service) putDatasource(ctx context.Context, tx Tx, ds *manta.Datasource) error {
	fk := []byte(ds.Name)
	pk, err := ds.ID.Encode()
	if err != nil {
		return err
	}

	// name index
	b, err := tx.Bucket(datasourceNameIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Put(fk, pk); err != nil {
		return errors.Wrap(err, "write name index failed")
	}

	// org index
	b, err = tx.Bucket(datasourceOrgIndexBucket)
	if err != nil {
		return err
	}

	fk, err = ds.OrgID.Encode()
	if err != nil {
		return err
	}

	key := IndexKey(fk, pk)
	if err = b.Put(key, pk); err != nil {
		return err
	}

	// put itself
	data, err := ds.Marshal()
	if err != nil {
		return err
	}

	b, err = tx.Bucket(datasourceBucket)
	if err != nil {
		return err
	}

	return b.Put(pk, data)
}

func (s *Service) UpdateDatasource(ctx context.Context, id manta.ID, udp manta.DatasourceUpdate) (*manta.Datasource, error) {
	panic("implement me")
}

func (s *Service) DeleteDatasource(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteDatasource(ctx, tx, id)
	})
}

func (s *Service) deleteDatasource(ctx context.Context, tx Tx, id manta.ID) error {
	ds, err := s.findDatasourceByID(ctx, tx, id)
	if err != nil {
		return err
	}

	fk := []byte(ds.Name)
	pk, _ := id.Encode()

	// name index
	b, err := tx.Bucket(datasourceNameIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Delete(fk); err != nil {
		return err
	}

	// orgID index
	fk, _ = ds.OrgID.Encode()
	indexKey := IndexKey(fk, pk)
	b, err = tx.Bucket(datasourceOrgIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Delete(indexKey); err != nil {
		return err
	}

	// delete itself
	b, err = tx.Bucket(datasourceBucket)
	if err != nil {
		return err
	}

	if err = b.Delete(pk); err != nil {
		return err
	}

	return nil
}
