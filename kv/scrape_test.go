package kv

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCraeteOne(t *testing.T) {
	var body = bytes.NewBufferString(`{
"orgID": "0a609b764eeae000",
"name": "demo",
"targets": [
"localhost:8088"
],
"labels": {
"foo": "bar"
}
}`)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8088/api/v1/scrapes", body)
	assert.NoError(t, err)

	req.Header.Set("Cookie", "manta_session=0a609b7653aae000")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	data, err := httputil.DumpResponse(resp, true)
	assert.NoError(t, err)
	fmt.Println(string(data))
}
