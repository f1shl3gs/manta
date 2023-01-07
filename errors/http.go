package errors

import (
	"context"
	"net/http"
)

// HTTPErrorHandler is the interface to handle http error.
type HTTPErrorHandler interface {
	HandleHTTPError(ctx context.Context, err error, w http.ResponseWriter)
}
