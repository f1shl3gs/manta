package kv

import (
	"bytes"
	"fmt"
	"github.com/f1shl3gs/manta"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httputil"
	"testing"
)

func TestCraeteOne(t *testing.T) {
	var orgID manta.ID
	err := orgID.DecodeFromString("0a5255af55928000")
	assert.NoError(t, err)

	var body = bytes.NewBufferString(`{
"orgID": "0a5255af55928000",
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

	req.Header.Set("Cookie", "manta_session=0a5255af59928000")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	data, err := httputil.DumpResponse(resp, true)
	assert.NoError(t, err)
	fmt.Println(string(data))
}
