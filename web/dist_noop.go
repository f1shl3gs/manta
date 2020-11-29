package web

import (
	"net/http"
	"os"
)

type noopFileSystem struct {
}

func (n *noopFileSystem) Open(name string) (http.File, error) {
	return nil, os.ErrNotExist
}

func NewFileSystem(dir, fallback string) http.FileSystem {
	return &noopFileSystem{}
}
