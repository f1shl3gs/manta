package wal

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/f1shl3gs/manta/pkg/mmap"

	"github.com/pkg/errors"
	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
)

// WAL is divided up into entryFiles. Each entry file stores maxNumEntries in
// the first logFileOffset bytes. Each entry takes a fixed entrySize bytes of
// space. The variable length data for these entries is written after
// logFileOffset from file beginning. Once snapshot is taken, all the files
// containing entries below snapshot index are deleted.

const (
	// maxNumEntries is maximum number of entries before rotating the file.
	maxNumEntries = 30000
	// logFileOffset is offset in the log file where data is stored.
	logFileOffset = 1 << 20 // 1MB

	// logFileSize is the initial size of the log file.
	logFileSize = 256 << 20 // 256MB
	// entrySize is the size in bytes of a single entry.
	entrySize = 32
	// logSuffix is the suffix for log files.
	logSuffix = ".wal"
)

var (
	emptyEntry = entry(make([]byte, entrySize))
)

type entry []byte

func (e entry) Term() uint64       { return binary.BigEndian.Uint64(e) }
func (e entry) Index() uint64      { return binary.BigEndian.Uint64(e[8:]) }
func (e entry) DataOffset() uint64 { return binary.BigEndian.Uint64(e[16:]) }
func (e entry) Type() uint64       { return binary.BigEndian.Uint64(e[24:]) }

func marshalEntry(b []byte, term, index, do, typ uint64) {
	if len(b) != entrySize {
		panic("entry buffer must be 32")
	}

	binary.BigEndian.PutUint64(b, term)
	binary.BigEndian.PutUint64(b[8:], index)
	binary.BigEndian.PutUint64(b[16:], do)
	binary.BigEndian.PutUint64(b[24:], typ)
}

// logFile represents a single log file.
type logFile struct {
	*mmap.File

	fid    int64
	logger *zap.Logger
}

func logFname(dir string, id int64) string {
	return filepath.Join(dir, fmt.Sprintf("%05d%s", id, logSuffix))
}

// openLogFile opens a logFile in the given directory. The filename is
// constructed based on the value of fid.
func openLogFile(dir string, fid int64, logger *zap.Logger) (*logFile, error) {
	fpath := logFname(dir, fid)
	lf := &logFile{
		fid:    fid,
		logger: logger.With(zap.String("dir", dir), zap.Int64("fid", fid)),
	}
	var err error
	// Open the file in read-write mode and create it if it doesn't exist yet.
	lf.File, err = mmap.OpenMmapFile(fpath, os.O_RDWR|os.O_CREATE, logFileSize)
	if err == mmap.ErrNewFile {
		// New file created
		mmap.ZeroOut(lf.Data, 0, logFileOffset)
	} else if err != nil {
		return nil, err
	}

	return lf, nil
}

// getEntry gets the entry at the slot idx.
func (lf *logFile) getEntry(idx int) entry {
	if lf == nil {
		return emptyEntry
	}

	if idx >= maxNumEntries {
		panic("entry index must less than maxNumEntries(30000)")
	}

	offset := idx * entrySize
	return lf.Data[offset : offset+entrySize]
}

// GetRaftEntry gets the entry at the index idx, reads the data from the appropriate
// offset and converts it to a raftpb.Entry object.
func (lf *logFile) GetRaftEntry(idx int) raftpb.Entry {
	entry := lf.getEntry(idx)
	re := raftpb.Entry{
		Term:  entry.Term(),
		Index: entry.Index(),
		Type:  raftpb.EntryType(int32(entry.Type())),
	}

	if entry.DataOffset() > 0 {
		if entry.DataOffset() >= uint64(len(lf.Data)) {
			panic("entry offset must less than mmap file length")
		}

		data := lf.Slice(int(entry.DataOffset()))
		if len(data) > 0 {
			// Copy the data over to allow the mmaped file to be deleted later.
			re.Data = append(re.Data, data...)
		}
	}

	return re
}

// firstIndex returns the first index in the file.
func (lf *logFile) firstIndex() uint64 {
	return lf.getEntry(0).Index()
}

// firstEmptySlot returns the index of the first empty slot in the file.
func (lf *logFile) firstEmptySlot() int {
	return sort.Search(maxNumEntries, func(i int) bool {
		e := lf.getEntry(i)
		return e.Index() == 0
	})
}

// lastEntry returns the last valid entry in the file.
func (lf *logFile) lastEntry() entry {
	// This would return the first pos, where e.Index() == 0.
	pos := lf.firstEmptySlot()
	if pos > 0 {
		pos--
	}
	return lf.getEntry(pos)
}

// slotGe would return -1 if raftIndex < firstIndex in this file.
// Would return maxNumEntries if raftIndex > lastIndex in this file.
// If raftIndex is found, or the entryFile has empty slots, the offset would be between
// [0, maxNumEntries).
func (lf *logFile) slotGe(raftIndex uint64) int {
	fi := lf.firstIndex()
	// If first index is zero or the first index is less than raftIndex, this
	// raftindex should be in a previous file.
	if fi == 0 || raftIndex < fi {
		return -1
	}

	// Look at the entry at slot diff. If the log has entries for all indices between
	// fi and raftIndex without any gaps, the entry should be there. This is an
	// optimization to avoid having to perform the search below.
	if diff := int(raftIndex - fi); diff < maxNumEntries && diff >= 0 {
		e := lf.getEntry(diff)
		if e.Index() == raftIndex {
			return diff
		}
	}

	// Find the first entry which has in index >= to raftIndex.
	return sort.Search(maxNumEntries, func(i int) bool {
		e := lf.getEntry(i)
		if e.Index() == 0 {
			// We reached too far to the right and found an empty slot.
			return true
		}
		return e.Index() >= raftIndex
	})
}

// delete unmaps and deletes the file.
func (lf *logFile) delete() error {
	lf.logger.Info("deleting file", zap.String("filename", lf.Fd.Name()))

	err := lf.Delete()
	if err != nil {
		lf.logger.Info("deleting file failed", zap.Error(err), zap.String("filename", lf.Fd.Name()))
	}
	return err
}

// getLogFiles returns all the log files in the directory sorted by the first
// index in each file.
func getLogFiles(dir string, logger *zap.Logger) ([]*logFile, error) {
	entryFiles, err := func() (list []string, err error) {
		err = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
			if fi.IsDir() {
				return nil
			}

			if strings.HasSuffix(path, logSuffix) {
				list = append(list, path)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		return list, nil
	}()
	if err != nil {
		return nil, err
	}

	var files []*logFile
	seen := make(map[int64]struct{})

	for _, fpath := range entryFiles {
		_, fname := filepath.Split(fpath)
		fname = strings.TrimSuffix(fname, logSuffix)

		fid, err := strconv.ParseInt(fname, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "while parsing: %s", fpath)
		}

		if _, ok := seen[fid]; ok {
			logger.Fatal("repeated entry file id found", zap.Int64("fid", fid))
		}
		seen[fid] = struct{}{}

		f, err := openLogFile(dir, fid, logger)
		if err != nil {
			return nil, err
		}

		logger.Info("found wal file", zap.Int64("fid", fid), zap.Uint64("first-index", f.firstIndex()))
		files = append(files, f)
	}

	// Sort files by the first index they store.
	sort.Slice(files, func(i, j int) bool {
		return files[i].getEntry(0).Index() < files[j].getEntry(0).Index()
	})
	return files, nil
}
