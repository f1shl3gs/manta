package http

import (
	"compress/gzip"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tarfs"
	"go.uber.org/zap"
)

const (
	assetsPrefix = "ui/build"
)

type AssetsHandler struct {
	fs     fs.FS
	logger *zap.Logger
}

func NewAssetsHandler(logger *zap.Logger) (*AssetsHandler, error) {
	f, err := manta.Assets.Open("assets.tgz")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	tfs, err := tarfs.New(gr)
	if err != nil {
		return nil, err
	}

	sfs, err := fs.Sub(tfs, assetsPrefix)
	if err != nil {
		return nil, err
	}

	return &AssetsHandler{
		fs:     sfs,
		logger: logger,
	}, nil
}

func (a *AssetsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileServer := http.FileServer(http.FS(a.fs))
	filename := strings.TrimPrefix(r.URL.Path, "/")
	f, err := a.fs.Open(filename)
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	if err == nil {
		if err = addCacheHeaderFromFile(f, w); err != nil {
			panic(err)
		}

		fileServer.ServeHTTP(w, r)
		return
	}

	r.URL.Path = "/"
	filename = "index.html"
	if err = addCacheHeaders(a.fs, filename, w); err != nil {
		panic(err)
	}
	fileServer.ServeHTTP(w, r)
}

// addCacheHeaders requests an hour of Cache-Control and sets an ETag based on file size and modtime
func addCacheHeaders(sfs fs.FS, filename string, w http.ResponseWriter) error {
	file, err := sfs.Open(filename)
	if err != nil {
		return err
	}

	return addCacheHeaderFromFile(file, w)
}

func addCacheHeaderFromFile(f fs.File, w http.ResponseWriter) error {
	fi, err := f.Stat()
	if err != nil {
		return err
	}

	w.Header().Add("Cache-Control", "public, max-age=3600")
	hour, minute, second := fi.ModTime().Clock()
	etag := fmt.Sprintf(`%d%d%d%d%d`, fi.Size(), fi.ModTime().Day(), hour, minute, second)

	w.Header().Set("ETag", etag)
	return nil
}
