package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f1shl3gs/manta"
	"github.com/stretchr/testify/require"
)

func TestSecretHandler_handlePut(t *testing.T) {
	tests := []struct {
		name             string
		urlGen           func(orgID manta.ID) string
		body             interface{}
		expectStatusCode int
	}{
		{
			name: "create a secret",
			urlGen: func(orgID manta.ID) string {
				return SecretPrefix + "?orgID=" + orgID.String()
			},
			body: map[string]string{
				"key":   "foo",
				"value": "bar",
			},
			expectStatusCode: http.StatusCreated,
		},
		{
			name: "create secret without org",
			urlGen: func(orgID manta.ID) string {
				return SecretPrefix
			},
			body: map[string]string{
				"key":   "foo",
				"value": "bar",
			},
			expectStatusCode: http.StatusBadRequest,
		},
		{
			name: "create secret with invalid body",
			urlGen: func(orgID manta.ID) string {
				return SecretPrefix + "?orgID=" + orgID.String()
			},
			body: map[string]string{
				"key":  "foo",
				"aaaa": "bar",
			},
			expectStatusCode: http.StatusBadRequest,
		},
		{
			name: "create secret with not exists org",
			urlGen: func(orgID manta.ID) string {
				return SecretPrefix + "?orgID=0000000000000001"
			},
			body: map[string]string{
				"key":   "foo",
				"value": "bar",
			},
			expectStatusCode: http.StatusNotFound,
		},
		{
			name: "create secret with invalid key",
			urlGen: func(orgID manta.ID) string {
				return SecretPrefix + "?orgID=" + orgID.String()
			},
			body: map[string]string{
				"key":   "fo o",
				"value": "bar",
			},
			expectStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init
			backend, logger, orgID, closer := NewTestBackendWithOrg(t)
			defer closer()
			router := NewRouter()
			NewSecretHandler(logger, router, backend.SecretService)

			url := tt.urlGen(orgID)
			// build request and start testing
			buf := bytes.NewBuffer(nil)
			err := json.NewEncoder(buf).Encode(tt.body)
			require.NoError(t, err)

			r := httptest.NewRequest(http.MethodPut, url, buf)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)
			require.Equalf(t, tt.expectStatusCode, w.Code, "body: %s", w.Body.String())
		})
	}
}
