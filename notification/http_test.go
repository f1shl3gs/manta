package notification_test

import (
	"testing"

	"github.com/f1shl3gs/manta/notification"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	const text = `{"id": "0aa7b90c86a2e002", "name":"sss","desc":"dfadsfadfad","orgID":"0aa7b90c86a2e001","status":"NotStarted","type":"http","method":"POST","url":"ssssfdafsafafdsfs","headers":{},"authMethod":"none"}`
	_, err := notification.UnmarshalJSON([]byte(text))
	assert.NoError(t, err)
}
