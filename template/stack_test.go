package template

import (
	"testing"

	"github.com/f1shl3gs/manta"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestUnmarshal(t *testing.T) {
	text := `apiVersion: manta/v1
kind: Scraper
name: Selfstat
spec:
  name: Selfstat
  desc: collect metrics from manta itself
  targets:
    - 127.0.0.1:8088
  labels:
    foo: bar`

	obj := &Object{}
	err := yaml.Unmarshal([]byte(text), obj)
	require.NoError(t, err)
	st, ok := obj.Spec.(*manta.ScrapeTarget)
	require.True(t, ok)
	require.Equal(t, "Selfstat", st.Name)
}
