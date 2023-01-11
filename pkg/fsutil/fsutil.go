package fsutil

import (
	"io/ioutil"
	"os"

	"path/filepath"
)

const (
	// PrivateDirMode grants owner to make/remove files inside the directory.
	PrivateDirMode = 0700

	// PrivateFileMode grants owner to read/write a file.
	PrivateFileMode = 0600
)

// IsDirWriteable checks if dir is writable by writing and removing a file
// to dir. It returns nil if dir is writable.
func IsDirWriteable(dir string) error {
	f, err := filepath.Abs(filepath.Join(dir, ".touch"))
	if err != nil {
		return err
	}
	if err := os.WriteFile(f, []byte(""), PrivateFileMode); err != nil {
		return err
	}
	return os.Remove(f)
}

// TouchDirAll is similar to os.MkdirAll. it creates directories with
// 0700 permission if any directory does not exists. touchDirAll also ensures
// the given directory is wrtable.
func TouchDirAll(path string) error {
	if err := os.MkdirAll(path, PrivateDirMode); err != nil {
		return err
	}

	f, err := filepath.Abs(filepath.Join(path, ".touch"))
	if err != nil {
		return err
	}

	if err := os.WriteFile(f, []byte(""), PrivateDirMode); err != nil {
		return err
	}

	return os.Remove(f)
}
