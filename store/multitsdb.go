package store

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/log"
	"github.com/f1shl3gs/manta/pkg/multierr"
)

type MultiTSDB struct {
	logger *zap.Logger
	// bucket objstore.Bucket

	dataDir               string
	tsdbOpts              *tsdb.Options
	labels                labels.Labels
	allowOutOfOrderUpload bool

	reg     prometheus.Registerer
	mtx     sync.RWMutex
	tenants map[manta.ID]*Tenant
}

func (m *MultiTSDB) TenantStorage(ctx context.Context, id manta.ID) (storage.Storage, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	ts := m.tenants[id]
	if ts == nil {
		return nil, ErrUnknownTenantStorage
	}

	return ts.readyS.Get(), nil
}

func NewMultiTSDB(
	dataDir string,
	logger *zap.Logger,
	reg prometheus.Registerer,
	tsdbOpts *tsdb.Options,
	labels labels.Labels,
	// bucket objstore.Bucket,
	allowOutOfOrderUpload bool,
) *MultiTSDB {
	return &MultiTSDB{
		logger:   logger,
		dataDir:  dataDir,
		reg:      reg,
		tsdbOpts: tsdbOpts,
		labels:   labels,
		// 	bucket:                bucket,
		allowOutOfOrderUpload: allowOutOfOrderUpload,
		tenants:               make(map[manta.ID]*Tenant),
	}
}

func (m *MultiTSDB) Open() error {
	if err := os.MkdirAll(m.dataDir, 0777); err != nil {
		return err
	}

	files, err := os.ReadDir(m.dataDir)
	if err != nil {
		return err
	}

	var g errgroup.Group
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		var id manta.ID
		err = id.DecodeFromString(f.Name())
		if err != nil {
			continue
		}

		g.Go(func() error {
			_, err := m.getOrLoadTenant(id, true)
			return err
		})
	}

	return g.Wait()
}

func (m *MultiTSDB) Flush() error {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	errs := &multierr.SyncedErrors{}
	wg := &sync.WaitGroup{}

	for id, tenant := range m.tenants {
		db := tenant.readyStorage().Get()
		if db == nil {
			m.logger.Error("flushing tsdb failed, not ready",
				zap.String("tenant", id.String()))
			continue
		}

		m.logger.Info("flushing tsdb",
			zap.String("tenant", id.String()))
		wg.Add(1)
		go func() {
			defer wg.Done()

			head := db.Head()
			if err := db.CompactHead(tsdb.NewRangeHead(head, head.MinTime(), head.MaxTime()-1)); err != nil {
				errs.Add(err)
			}
		}()
	}

	wg.Wait()

	return errs.Unwrap()
}

func (m *MultiTSDB) Close() error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	errs := multierr.Errors{}
	for id, tenant := range m.tenants {
		db := tenant.readyStorage().Get()
		if db == nil {
			m.logger.Error("closing tsdb failed, not ready",
				zap.String("tenant", id.String()))
			continue
		}

		if err := db.Close(); err != nil {
			errs.Add(db.Close())
		}
	}

	return errs.Unwrap()
}

func (m *MultiTSDB) Sync(ctx context.Context) (int, error) {
	return 0, nil
	/*	if m.bucket == nil {
			return 0, errors.New("bucket is not specified, sync should not be invoked")
		}

		var (
			uploaded atomic.Int64
			errs     = &multierr.SyncedErrors{}
			wg       = sync.WaitGroup{}
		)

		m.mtx.RLock()
		defer m.mtx.RUnlock()

		for tid, tenant := range m.tenants {
			m.logger.Debug("uploading block for tenant",
				zap.String("tenant", tid.String()))

			s := tenant.shipper()
			if s == nil {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				up, err := s.Sync(ctx)
				if err != nil {
					errs.Add(errors.Wrap(err, "upload"))
				}

				uploaded.Add(int64(up))
			}()
		}

		wg.Wait()

		return int(uploaded.Load()), errs.Unwrap()
	*/
}

func (m *MultiTSDB) RemoveLockFilesIfAny() error {
	fis, err := os.ReadDir(m.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	errs := &multierr.Errors{}
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}

		if err := os.Remove(filepath.Join(m.defaultTenantDataDir(fi.Name()), "lock")); err != nil {
			if os.IsNotExist(err) {
				continue
			}

			errs.Add(err)
			continue
		}

		m.logger.Info("a leftover lockfile found and removed",
			zap.String("tenant", fi.Name()))
	}

	return errs.Unwrap()
}

func (m *MultiTSDB) getOrLoadTenant(tenantID manta.ID, blockingStart bool) (*Tenant, error) {
	// fast path, as creating tenants is a very rare operation
	m.mtx.RLock()
	tenant, exist := m.tenants[tenantID]
	m.mtx.RUnlock()
	if exist {
		return tenant, nil
	}

	// Slow path needs to lock fully and attempt to read again to prevent race conditions,
	// where since the fast path was tried, there may have actually
	// been the same tenant inserted in the map
	m.mtx.Lock()
	tenant, exist = m.tenants[tenantID]
	if exist {
		m.mtx.Unlock()
		return tenant, nil
	}

	tenant = newTenant()
	m.tenants[tenantID] = tenant
	m.mtx.Unlock()

	logger := m.logger.With(zap.String("tenant", tenantID.String()))
	if !blockingStart {
		go func() {
			err := m.startTSDB(logger, tenantID, tenant)
			if err != nil {
				logger.Warn("failed to start tsdb asynchronously",
					zap.Error(err))
			}
		}()

		return tenant, nil
	}

	return tenant, m.startTSDB(logger, tenantID, tenant)
}

func (m *MultiTSDB) startTSDB(zl *zap.Logger, tenantID manta.ID, tenant *Tenant) error {
	reg := prometheus.WrapRegistererWith(prometheus.Labels{"tenant": tenantID.String()}, m.reg)
	dataDir := m.defaultTenantDataDir(tenantID.String())
	opts := *m.tsdbOpts
	kitlog := log.NewZapToGokitLogAdapter(zl)

	s, err := tsdb.Open(dataDir, kitlog, &UnRegisterer{Registerer: reg}, &opts)
	if err != nil {
		m.mtx.Lock()
		delete(m.tenants, tenantID)
		m.mtx.Unlock()
		return err
	}

	/*	var ship *shipper.Shipper
		if m.bucket != nil {
			ship = shipper.New(
				kitlog,
				reg,
				dataDir,
				m.bucket,
				func() labels.Labels {
					return nil
				},
				metadata.ReceiveSource,
				false,
				m.allowOutOfOrderUpload,
			)
		}
		tenant.set(s, ship)
	*/

	tenant.set(s)

	zl.Info("TSDB is now ready")

	return nil
}

func emptyLabels() labels.Labels {
	return nil
}

func (m *MultiTSDB) defaultTenantDataDir(tenantID string) string {
	return path.Join(m.dataDir, tenantID)
}

func (m *MultiTSDB) TenantAppendable(id manta.ID) (storage.Appendable, error) {
	ts, err := m.getOrLoadTenant(id, false)
	if err != nil {
		return nil, err
	}

	return ts.readyStorage().Get(), nil
}

type Tenant struct {
	logger *zap.Logger

	readyS *ReadyStorage
	// ship   *shipper.Shipper

	mtx *sync.RWMutex
}

func newTenant() *Tenant {
	return &Tenant{
		readyS: &ReadyStorage{},
		mtx:    &sync.RWMutex{},
	}
}

func (t *Tenant) readyStorage() *ReadyStorage {
	return t.readyS
}

/*
func (t *Tenant) shipper() *shipper.Shipper {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	return t.ship
}

func (t *Tenant) set(tenantTSDB *tsdb.DB, ship *shipper.Shipper) {
	t.readyS.Set(tenantTSDB)

	t.mtx.Lock()
	t.ship = ship
	t.mtx.Unlock()
}*/

func (t *Tenant) set(db *tsdb.DB) {
	t.readyS.Set(db)
}

// ErrNotReady is returned if the underlying storage is not ready yet.
var ErrNotReady = errors.New("TSDB not ready")

// ReadyStorage implements the Storage interface while allowing to set the actual
// storage at a later point in time.
// TODO: Replace this with upstream Prometheus implementation when it is exposed.
type ReadyStorage struct {
	mtx sync.RWMutex
	a   *adapter
}

// Set the storage.
func (s *ReadyStorage) Set(db *tsdb.DB) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.a = &adapter{db: db}
}

// Get the storage.
func (s *ReadyStorage) Get() *tsdb.DB {
	if x := s.get(); x != nil {
		return x.db
	}
	return nil
}

func (s *ReadyStorage) get() *adapter {
	s.mtx.RLock()
	x := s.a
	s.mtx.RUnlock()
	return x
}

// StartTime implements the Storage interface.
func (s *ReadyStorage) StartTime() (int64, error) {
	return 0, errors.New("not implemented")
}

// Querier implements the Storage interface.
func (s *ReadyStorage) Querier(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
	if x := s.get(); x != nil {
		return x.Querier(ctx, mint, maxt)
	}
	return nil, ErrNotReady
}

// Appender implements the Storage interface.
func (s *ReadyStorage) Appender(ctx context.Context) (storage.Appender, error) {
	if x := s.get(); x != nil {
		return x.Appender(ctx)
	}
	return nil, ErrNotReady
}

// Close implements the Storage interface.
func (s *ReadyStorage) Close() error {
	if x := s.Get(); x != nil {
		return x.Close()
	}
	return nil
}

// adapter implements a storage.Storage around TSDB.
type adapter struct {
	db *tsdb.DB
}

// StartTime implements the Storage interface.
func (a adapter) StartTime() (int64, error) {
	return 0, errors.New("not implemented")
}

func (a adapter) Querier(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
	q, err := a.db.Querier(ctx, mint, maxt)
	if err != nil {
		return nil, err
	}
	return q, nil
}

// Appender returns a new appender against the storage.
func (a adapter) Appender(ctx context.Context) (storage.Appender, error) {
	return a.db.Appender(ctx), nil
}

// Close closes the storage and all its underlying resources.
func (a adapter) Close() error {
	return a.db.Close()
}

// UnRegisterer is a Prometheus registerer that
// ensures that collectors can be registered
// by unregistering already-registered collectors.
// FlushableStorage uses this registerer in order
// to not lose metric values between DB flushes.
type UnRegisterer struct {
	prometheus.Registerer
}

func (u *UnRegisterer) MustRegister(cs ...prometheus.Collector) {
	for _, c := range cs {
		if err := u.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
				if ok = u.Unregister(c); !ok {
					panic("unable to unregister existing collector")
				}
				u.Registerer.MustRegister(c)
				continue
			}
			panic(err)
		}
	}
}
