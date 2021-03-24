package vfs

import (
	"io"
	"os"
)

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

type Opener interface {
	Open(name string) (ReadSeekCloser, error)
}

type FileSystem interface {
	Opener

	Lstat(path string) (os.FileInfo, error)
	Stat(path string) (os.FileInfo, error)
	ReadDir(path string) ([]os.FileInfo, error)
	String() string
}
