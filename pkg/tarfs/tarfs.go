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

var (
	ErrInvalidCharacterInPath = errors.New("invalid character in path")
)

// TarFS implement fs.FS, so it can be wrapped as http.FileSystem
type TarFS struct {
	files map[string]file
}

// Open implement fs.FS
func (tarfs *TarFS) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return nil, ErrInvalidCharacterInPath
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

		var data []byte
		if hdr.Size != 0 {
			data, err = readAll(tr, hdr)
			_, err = tr.Read(data)
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
			}
		}

		files[path.Join("/", hdr.Name)] = file{
			data: data,
			fi:   hdr.FileInfo(),
		}
	}

	return &TarFS{files: files}, nil
}

// readAll reads from r until an error or EOF and returns the data it read.
// It only allocate 1 time per call, so there is no slice growth overhead.
func readAll(r io.Reader, hdr *tar.Header) ([]byte, error) {
	buf := make([]byte, 0, hdr.Size)

	for {
		n, err := r.Read(buf[len(buf):cap(buf)])
		buf = buf[:len(buf)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}

			return buf, err
		}
	}
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
