package web

import (
	"bytes"
	"encoding/json"
	"github.com/f1shl3gs/manta"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrganization(t *testing.T) {
	backend, logger, closer := NewTestBackend(t)
	defer closer()

	router := NewRouter()
	NewOrganizationHandler(logger, router, backend)

	t.Run("create organization", func(t *testing.T) {
		org := &manta.Organization{
			Name: "foo",
			Desc: "foo desc",
		}
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(org)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, orgsPrefix, buf)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)

		created, err := decodeOrganization(w.Body)
		require.NoError(t, err, "decode response organization failed")

		require.Equal(t, org.Name, created.Name)
		require.Equal(t, org.Desc, created.Desc)
		require.Equal(t, true, created.ID.Valid())
	})
}

func TestDecode(t *testing.T) {
	data := `{"id":"0698cab45e059000","created":"2020-11-10T17:21:40.216261892Z","modified":"2020-11-10T17:21:40.216261962Z","name":"name","desc":"desc"}`
	buf := bytes.NewBuffer([]byte(data))

	org, err := decodeOrganization(buf)
	require.NoError(t, err)
	require.Equal(t, "name", org.Name)
}
