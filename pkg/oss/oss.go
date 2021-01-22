package oss

import (
	"golang.org/x/tools/godoc/vfs"
	"os"
)

type oss struct {
}

func (o oss) Open(name string) (vfs.ReadSeekCloser, error) {
	panic("implement me")
}

func (o oss) Lstat(path string) (os.FileInfo, error) {
	panic("implement me")
}

func (o oss) Stat(path string) (os.FileInfo, error) {
	panic("implement me")
}

func (o oss) ReadDir(path string) ([]os.FileInfo, error) {
	panic("implement me")
}

func (o oss) RootType(path string) vfs.RootType {
	panic("implement me")
}

func (o oss) String() string {
	panic("implement me")
}
