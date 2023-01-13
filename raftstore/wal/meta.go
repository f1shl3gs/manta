package wal

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/f1shl3gs/manta/pkg/mmap"

	"github.com/pkg/errors"
	"go.etcd.io/raft/v3"
	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
)

type MetaInfo int

const (
	RaftId MetaInfo = iota
	ClusterId
	CheckpointIndex
	SnapshotIndex
	SnapshotTerm
)

// getOffset returns offsets in wal.meta file.
func getOffset(info MetaInfo) int {
	switch info {
	case RaftId:
		return 0
	case ClusterId:
		return 8
	case CheckpointIndex:
		return 16
	case SnapshotIndex:
		return snapshotIndex
	case SnapshotTerm:
		return snapshotIndex + 8
	default:
		panic("Invalid info: " + fmt.Sprint(info))
	}
}

// Constants to use when writing to mmap'ed meta and entry files.
const (
	// metaName is the name of the file used to store metadata (e.g raft ID, checkpoint).
	metaName = "wal.meta"
	// metaFileSize is the size of the wal.meta file.
	metaFileSize = 1 << 20
	//hardStateOffset is the offset of the hard sate within the wal.meta file.
	hardStateOffset = 512
	// snapshotIndex stores the index and term corresponding to the snapshot.
	snapshotIndex = 1024
	// snapshotOffest is the offset of the snapshot within the wal.meta file.
	snapshotOffset = snapshotIndex + 16
)

func readSlice(dst []byte, offset int) []byte {
	b := dst[offset:]
	sz := binary.BigEndian.Uint32(b)
	return b[4 : 4+sz]
}
func writeSlice(dst []byte, src []byte) {
	binary.BigEndian.PutUint32(dst[:4], uint32(len(src)))
	copy(dst[4:], src)
}
func allocateSlice(dst []byte, sz int) []byte {
	binary.BigEndian.PutUint32(dst[:4], uint32(sz))
	return dst[4 : 4+sz]
}
func sliceSize(dst []byte, offset int) int {
	sz := binary.BigEndian.Uint32(dst[offset:])
	return 4 + int(sz)
}

// metaFile stores the RAFT metadata (e.g RAFT ID, snapshot, hard state).
type metaFile struct {
	*mmap.MmapFile

	logger *zap.Logger
}

// newMetaFile opens the meta file in the given directory.
func newMetaFile(dir string, logger *zap.Logger) (*metaFile, error) {
	fname := filepath.Join(dir, metaName)
	// Open the file in read-write mode and creates it if it doesn't exist.
	mf, err := mmap.OpenMmapFile(fname, os.O_RDWR|os.O_CREATE, metaFileSize)
	if err == mmap.ErrNewFile {
		mmap.ZeroOut(mf.Data, 0, snapshotOffset+4)
	} else if err != nil {
		return nil, errors.Wrapf(err, "unable to open meta file")
	}
	return &metaFile{MmapFile: mf, logger: logger}, nil
}

func (m *metaFile) bufAt(info MetaInfo) []byte {
	pos := getOffset(info)
	return m.Data[pos : pos+8]
}
func (m *metaFile) Uint(info MetaInfo) uint64 {
	return binary.BigEndian.Uint64(m.bufAt(info))
}
func (m *metaFile) SetUint(info MetaInfo, id uint64) {
	binary.BigEndian.PutUint64(m.bufAt(info), id)
}

func (m *metaFile) StoreHardState(hs *raftpb.HardState) error {
	if hs == nil || raft.IsEmptyHardState(*hs) {
		return nil
	}
	buf, err := hs.Marshal()
	if err != nil {
		return errors.Wrapf(err, "cannot marshal hard state")
	}
	if !(len(buf) < snapshotIndex-hardStateOffset) {
		m.logger.Fatal("buf length should less than snapshotIndex - hardStateOffset",
			zap.Int("buf", len(buf)),
			zap.Int("remain", snapshotIndex-hardStateOffset))
	}

	writeSlice(m.Data[hardStateOffset:], buf)
	return nil
}

func (m *metaFile) HardState() (raftpb.HardState, error) {
	val := readSlice(m.Data, hardStateOffset)
	var hs raftpb.HardState

	if len(val) == 0 {
		return hs, nil
	}
	if err := hs.Unmarshal(val); err != nil {
		return hs, errors.Wrapf(err, "cannot parse hardState")
	}
	return hs, nil
}

// StoreSnapshot stores the snapshot and return the snapshot length and error.
func (m *metaFile) StoreSnapshot(snap *raftpb.Snapshot) error {
	if snap == nil || raft.IsEmptySnap(*snap) {
		return nil
	}
	m.SetUint(SnapshotIndex, snap.Metadata.Index)
	m.SetUint(SnapshotTerm, snap.Metadata.Term)

	buf, err := snap.Marshal()
	if err != nil {
		return errors.Wrapf(err, "cannot marshal snapshot")
	}
	m.logger.Info("got valid snapshot to store",
		zap.Int("length", len(buf)))

	for len(m.Data)-snapshotOffset < len(buf) {
		if err := m.Truncate(2 * int64(len(m.Data))); err != nil {
			return errors.Wrapf(err, "while truncating: %s", m.Fd.Name())
		}
	}
	writeSlice(m.Data[snapshotOffset:], buf)

	return nil
}

func (m *metaFile) snapshot() (raftpb.Snapshot, error) {
	val := readSlice(m.Data, snapshotOffset)

	var snap raftpb.Snapshot
	if len(val) == 0 {
		return snap, nil
	}

	if err := snap.Unmarshal(val); err != nil {
		return snap, errors.Wrapf(err, "cannot parse snapshot")
	}
	return snap, nil
}
