package manta

import (
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func BenchmarkReadDashboard(b *testing.B) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/dashboards/06d4010163e3d000", nil)
	require.NoError(b, err)

	fetch := func() error {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return errors.New("unexpected status code")
		}

		return nil
	}

	for i := 0; i < b.N; i++ {
		err := fetch()
		if err != nil {
			b.Fatal(err)
		}
	}
}
