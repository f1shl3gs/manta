package kv

import (
    "context"
    "time"

    "github.com/f1shl3gs/manta"
)

var (
    ConfigurationBucket = []byte("configurations")

    ConfigurationOrgIndexBucket = []byte("configurationorgindex")
)

func (s *Service) CreateConfiguration(ctx context.Context, cf *manta.Configuration) error {
    now := time.Now()
    cf.ID = s.idGen.ID()
    cf.Created = now
    cf.Updated = now

    return s.kv.Update(ctx, func(tx Tx) error {
        return s.putConfiguration(ctx, tx, cf)
    })
}

func (s *Service) putConfiguration(ctx context.Context, tx Tx, cf *manta.Configuration) error {
    pk, err := cf.ID.Encode()
    if err != nil {
        return err
    }

    fk, err := cf.OrgID.Encode()
    if err != nil {
        return err
    }

    // org index
    indexKey := IndexKey(fk, pk)
    b, err := tx.Bucket(ConfigurationOrgIndexBucket)
    if err != nil {
        return err
    }

    if err = b.Put(indexKey, pk); err != nil {
        return err
    }

    // configuration
    b, err = tx.Bucket(ConfigurationBucket)
    if err != nil {
        return err
    }

    val, err := cf.Marshal()
    if err != nil {
        return err
    }

    return b.Put(pk, val)
}

func (s *Service) GetConfiguration(ctx context.Context, id manta.ID) (*manta.Configuration, error) {
    var (
        cf *manta.Configuration
        err error
    )

    err = s.kv.View(ctx, func(tx Tx) error {
        cf, err = s.getConfiguration(ctx, tx, id)
        return err
    })
    if err != nil {
        return nil, err
    }

    return cf, nil
}

func (s *Service) getConfiguration(ctx context.Context, tx Tx, id manta.ID) (*manta.Configuration, error) {
    b, err := tx.Bucket(ConfigurationBucket)
    if err != nil {
        return nil, err
    }

    key, err := id.Encode()
    if err != nil {
        return nil, err
    }

    val, err := b.Get(key)
    if err != nil {
        return nil, err
    }

    var cf = &manta.Configuration{}
    if err := cf.Unmarshal(val); err != nil {
        return nil, err
    }

    return cf, nil
}

func (s *Service) FindConfigurations(ctx context.Context, filter manta.ConfigurationFilter) ([]*manta.Configuration, error) {
    var (
        cs []*manta.Configuration
        err error
    )

    err = s.kv.View(ctx, func(tx Tx) error {
        cs, err = s.findConfigurations(ctx, tx, filter)
        return err
    })

    if err != nil {
        return nil, err
    }

    return cs, nil
}

func (s *Service) findConfigurations(ctx context.Context, tx Tx, filter manta.ConfigurationFilter) ([]*manta.Configuration, error) {
    var (
        err error
    )

    fk, err := filter.OrgID.Encode()
    if err != nil {
        return nil, err
    }

    b, err := tx.Bucket(ConfigurationOrgIndexBucket)
    if err != nil {
        return nil, err
    }

    cursor, err := b.ForwardCursor(fk)
    if err != nil {
        return nil, err
    }

    keys := make([][]byte, 0, 16)
    err = WalkCursor(ctx, cursor, func(k, v []byte) error {
        keys = append(keys, v)
        return nil
    })
    if err != nil {
        return nil, err
    }

    b, err = tx.Bucket(ConfigurationBucket)
    if err != nil {
        return nil, err
    }

    values, err := b.GetBatch(keys...)
    if err != nil {
        return nil, err
    }

    cs := make([]*manta.Configuration, 0, len(values))
    for _, val := range values {
        if val == nil {
            continue
        }

        c := &manta.Configuration{}
        if err = c.Unmarshal(val); err != nil {
            continue
        }

        cs = append(cs, c)
    }

    return cs, nil
}

func (s *Service) UpdateConfiguration(ctx context.Context, id manta.ID, upd manta.ConfigurationUpdate) error {
    return s.kv.Update(ctx, func(tx Tx) error {
        c, err := s.getConfiguration(ctx, tx, id)
        if err != nil {
            return err
        }

        upd.Apply(c)
        c.Updated = time.Now()

        return nil
    })
}

func (s *Service) DeleteConfiguration(ctx context.Context, id manta.ID) error {
    return s.kv.Update(ctx, func(tx Tx) error {
        return s.deleteConfiguration(ctx, tx, id)
    })
}

func (s *Service) deleteConfiguration(ctx context.Context, tx Tx, id manta.ID) error {
    c, err := s.getConfiguration(ctx, tx, id)
    if err != nil {
        return err
    }

    pk, err := c.ID.Encode()
    if err != nil {
        return err
    }

    fk, err := c.OrgID.Encode()
    if err != nil {
        return err
    }

    // delete configuration
    b, err := tx.Bucket(ConfigurationBucket)
    if err != nil {
        return err
    }

    if err = b.Delete(pk); err != nil {
        return err
    }

    // delete org index
    b, err = tx.Bucket(ConfigurationOrgIndexBucket)
    if err != nil {
        return err
    }

    if err = b.Delete(IndexKey(fk, pk)); err != nil {
        return err
    }

    return nil
}