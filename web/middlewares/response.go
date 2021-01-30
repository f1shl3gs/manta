package middlewares

import "net/http"

type RecordableHTTPResponse interface {
	http.ResponseWriter

	Status() int

	Written() int
}

type recordableHTTPResponse struct {
	http.ResponseWriter
	status  int
	written int
}

func (r *recordableHTTPResponse) Write(data []byte) (int, error) {
	n, err := r.ResponseWriter.Write(data)
	r.written += n

	return n, err
}

func (r *recordableHTTPResponse) Status() int {
	if r.status == 0 {
		return http.StatusOK
	}

	return r.status
}

func (r *recordableHTTPResponse) Written() int {
	return r.written
}

func newRecordableResponse(w http.ResponseWriter) RecordableHTTPResponse {
	if r, ok := w.(RecordableHTTPResponse); ok {
		return r
	}

	return &recordableHTTPResponse{
		ResponseWriter: w,
		status:         0,
		written:        0,
	}
}
