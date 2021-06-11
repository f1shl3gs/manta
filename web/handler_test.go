package web

import (
	"net/http"
	"testing"

	"github.com/f1shl3gs/manta"
	"github.com/stretchr/testify/require"
)

func mustIDFromString(s string) manta.ID {
	id, err := manta.IDFromString(s)
	if err != nil {
		panic(err)
	}

	return *id
}

func TestOrgIDFromRequest(t *testing.T) {
	tests := []struct {
		name string
		url  string
		id   manta.ID
		err  error
	}{
		{
			name: "valid",
			url:  "/foo/bar?orgID=0767fada06e3d000",
			id:   mustIDFromString("0767fada06e3d000"),
			err:  nil,
		},
		{
			name: "no orgID",
			url:  "/foo/bar",
			id:   0,
			err: &manta.Error{
				Code: manta.EInvalid,
				Msg:  "invalid orgID from url",
				Op:   "decode orgID",
				Err:  manta.ErrInvalidIDLength,
			},
		},
		{
			name: "invalid id length",
			url:  "/foo/bar?orgID=0767fada06e3d0",
			id:   0,
			err: &manta.Error{
				Code: manta.EInvalid,
				Msg:  "invalid orgID from url",
				Op:   "decode orgID",
				Err:  manta.ErrInvalidIDLength,
			},
		},
		{
			name: "zero",
			url:  "/foo/bar?orgID=0",
			id:   0,
			err: &manta.Error{
				Code: manta.EInvalid,
				Msg:  "invalid orgID from url",
				Op:   "decode orgID",
				Err:  manta.ErrInvalidIDLength,
			},
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodGet, tt.url, nil)
		require.NoError(t, err)

		id, err := orgIDFromRequest(req)
		require.EqualValues(t, tt.err, err)
		require.Equal(t, tt.id, id)
	}
}
