package wal

import (
	"fmt"
	"math"
	"sync"

	"github.com/pkg/errors"
	"go.etcd.io/raft/v3"
	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
)

// versionKey is hardcoded into the special key used to fetch the maximum version from the DB.
const versionKey = 1

// DiskStorage handles disk access and writing for the RAFT write-ahead log.
// Dir contains wal.meta file and <start idx zero padded>.wal files.
//
// === wal.meta file ===
// This file is generally around 4KB, so it can fit nicely in one Linux page.
//
//	Layout:
//
// 00-08 Bytes: Raft ID
// 08-16 Bytes: Group ID
// 16-24 Bytes: Checkpoint Index
// 512 Bytes: Hard State (Marshalled)
// 1024-1032 Bytes: Snapshot Index
// 1032-1040 Bytes: Snapshot Term
// 1040 Bytes: Snapshot (Marshalled)
//
// --- <0000i>.wal files ---
// These files contain raftpb.Entry protos. Each entry is composed of term, index, type and data.
//
// Term takes 8 bytes. Index takes 8 bytes. Type takes 8 bytes. And for data, we store an offset to
// the actual slice, which is 8 bytes. Size of entry = 32 bytes.
// First 30K entries would consume 960KB, hence fitting on the first MB of the file (logFileOffset).
//
// Pre-allocate 1MB in each file just for these entries, and zero them out explicitly. Zeroing them
// out ensures that we know when these entries end, in case of a restart.
//
// And the data for these entries are laid out starting logFileOffset. Those are the offsets you
// store in the Entry for Data field.
// After 30K entries, we rotate the file.
//
// --- clean up ---
// If snapshot idx = Idx_s. We find the first log file whose first entry is
// less than Idx_s. This file and anything above MUST be kept. All the log
// files lower than this file can be deleted.
//
// --- sync ---
// mmap fares well with process crashes without doing anything. In case
// HardSync is set, msync is called after every write, which flushes those
// writes to disk.
type DiskStorage struct {
	dir    string
	logger *zap.Logger

	meta *metaFile
	wal  *wal
	lock sync.Mutex
}

// Init initializes an instance of DiskStorage
func Init(dir string, logger *zap.Logger) (*DiskStorage, error) {
	s := &DiskStorage{
		dir:    dir,
		logger: logger,
	}

	var err error
	if s.meta, err = newMetaFile(dir, logger); err != nil {
		return nil, err
	}
	// fmt.Printf("meta: %s\n", hex.Dump(s.meta.data[1024:2048]))
	// fmt.Printf("found snapshot of size: %d\n", sliceSize(s.meta.data, snapshotOffset))

	if s.wal, err = openWal(dir, s.logger); err != nil {
		return nil, err
	}

	snap, err := s.meta.snapshot()
	if err != nil {
		return nil, err
	}

	first, _ := s.FirstIndex()
	if !raft.IsEmptySnap(snap) {
		if snap.Metadata.Index+1 != first {
			panic(fmt.Sprintf("snap index: %d + 1 should be equal to first: %d\n", snap.Metadata.Index, first))
		}
	}

	// If db is not closed properly, there might be index ranges for which delete entries are not
	// inserted. So insert delete entries for those ranges starting from 0 to (first-1).
	s.wal.deleteBefore(first - 1)
	last := s.wal.LastIndex()

	s.logger.Info("init raft storage",
		zap.Uint64("snap-term", snap.Metadata.Term),
		zap.Uint64("snap-index", snap.Metadata.Index),
		zap.Uint64("first", first),
		zap.Uint64("last", last))

	return s, nil
}

func (w *DiskStorage) SetUint(info MetaInfo, id uint64) { w.meta.SetUint(info, id) }
func (w *DiskStorage) Uint(info MetaInfo) uint64        { return w.meta.Uint(info) }

func (w *DiskStorage) NodeID() uint64 {
	return w.meta.Uint(RaftId)
}

func (w *DiskStorage) SetNodeID(id uint64) {
	w.meta.SetUint(RaftId, id)
}

// reset resets the entries. Used for testing.
func (w *DiskStorage) reset(es []raftpb.Entry) error {
	// Clean out the state.
	if err := w.wal.reset(); err != nil {
		return err
	}
	return w.addEntries(es)
}

func (w *DiskStorage) HardState() (raftpb.HardState, error) {
	if w.meta == nil {
		return raftpb.HardState{}, errors.Errorf("uninitialized meta file")
	}
	return w.meta.HardState()
}

// Checkpoint returns the Raft index corresponding to the checkpoint.
func (w *DiskStorage) Checkpoint() (uint64, error) {
	if w.meta == nil {
		return 0, errors.Errorf("uninitialized meta file")
	}
	return w.meta.Uint(CheckpointIndex), nil
}

// Implement the Raft.Storage interface.
// -------------------------------------

// InitialState returns the saved HardState and ConfState information.
func (w *DiskStorage) InitialState() (hs raftpb.HardState, cs raftpb.ConfState, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	hs, err = w.meta.HardState()
	if err != nil {
		return
	}
	var snap raftpb.Snapshot
	snap, err = w.meta.snapshot()
	if err != nil {
		return
	}
	return hs, snap.Metadata.ConfState, nil
}

func (w *DiskStorage) NumEntries() int {
	w.lock.Lock()
	defer w.lock.Unlock()

	start := w.wal.firstIndex()

	var count int
	for {
		ents := w.wal.allEntries(start, math.MaxUint64, 64<<20)
		if len(ents) == 0 {
			return count
		}
		count += len(ents)
		start = ents[len(ents)-1].Index + 1
	}
}

// Entries returns a slice of log entries in the range [lo,hi).
// MaxSize limits the total size of the log entries returned, but
// Entries returns at least one entry if any.
func (w *DiskStorage) Entries(lo, hi, maxSize uint64) (es []raftpb.Entry, rerr error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// glog.Infof("Entries after lock: [%d, %d) maxSize:%d", lo, hi, maxSize)

	first := w.wal.firstIndex()
	if lo < first {
		w.logger.Error("lo should large or equal to first",
			zap.Uint64("lo", lo),
			zap.Uint64("first", first))

		return nil, raft.ErrCompacted
	}

	last := w.wal.LastIndex()
	if hi > last+1 {
		w.logger.Error("hi must less or equal to last+1",
			zap.Uint64("hi", hi),
			zap.Uint64("last", last))

		return nil, raft.ErrUnavailable
	}

	ents := w.wal.allEntries(lo, hi, maxSize)
	// glog.Infof("got entries [%d, %d): %+v\n", lo, hi, ents)
	return ents, nil
}

func (w *DiskStorage) Term(idx uint64) (uint64, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	si := w.meta.Uint(SnapshotIndex)
	if idx < si {
		w.logger.Error("term index < snap index",
			zap.Uint64("term", idx),
			zap.Uint64("snap", si))
		return 0, raft.ErrCompacted
	}
	if idx == si {
		return w.meta.Uint(SnapshotTerm), nil
	}

	term, err := w.wal.Term(idx)
	if err != nil {
		w.logger.Error("get term failed", zap.Error(err), zap.Uint64("idx", idx))
	}
	// glog.Errorf("Got term: %d for index: %d\n", term, idx)
	return term, err
}

func (w *DiskStorage) LastIndex() (uint64, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	li := w.wal.LastIndex()
	si := w.meta.Uint(SnapshotIndex)
	if li < si {
		return si, nil
	}
	return li, nil
}

func (w *DiskStorage) firstIndex() uint64 {
	if si := w.Uint(SnapshotIndex); si > 0 {
		return si + 1
	}
	return w.wal.firstIndex()
}

// FirstIndex returns the first index. It is typically SnapshotIndex+1.
func (w *DiskStorage) FirstIndex() (uint64, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.firstIndex(), nil
}

// Snapshot returns the most recent snapshot.  If snapshot is temporarily
// unavailable, it should return ErrSnapshotTemporarilyUnavailable, so raft
// state machine could know that Storage needs some time to prepare snapshot
// and call Snapshot later.
func (w *DiskStorage) Snapshot() (raftpb.Snapshot, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.meta.snapshot()
}

// ---------------- Raft.Storage interface complete.

// CreateSnapshot generates a snapshot with the given ConfState and data and writes it to disk.
func (w *DiskStorage) CreateSnapshot(i uint64, cs *raftpb.ConfState, data []byte) error {
	w.logger.Info("create snapshot",
		zap.Uint64("index", i),
		zap.String("cs", cs.String()))

	w.lock.Lock()
	defer w.lock.Unlock()

	first := w.firstIndex()
	if i < first {
		w.logger.Error("snapshot outofdate, index < first",
			zap.Uint64("index", i),
			zap.Uint64("first", first))
		return raft.ErrSnapOutOfDate
	}

	e, err := w.wal.seekEntry(i)
	if err != nil {
		return err
	}

	var snap raftpb.Snapshot
	snap.Metadata.Index = i
	snap.Metadata.Term = e.Term()
	if cs == nil {
		panic("confState cannot be nil")
	}
	snap.Metadata.ConfState = *cs
	snap.Data = data

	if err := w.meta.StoreSnapshot(&snap); err != nil {
		return err
	}
	// Now we delete all the files which are below the snapshot index.
	w.wal.deleteBefore(snap.Metadata.Index)
	return nil
}

// Save would write Entries, HardState and Snapshot to persistent storage in order, i.e. Entries
// first, then HardState and Snapshot if they are not empty. If persistent storage supports atomic
// writes then all of them can be written together. Note that when writing an Entry with Index i,
// any previously-persisted entries with Index >= i must be discarded.
func (w *DiskStorage) Save(h *raftpb.HardState, es []raftpb.Entry, snap *raftpb.Snapshot) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if err := w.wal.AddEntries(es); err != nil {
		return err
	}
	if err := w.meta.StoreHardState(h); err != nil {
		return err
	}
	if err := w.meta.StoreSnapshot(snap); err != nil {
		return err
	}
	return nil
}

// Append the new entries to storage.
func (w *DiskStorage) addEntries(entries []raftpb.Entry) error {
	if len(entries) == 0 {
		return nil
	}

	first, err := w.FirstIndex()
	if err != nil {
		return err
	}
	firste := entries[0].Index
	if firste+uint64(len(entries))-1 < first {
		// All of these entries have already been compacted.
		return nil
	}
	if first > firste {
		// Truncate compacted entries
		entries = entries[first-firste:]
	}

	// AddEntries would zero out all the entries starting entries[0].Index before writing.
	if err := w.wal.AddEntries(entries); err != nil {
		return errors.Wrapf(err, "while adding entries")
	}
	return nil
}

// TruncateEntriesUntil deletes the data field of every raft entry
// of type EntryNormal and index ∈ [0, lastIdx).
func (w *DiskStorage) TruncateEntriesUntil(lastIdx uint64) {
	w.wal.truncateEntriesUntil(lastIdx)
}

func (w *DiskStorage) NumLogFiles() int {
	return len(w.wal.files)
}

// Sync calls the Sync method in the underlying badger instance to write all the contents to disk.
func (w *DiskStorage) Sync() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if err := w.meta.Sync(); err != nil {
		return errors.Wrapf(err, "while syncing meta")
	}
	if err := w.wal.current.Sync(); err != nil {
		return errors.Wrapf(err, "while syncing current file")
	}
	return nil
}

// Close closes the DiskStorage.
func (w *DiskStorage) Close() error {
	return w.Sync()
}
