package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/f1shl3gs/manta"
)

// PlatformErrorCodeHeader shows the error code of platform error.
const PlatformErrorCodeHeader = "X-Platform-Error-Code"

// ErrorHandler is the error handler in http package.
type ErrorHandler int

// HandleHTTPError encodes err with the appropriate status code and format,
// sets the X-Platform-Error-Code headers on the response.
// We're no longer using X-Manta-Error and X-Manta-Reference.
// and sets the response status to the corresponding status code.
func (h ErrorHandler) HandleHTTPError(ctx context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		return
	}

	code := manta.ErrorCode(err)
	w.Header().Set(PlatformErrorCodeHeader, code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(ErrorCodeToStatusCode(ctx, code))
	var e struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	e.Code = manta.ErrorCode(err)
	if err, ok := err.(*manta.Error); ok {
		e.Message = err.Error()
	} else {
		e.Message = "An internal error has occurred"
	}
	b, _ := json.Marshal(e)
	_, _ = w.Write(b)
}

// StatusCodeToErrorCode maps a http status code integer to an
// manta error code string.
func StatusCodeToErrorCode(statusCode int) string {
	errorCode, ok := httpStatusCodeToMantaError[statusCode]
	if ok {
		return errorCode
	}

	return manta.EInternal
}

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

// CheckErrorStatus for status and any error in the response.
func CheckErrorStatus(code int, res *http.Response) error {
	err := CheckError(res)
	if err != nil {
		return err
	}

	if res.StatusCode != code {
		return fmt.Errorf("unexpected status code: %s", res.Status)
	}

	return nil
}

// CheckError reads the http.Response and returns an error if one exists.
// It will automatically recognize the errors returned by Manta services
// and decode the error into an internal error type. If the error cannot
// be determined in that way, it will create a generic error message.
//
// If there is no error, then this returns nil.
func CheckError(resp *http.Response) (err error) {
	switch resp.StatusCode / 100 {
	case 4, 5:
		// We will attempt to parse this error outside of this block.
	case 2:
		return nil
	default:
		// TODO(jsternberg): Figure out what to do here?
		return &manta.Error{
			Code: manta.EInternal,
			Msg:  fmt.Sprintf("unexpected status code: %d %s", resp.StatusCode, resp.Status),
		}
	}

	perr := &manta.Error{
		Code: StatusCodeToErrorCode(resp.StatusCode),
	}

	if resp.StatusCode == http.StatusUnsupportedMediaType {
		perr.Msg = fmt.Sprintf("invalid media type: %q", resp.Header.Get("Content-Type"))
		return perr
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		// Assume JSON if there is no content-type.
		contentType = "application/json"
	}
	mediatype, _, _ := mime.ParseMediaType(contentType)

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		perr.Msg = "failed to read error response"
		perr.Err = err
		return perr
	}

	switch mediatype {
	case "application/json":
		if err := json.Unmarshal(buf.Bytes(), perr); err != nil {
			perr.Msg = fmt.Sprintf("attempted to unmarshal error as JSON but failed: %q", err)
			perr.Err = firstLineAsError(buf)
		}
	default:
		perr.Err = firstLineAsError(buf)
	}

	if perr.Code == "" {
		// given it was unset during attempt to unmarshal as JSON
		perr.Code = StatusCodeToErrorCode(resp.StatusCode)
	}

	return perr
}

func firstLineAsError(buf bytes.Buffer) error {
	line, _ := buf.ReadString('\n')
	return errors.New(strings.TrimSuffix(line, "\n"))
}
