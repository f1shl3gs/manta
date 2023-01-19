package pb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternalRequestMarshal(t *testing.T) {
	req := &InternalRequest{
		ID: 1010,
		CreateBucket: &CreateBucket{
			Name: []byte("foobar"),
		},
	}

	data, err := req.Marshal()
	assert.NoError(t, err)

	decodeReq := &InternalRequest{}
	err = decodeReq.Unmarshal(data)
	assert.NoError(t, err)

	assert.Equal(t, req, decodeReq)
}
