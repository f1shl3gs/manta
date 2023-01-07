package router

import (
	"context"
	"net/http"

	"github.com/f1shl3gs/manta/errors"
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
	errors.EInternal:            http.StatusInternalServerError,
	errors.ENotImplemented:      http.StatusNotImplemented,
	errors.EInvalid:             http.StatusBadRequest,
	errors.EUnprocessableEntity: http.StatusUnprocessableEntity,
	errors.EEmptyValue:          http.StatusBadRequest,
	errors.EConflict:            http.StatusUnprocessableEntity,
	errors.ENotFound:            http.StatusNotFound,
	errors.EUnavailable:         http.StatusServiceUnavailable,
	errors.EForbidden:           http.StatusForbidden,
	errors.ETooManyRequests:     http.StatusTooManyRequests,
	errors.EUnauthorized:        http.StatusUnauthorized,
	errors.EMethodNotAllowed:    http.StatusMethodNotAllowed,
	errors.ETooLarge:            http.StatusRequestEntityTooLarge,
}

var httpStatusCodeToMantaError = map[int]string{}

func init() {
	for k, v := range mantaErrorToStatusCode {
		httpStatusCodeToMantaError[v] = k
	}
}
