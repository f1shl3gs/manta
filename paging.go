package manta

import (
	"net/http"
	"strconv"
)

const (
	DefaultPageSize = 20
)

// FindOptions represents options passed to all find methods with multiple results.
type FindOptions struct {
	Limit      int
	Offset     int
	SortBy     string
	Descending bool
}

// DecodeFindOptions returns a FindOptions decoded from a http request.
func DecodeFindOptions(r *http.Request) (FindOptions, error) {
	opts := FindOptions{}
	query := r.URL.Query()

	if value := query.Get("offset"); value != "" {
		offset, err := strconv.Atoi(value)
		if err != nil {
			return FindOptions{}, &Error{
				Code: EInvalid,
				Err:  err,
				Msg:  "invalid offset",
			}
		}

		opts.Offset = offset
	}

	if value := query.Get("limit"); value != "" {
		limit, err := strconv.Atoi(value)
		if err != nil {
			return FindOptions{}, &Error{
				Code: EInvalid,
				Msg:  "invalid limit",
				Err:  err,
			}
		}

		opts.Limit = limit
	} else {
		opts.Limit = DefaultPageSize
	}

	if value := query.Get("descending"); value != "" {
		descending, err := strconv.ParseBool(value)
		if err != nil {
			return FindOptions{}, &Error{
				Code: EInvalid,
				Msg:  "invalid descending",
				Err:  err,
			}
		}

		opts.Descending = descending
	}

	if value := query.Get("sortBy"); value != "" {
		opts.SortBy = value
	}

	return opts, nil
}
