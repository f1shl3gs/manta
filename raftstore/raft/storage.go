package raft

import (
	"encoding/binary"
	"io"
	"sync"

	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
)

type Storage interface {
	// Save function saves ents and state to the underlying stable storage.
	// Save MUST block until st and ents are on stable storage.
	Save(st raftpb.HardState, ents []raftpb.Entry) error
	// SaveSnap function saves snapshot to the underlying stable storage.
	SaveSnap(snap raftpb.Snapshot) error
	// Close closes the Storage and performs finalization.
	Close() error
	// Release releases the locked wal files older than the provided snapshot.
	Release(snap raftpb.Snapshot) error
	// Sync WAL
	Sync() error
}

type storage struct {
	lg *zap.Logger
	s  *snap.Snapshotter

	// Mutex protected variables
	mux sync.RWMutex
	w   *wal.WAL
}

func NewStorage(lg *zap.Logger, w *wal.WAL, s *snap.Snapshotter) Storage {
	return &storage{lg: lg, w: w, s: s}
}

// SaveSnap saves the snapshot file to disk and writes the WAL snapshot entry.
func (st *storage) SaveSnap(snap raftpb.Snapshot) error {
	st.mux.RLock()
	defer st.mux.RUnlock()
	walsnap := walpb.Snapshot{
		Index:     snap.Metadata.Index,
		Term:      snap.Metadata.Term,
		ConfState: &snap.Metadata.ConfState,
	}
	// save the snapshot file before writing the snapshot to the wal.
	// This makes it possible for the snapshot file to become orphaned, but prevents
	// a WAL snapshot entry from having no corresponding snapshot file.
	err := st.s.SaveSnap(snap)
	if err != nil {
		return err
	}
	// gofail: var raftBeforeWALSaveSnaphot struct{}

	return st.w.SaveSnapshot(walsnap)
}

// Release releases resources older than the given snap and are no longer needed:
// - releases the locks to the wal files that are older than the provided wal for the given snap.
// - deletes any .snap.db files that are older than the given snap.
func (st *storage) Release(snap raftpb.Snapshot) error {
	st.mux.RLock()
	defer st.mux.RUnlock()
	if err := st.w.ReleaseLockTo(snap.Metadata.Index); err != nil {
		return err
	}
	return st.s.ReleaseSnapDBs(snap)
}

func (st *storage) Save(s raftpb.HardState, ents []raftpb.Entry) error {
	st.mux.RLock()
	defer st.mux.RUnlock()
	return st.w.Save(s, ents)
}

func (st *storage) Close() error {
	st.mux.Lock()
	defer st.mux.Unlock()
	return st.w.Close()
}

func (st *storage) Sync() error {
	st.mux.RLock()
	defer st.mux.RUnlock()
	return st.w.Sync()
}

// readWAL reads the WAL at the given snap and returns the wal, its latest HardState
// and cluster ID, and all entries that appear after the position of the given snap
// in the WAL.
// The snap must have been previously saved to the WAL, or this call will panic.
func readWAL(logger *zap.Logger, walDir string, snap walpb.Snapshot) (
	w *wal.WAL,
	id, cid uint64,
	st raftpb.HardState,
	ents []raftpb.Entry,
) {
	var (
		metadata []byte
		err      error
	)

	repaired := false
	for {
		if w, err = wal.Open(logger, walDir, snap); err != nil {
			logger.Fatal("failed to open WAL", zap.Error(err))
		}

		if metadata, st, ents, err = w.ReadAll(); err != nil {
			w.Close()

			// we can only repair ErrUnexpectedEOF and we never repair twice.
			if repaired || err != io.ErrUnexpectedEOF {
				logger.Fatal("failed to read WAL, cannot be repaired", zap.Error(err))
			}

			if !wal.Repair(logger, walDir) {
				logger.Fatal("failed to repair WAL", zap.Error(err))
			} else {
				logger.Info("repaired WAL", zap.Error(err))
				repaired = true
			}

			continue
		}

		break
	}

	id = binary.BigEndian.Uint64(metadata)
	cid = binary.BigEndian.Uint64(metadata[8:])

	return w, id, cid, st, ents
}
