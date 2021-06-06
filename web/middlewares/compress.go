package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	acceptEncoding = "Accept-Encoding"
	gzipEncoding   = "gzip"

	contentEncoding = "Content-Encoding"
	contentLength   = "Content-Length"
)

type compressResponseWriter struct {
	// http.ResponseWriter
	RecordableHTTPResponse
	io.Writer
}

func (cw *compressResponseWriter) Write(b []byte) (int, error) {
	return cw.Writer.Write(b)
}

func Gzip(next http.Handler) http.Handler {
	pool := sync.Pool{
		New: func() interface{} {
			w, _ := gzip.NewWriterLevel(nil, gzip.DefaultCompression)
			return w
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// detect whether compress is needed
		if !strings.Contains(r.Header.Get(acceptEncoding), gzipEncoding) {
			next.ServeHTTP(w, r)
			return
		}

		rw := newRecordableResponse(w)
		rw.Header().Add("Vary", acceptEncoding)

		gw := pool.Get().(*gzip.Writer)
		gw.Reset(rw)
		defer func() {
			gw.Close()
			gw.Reset(nil)
			pool.Put(gw)
			rw.Header().Set(contentLength, strconv.Itoa(rw.Written()))
		}()

		w.Header().Set(contentEncoding, gzipEncoding)
		r.Header.Del(acceptEncoding)

		next.ServeHTTP(&compressResponseWriter{
			RecordableHTTPResponse: rw,
			Writer:                 gw,
		}, r)
	})
}
