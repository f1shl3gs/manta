package kv

import (
	"context"
	"encoding/json"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	DashboardsBucket        = []byte("dashboards")
	DashboardOrgIndexBucket = []byte("dashboardorgindex")

	CellsBucket              = []byte("cells")
	CellDashboardIndexBucket = []byte("celldashboardindex")
)

func (s *Service) FindDashboardByID(ctx context.Context, id manta.ID) (*manta.Dashboard, error) {
	var (
		d   *manta.Dashboard
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		d, err = s.findDashboardByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *Service) findDashboardByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Dashboard, error) {
	pk, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(DashboardsBucket)
	if err != nil {
		return nil, err
	}

	val, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	d := &manta.Dashboard{}
	if err = json.Unmarshal(val, d); err != nil {
		return nil, err
	}

	return d, nil
}

func (s *Service) FindDashboards(ctx context.Context, filter manta.DashboardFilter) ([]*manta.Dashboard, error) {
	var (
		list []*manta.Dashboard
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.OrganizationID != nil {
			list, err = s.findDashboardByOrg(ctx, tx, *filter.OrganizationID)
		} else {
			list, err = s.findAllDashboards(ctx, tx)
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

// todo: for now
func (s *Service) findAllDashboards(ctx context.Context, tx Tx) ([]*manta.Dashboard, error) {
	b, err := tx.Bucket(DashboardsBucket)
	if err != nil {
		return nil, err
	}

	c, err := b.ForwardCursor(nil)
	if err != nil {
		return nil, err
	}

	list := make([]*manta.Dashboard, 0, 8)
	err = WalkCursor(ctx, c, func(k, v []byte) error {
		dash := &manta.Dashboard{}
		if err := json.Unmarshal(v, dash); err != nil {
			return err
		}

		list = append(list, dash)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Service) findDashboardByOrg(ctx context.Context, tx Tx, orgID manta.ID) ([]*manta.Dashboard, error) {
	fk, err := orgID.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(DashboardOrgIndexBucket)
	if err != nil {
		return nil, err
	}

	c, err := b.ForwardCursor(fk)
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

	b, err = tx.Bucket(DashboardsBucket)
	if err != nil {
		return nil, err
	}

	values, err := b.GetBatch(keys...)
	if err != nil {
		return nil, err
	}

	ds := make([]*manta.Dashboard, 0, len(values))
	for _, val := range values {
		if val == nil {
			continue
		}

		d := &manta.Dashboard{}
		if err = json.Unmarshal(val, d); err != nil {
			return nil, err
		}

		ds = append(ds, d)
	}

	return ds, nil
}

func (s *Service) CreateDashboard(ctx context.Context, d *manta.Dashboard) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.createDashboard(ctx, tx, d)
	})
}

func (s *Service) createDashboard(ctx context.Context, tx Tx, d *manta.Dashboard) error {
	now := time.Now()

	d.ID = s.idGen.ID()
	d.Created = now
	d.Updated = now

	return s.putDashboard(ctx, tx, d)
}

func (s *Service) putDashboard(ctx context.Context, tx Tx, d *manta.Dashboard) error {
	pk, err := d.ID.Encode()
	if err != nil {
		return err
	}

	fk, err := d.OrgID.Encode()
	if err != nil {
		return err
	}

	// org index
	indexKey := IndexKey(fk, pk)
	b, err := tx.Bucket(DashboardOrgIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Put(indexKey, pk); err != nil {
		return err
	}

	// dashboard
	b, err = tx.Bucket(DashboardsBucket)
	if err != nil {
		return err
	}
	val, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return b.Put(pk, val)
}

func (s *Service) UpdateDashboard(ctx context.Context, id manta.ID, udp manta.DashboardUpdate) (*manta.Dashboard, error) {
	var (
		dash *manta.Dashboard
		err  error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		dash, err = s.findDashboardByID(ctx, tx, id)
		if err != nil {
			return err
		}

		udp.Apply(dash)
		dash.Updated = time.Now()

		return s.putDashboard(ctx, tx, dash)
	})

	if err != nil {
		return nil, err
	}

	return dash, nil
}

func (s *Service) AddDashboardCell(ctx context.Context, id manta.ID, cell *manta.Cell) error {
	var (
		dash *manta.Dashboard
		err  error
	)

	// initial
	cell.ID = s.idGen.ID()

	err = s.kv.Update(ctx, func(tx Tx) error {
		dash, err = s.findDashboardByID(ctx, tx, id)
		if err != nil {
			return err
		}

		dash.Cells = append(dash.Cells, *cell)
		dash.Updated = time.Now()

		return s.putDashboard(ctx, tx, dash)
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveDashboardCell(ctx context.Context, did, cellID manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		dash, err := s.findDashboardByID(ctx, tx, did)
		if err != nil {
			return err
		}

		for i := 0; i < len(dash.Cells); i++ {
			if dash.Cells[i].ID == cellID {
				dash.Cells = append(dash.Cells[:i], dash.Cells[i+1:]...)
				break
			}
		}

		dash.Updated = time.Now()

		return s.putDashboard(ctx, tx, dash)
	})
}

func (s *Service) UpdateDashboardCell(ctx context.Context, did, cellID manta.ID, udp manta.DashboardCellUpdate) (*manta.Cell, error) {
	var (
		cell *manta.Cell
		err  error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		dash, err := s.findDashboardByID(ctx, tx, did)
		if err != nil {
			return err
		}

		for i := 0; i < len(dash.Cells); i++ {
			if dash.Cells[i].ID == cellID {
				cell = &dash.Cells[i]

				udp.Apply(cell)
				break
			}
		}

		dash.Updated = time.Now()
		return s.putDashboard(ctx, tx, dash)
	})

	if err != nil {
		return nil, err
	}

	return cell, nil
}

func (s *Service) DeleteDashboard(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteDashboard(ctx, tx, id)
	})
}

func (s *Service) deleteDashboard(ctx context.Context, tx Tx, id manta.ID) error {
	pk, err := id.Encode()
	if err != nil {
		return err
	}

	d, err := s.findDashboardByID(ctx, tx, id)
	if err != nil {
		return err
	}

	// delete org index
	fk, err := d.OrgID.Encode()
	if err != nil {
		return err
	}

	b, err := tx.Bucket(DashboardOrgIndexBucket)
	if err != nil {
		return err
	}

	indexKey := IndexKey(fk, pk)
	if err = b.Delete(indexKey); err != nil {
		return err
	}

	// delete dashboard
	b, err = tx.Bucket(DashboardsBucket)
	if err != nil {
		return err
	}

	return b.Delete(pk)
}

func (s *Service) ReplaceDashboardCells(ctx context.Context, did manta.ID, cells []manta.Cell) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.kv.Update(ctx, func(tx Tx) error {
		dash, err := s.findDashboardByID(ctx, tx, did)
		if err != nil {
			return err
		}

		dash.Cells = cells
		dash.Updated = time.Now()

		return s.putDashboard(ctx, tx, dash)
	})
}

func (s *Service) findDashboardCell(ctx context.Context, tx Tx, did, cid manta.ID) (*manta.Cell, error) {
	dash, err := s.findDashboardByID(ctx, tx, did)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(dash.Cells); i++ {
		if dash.Cells[i].ID == cid {
			return &dash.Cells[i], nil
		}
	}

	return nil, ErrKeyNotFound
}

func (s *Service) GetDashboardCell(ctx context.Context, did, cid manta.ID) (*manta.Cell, error) {
	var (
		cell *manta.Cell
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		cell, err = s.findDashboardCell(ctx, tx, did, cid)
		return err
	})

	if err != nil {
		return nil, err
	}

	return cell, nil
}

// func (s *Service) GetDashboardCellView(ctx context.Context, did, cid manta.ID) (*manta.View, error) {
// 	var (
// 		cell *manta.Prop
// 		err error
// 	)
//
// 	err = s.kv.View(ctx, func(tx Tx) error {
// 		cell, err = s.findDashboardCell(ctx, tx, did, cid)
// 		if err != nil {
// 			return err
// 		}
//
// 		return nil
// 	})
// }
//
// func (s *Service) UpdateDashboardCellView(ctx context.Context, did, cid manta.ID, udp manta.ViewUpdate) (*manta.View, error) {
// 	panic("implement me")
// }
