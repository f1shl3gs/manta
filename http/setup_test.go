package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/f1shl3gs/manta"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

func TestSetup(t *testing.T) {
	t.Run("setup", func(t *testing.T) {
		var (
			ctx      = context.Background()
			orgName  = "org"
			username = "foo"
			password = "password"
			service  = NewTestHTTPService(t)
		)

		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(&manta.OnBoardingRequest{
			Username:     username,
			Password:     password,
			Organization: orgName,
		})
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, setupPath, buf)
		w := httptest.NewRecorder()

		service.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())

		// check org and user
		backend := service.backend

		onboarded, err := backend.OnBoardingService.Onboarded(ctx)
		assert.NoError(t, err)
		assert.True(t, onboarded)

		org, err := backend.OrganizationService.FindOrganization(ctx, manta.OrganizationFilter{Name: &orgName})
		assert.NoError(t, err)
		assert.Equal(t, org.Name, orgName)

		user, err := backend.UserService.FindUser(ctx, manta.UserFilter{Name: &username})
		assert.NoError(t, err)
		assert.Equal(t, username, user.Name)

		err = backend.PasswordService.ComparePassword(ctx, user.ID, password)
		assert.NoError(t, err)
	})
}
