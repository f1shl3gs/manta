package raftstore

import (
	"encoding/binary"
	"io"

	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.etcd.io/etcd/server/v3/wal"
	"go.etcd.io/etcd/server/v3/wal/walpb"
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
	*wal.WAL
	*snap.Snapshotter
}

func NewStorage(w *wal.WAL, s *snap.Snapshotter) Storage {
	return &storage{w, s}
}

// SaveSnap saves the snapshot file to disk and writes the WAL snapshot entry.
func (st *storage) SaveSnap(snap raftpb.Snapshot) error {
	walsnap := walpb.Snapshot{
		Index: snap.Metadata.Index,
		Term:  snap.Metadata.Term,
		// ConfState: &snap.Metadata.ConfState,
	}
	// save the snapshot file before writing the snapshot to the wal.
	// This makes it possible for the snapshot file to become orphaned, but prevents
	// a WAL snapshot entry from having no corresponding snapshot file.
	err := st.Snapshotter.SaveSnap(snap)
	if err != nil {
		return err
	}
	// gofail: var raftBeforeWALSaveSnaphot struct{}

	return st.WAL.SaveSnapshot(walsnap)
}

// Release releases resources older than the given snap and are no longer needed:
// - releases the locks to the wal files that are older than the provided wal for the given snap.
// - deletes any .snap.db files that are older than the given snap.
func (st *storage) Release(snap raftpb.Snapshot) error {
	if err := st.WAL.ReleaseLockTo(snap.Metadata.Index); err != nil {
		return err
	}
	return st.Snapshotter.ReleaseSnapDBs(snap)
}

// readWAL reads the WAL at the given snap and returns the wal,
// it's latest HardState and cluster ID, and all entries that
// appear after the position of the given snap in the WAL. The
// snap must have been previously saved to the WAL, or this call
// will panic
func readWAL(
	logger *zap.Logger,
	waldir string,
	snapshot walpb.Snapshot,
) (w *wal.WAL, id, cid uint64, st raftpb.HardState, ents []raftpb.Entry) {
	var (
		err      error
		metadata []byte
	)

	repaired := false
	for {
		if w, err = wal.Open(logger, waldir, snapshot); err != nil {
			logger.Fatal("failed to open WAL",
				zap.String("dir", waldir),
				zap.Error(err))
		}

		metadata, st, ents, err = w.ReadAll()
		if err == nil {
			break
		}

		// trying to repair WAL
		w.Close()
		// we can only repair ErrUnexpectedEOF and we never repair twice
		if repaired || err != io.ErrUnexpectedEOF {
			logger.Fatal("failed to read WAL, cannot be repaired",
				zap.Error(err))
		}

		if !wal.Repair(logger, waldir) {
			logger.Fatal("failed to repair WAL",
				zap.Error(err))
		} else {
			logger.Info("repaired WAL",
				zap.Error(err))
			repaired = true
		}
	}

	id = binary.BigEndian.Uint64(metadata)
	cid = binary.BigEndian.Uint64(metadata[8:])

	return w, id, cid, st, ents
}
