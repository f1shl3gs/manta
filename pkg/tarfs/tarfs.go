package tarfs

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// TarFS implement fs.FS
type TarFS struct {
	files map[string]file
}

func (tarfs *TarFS) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return nil, errors.New("http: invalid character in file path")
	}

	f, exists := tarfs.files[path.Join("/", name)]
	if !exists {
		return nil, os.ErrNotExist
	}

	f.buf = bytes.NewBuffer(f.data)

	return &f, nil
}

func New(r io.Reader) (*TarFS, error) {
	tr := tar.NewReader(r)
	files := map[string]file{}

	// Iterate through the files in the archive
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}

		files[path.Join("/", hdr.Name)] = file{
			data: data,
			fi:   hdr.FileInfo(),
		}
	}

	return &TarFS{files: files}, nil
}

func (tarfs *TarFS) Walk(fn func(name string, fi os.FileInfo)) {
	for name, file := range tarfs.files {
		fn(name, file.fi)
	}
}

// file implement fs.File
type file struct {
	buf   *bytes.Buffer
	data  []byte
	fi    os.FileInfo
	files []os.FileInfo
}

func (f *file) Stat() (fs.FileInfo, error) {
	return f.fi, nil
}

func (f *file) Read(bytes []byte) (int, error) {
	return f.buf.Read(bytes)
}

func (f *file) Close() error {
	return nil
}
