package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWatch(t *testing.T) {
	t.SkipNow()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8088/api/v1/configs/0a6a1631da5aa000?orgID=0a6723777c969000&watch=true", nil)
	assert.NoError(t, err)

	req.Header.Set("Cookie", "manta_session=0a6724103b569000")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	reader := httputil.NewChunkedReader(resp.Body)
	buf := make([]byte, 16*1024)
	for {
		n, err := reader.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				return
			}

			panic(err)
		}

		fmt.Println(string(buf[:n]))
	}
}
