package http

import (
	"context"
	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
	"net/http"
)

// ErrorCodeToStatusCode maps an manta error code string to a
// http status code integer.
func ErrorCodeToStatusCode(ctx context.Context, code string) int {
	// If the client disconnects early or times out then return a different
	// error than the passed in error code. Client timeouts return a 408
	// while disconnections return a non-standard Nginx HTTP 499 code.
	if err := ctx.Err(); err == context.DeadlineExceeded {
		return http.StatusRequestTimeout
	} else if err == context.Canceled {
		return 499 // https://httpstatuses.com/499
	}

	// Otherwise map internal error codes to HTTP status codes.
	statusCode, ok := mantaErrorToStatusCode[code]
	if ok {
		return statusCode
	}
	return http.StatusInternalServerError
}

// mantaErrorToStatusCode is a mapping of ErrorCode to http status code.
var mantaErrorToStatusCode = map[string]int{
	manta.EInternal:            http.StatusInternalServerError,
	manta.ENotImplemented:      http.StatusNotImplemented,
	manta.EInvalid:             http.StatusBadRequest,
	manta.EUnprocessableEntity: http.StatusUnprocessableEntity,
	manta.EEmptyValue:          http.StatusBadRequest,
	manta.EConflict:            http.StatusUnprocessableEntity,
	manta.ENotFound:            http.StatusNotFound,
	manta.EUnavailable:         http.StatusServiceUnavailable,
	manta.EForbidden:           http.StatusForbidden,
	manta.ETooManyRequests:     http.StatusTooManyRequests,
	manta.EUnauthorized:        http.StatusUnauthorized,
	manta.EMethodNotAllowed:    http.StatusMethodNotAllowed,
	manta.ETooLarge:            http.StatusRequestEntityTooLarge,
}

var httpStatusCodeToMantaError = map[int]string{}

func init() {
	for k, v := range mantaErrorToStatusCode {
		httpStatusCodeToMantaError[v] = k
	}
}

func logEncodingError(log *zap.Logger, r *http.Request, err error) {
	// If we encounter an error while encoding the response to an http request
	// the best thing we can do is log that error, as we may have already written
	// the headers for the http request in question.
	log.Info("Error encoding response",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
		zap.String("remote", r.RemoteAddr),
		zap.Error(err))
}
