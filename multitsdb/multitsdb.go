package multitsdb

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/log"
	"github.com/f1shl3gs/manta/pkg/multierr"
)

// ErrNotReady is returned if the underlying storage is not ready yet.
var ErrNotReady = errors.New("TSDB not ready")

type MultiTSDB struct {
	dataDir  string
	logger   *zap.Logger
	reg      prometheus.Registerer
	tsdbOpts *tsdb.Options
	labels   labels.Labels

	mtx                   sync.RWMutex
	tenants               map[manta.ID]*tenant
	allowOutOfOrderUpload bool
}

func (m *MultiTSDB) Queryable(ctx context.Context, id manta.ID) (storage.Queryable, error) {
	t, err := m.getOrLoadTenant(id, true)
	if err != nil {
		return nil, err
	}

	return t.readyS.Get(), nil
}

func (m *MultiTSDB) Appendable(ctx context.Context, id manta.ID) (storage.Appendable, error) {
	t, err := m.getOrLoadTenant(id, true)
	if err != nil {
		return nil, err
	}

	return t.readyS.Get(), nil
}

func NewMultiTSDB(
	dataDir string,
	logger *zap.Logger,
	reg prometheus.Registerer,
	tsdbOpts *tsdb.Options,
	labels labels.Labels,
	allowOutOfOrderUpload bool,
) *MultiTSDB {
	return &MultiTSDB{
		logger:                logger.Named("multitsdb"),
		dataDir:               dataDir,
		reg:                   reg,
		tsdbOpts:              tsdbOpts,
		labels:                labels,
		allowOutOfOrderUpload: allowOutOfOrderUpload,
		tenants:               make(map[manta.ID]*tenant),
	}
}

func (m *MultiTSDB) Open() error {
	if err := os.MkdirAll(m.dataDir, 0750); err != nil {
		return err
	}

	files, err := os.ReadDir(m.dataDir)
	if err != nil {
		return err
	}

	var g errgroup.Group
	for _, f := range files {
		file := f
		if !file.IsDir() {
			continue
		}

		var id manta.ID
		if err = id.DecodeFromString(f.Name()); err != nil {
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

	errs := &multierr.List{}
	wg := &sync.WaitGroup{}

	for id, tenant := range m.tenants {
		db := tenant.readyS.Get()
		if db == nil {
			m.logger.Error("Flushing TSDB failed, not ready",
				zap.String("tenant", id.String()))
			continue
		}

		m.logger.Info("Flushing TSDB", zap.String("tenant", id.String()))

		wg.Add(1)
		go func() {
			defer wg.Done()

			head := db.Head()
			if err := db.CompactHead(tsdb.NewRangeHead(head, head.MinTime(), head.MaxTime()-1)); err != nil {
				errs.Append(err)
			}
		}()
	}

	wg.Wait()

	return errs.Err()
}

func (m *MultiTSDB) Close() error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	errs := &multierr.List{}
	for id, tenant := range m.tenants {
		db := tenant.readyS.Get()
		if db == nil {
			m.logger.Error("Closing TSDB failed, not ready",
				zap.String("tenant", id.String()))
			continue
		}

		if err := db.Close(); err != nil {
			errs.Append(err)
		}
	}

	return errs.Err()
}

func (m *MultiTSDB) RemoveLockFilesIfAny() error {
	fis, err := os.ReadDir(m.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	errs := &multierr.List{}
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}

		if err := os.Remove(filepath.Join(m.defaultTenantDataDir(fi.Name()), "lock")); err != nil {
			if os.IsNotExist(err) {
				continue
			}

			errs.Append(err)
			continue
		}

		m.logger.Info("a leftover lockfile found and removed", zap.String("tenant", fi.Name()))
	}

	return errs.Err()
}

func (m *MultiTSDB) defaultTenantDataDir(tenantID string) string {
	return path.Join(m.dataDir, tenantID)
}

func (m *MultiTSDB) getOrLoadTenant(id manta.ID, blockingStart bool) (*tenant, error) {
	// Fast path, as creating tenants is a very rare operation
	m.mtx.RLock()
	tenant, exist := m.tenants[id]
	m.mtx.RUnlock()
	if exist {
		return tenant, nil
	}

	// Slow path needs to lock fully and attempt to read again to prevent race conditions,
	// where since the fast path was tried, there may have actually been the same tenant
	// inserted in the map.
	m.mtx.Lock()
	tenant, exist = m.tenants[id]
	if exist {
		m.mtx.Unlock()
		return tenant, nil
	}

	tenant = newTenant()
	m.tenants[id] = tenant
	m.mtx.Unlock()

	logger := m.logger.With(zap.String("tenant", id.String()))
	if !blockingStart {
		go func() {
			if err := m.startTSDB(logger, id, tenant); err != nil {
				m.logger.Error("failed to start tsdb asynchronously",
					zap.Error(err))
			}
		}()

		return tenant, nil
	}

	return tenant, m.startTSDB(logger, id, tenant)
}

func (m *MultiTSDB) startTSDB(zl *zap.Logger, tenantID manta.ID, tenant *tenant) error {
	reg := prometheus.WrapRegistererWith(prometheus.Labels{"tenant": tenantID.String()}, m.reg)
	dataDir := m.defaultTenantDataDir(tenantID.String())
	opts := *m.tsdbOpts
	kitlog := log.NewZapToGokitLogAdapter(zl)

	db, err := tsdb.Open(dataDir, kitlog, &UnRegisterer{reg}, &opts, nil)
	if err != nil {
		m.mtx.Lock()
		delete(m.tenants, tenantID)
		m.mtx.Unlock()

		return err
	}

	tenant.readyS.Set(db)

	zl.Info("TSDB is now ready")

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

type tenant struct {
	mtx *sync.RWMutex

	readyS *ReadyStorage
}

func newTenant() *tenant {
	return &tenant{
		readyS: &ReadyStorage{},
		mtx:    &sync.RWMutex{},
	}
}
