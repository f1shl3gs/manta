package profiler

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultMaxSize = uint64(128 * 1024 * 1024)
)

// dumpStore represents the directory where heap profiles are stored.
// It supports automatic garbage collection of old profiles
type dumpStore struct {
	maxSize        uint64
	dir            string
	prefix, suffix string

	cleanupFn func(string) error
}

func (s *dumpStore) format(ts time.Time, current uint64) string {
	filename := fmt.Sprintf("%s.%s.%d%s",
		s.prefix, ts.Format(timestampFormat), current, s.suffix)
	return filepath.Join(s.dir, filename)
}

// GC runs the GC policy on this dumpStore
func (s *dumpStore) GC() error {
	// NB: ioutil.ReadDir sorts the file names in ascending order
	// This brings the oldest files first.
	files, err := ioutil.ReadDir(s.dir)
	if err != nil {
		return err
	}

	maxSize := s.maxSize
	if maxSize <= 0 {
		maxSize = defaultMaxSize
	}

	actualSize := uint64(0)
	for i := len(files) - 1; i >= 0; i-- {
		fi := files[i]
		if !s.checkOwnsFile(fi.Name()) {
			// Not for this dumpStore. Ignore it silently
			continue
		}

		// Note: we are counting preserved files against the maximum
		actualSize += uint64(fi.Size())

		// Does the entry make the directory grow past the maximum?
		if actualSize < maxSize {
			continue
		}

		if err := s.cleanupFn(fi.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (s *dumpStore) checkOwnsFile(name string) bool {
	return strings.HasPrefix(name, s.prefix) && strings.HasSuffix(name, s.suffix)
}
