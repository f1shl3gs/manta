package manta

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck_Unmarshal(t *testing.T) {
	const text = `{"name":"name this check","desc":"desc this check","status":"active","cron":"@every 5s","orgID":"0a712029e1880000","expr":"rate(process_cpu_seconds_total[1m]) * 100","conditions":[{"status":"warn","pending":"10m","threshold":{"type":"gt","value":100}}]}`
	var check Check
	err := json.Unmarshal([]byte(text), &check)
	assert.NoError(t, err)
}
