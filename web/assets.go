package web

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

type Assets struct {
	Prefix  string
	Default string

	logger *zap.Logger
}

func (a *Assets) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sfs, err := fs.Sub(manta.Assets, "ui/build")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(sfs))
	filename := strings.TrimPrefix(r.URL.Path, "/")
	f, err := sfs.Open(filename)
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	if err == nil {
		if err = addCacheHeaders(sfs, filename, w); err != nil {
			panic(err)
		}

		fileServer.ServeHTTP(w, r)
		return
	}

	r.URL.Path = "/"
	filename = "index.html"
	if err = addCacheHeaders(sfs, filename, w); err != nil {
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

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	w.Header().Add("Cache-Control", "public, max-age=3600")
	hour, minute, second := fi.ModTime().Clock()
	etag := fmt.Sprintf(`%d%d%d%d%d`, fi.Size(), fi.ModTime().Day(), hour, minute, second)

	w.Header().Set("ETag", etag)
	return nil
}
