package labelset

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalJSONLabelSet(t *testing.T) {
	type testConfig struct {
		LabelSet LabelSet `yaml:"labelSet,omitempty"`
	}

	// valid LabelSet JSON
	labelSetJSON := `{
	"labelSet": {
		"monitor": "codelab",
		"foo": "bar"
	}
}`
	var c testConfig
	err := json.Unmarshal([]byte(labelSetJSON), &c)

	if err != nil {
		t.Errorf("unexpected error while marshalling JSON : %s", err.Error())
	}

	labelSetString := c.LabelSet.String()

	expected := `{foo="bar", monitor="codelab"}`

	if expected != labelSetString {
		t.Errorf("expected %s but got %s", expected, labelSetString)
	}

	// invalid LabelSet JSON
	invalidlabelSetJSON := `{
	"labelSet": {
		"1nvalid_23name": "codelab",
		"foo": "bar"
	}
}`

	err = json.Unmarshal([]byte(invalidlabelSetJSON), &c)
	expectedErr := `"1nvalid_23name" is not a valid label name`
	if err == nil || err.Error() != expectedErr {
		t.Errorf("expected an error with message '%s' to be thrown", expectedErr)
	}
}

func TestLabelSetClone(t *testing.T) {
	labelSet := LabelSet{
		"monitor": "codelab",
		"foo":     "bar",
		"bar":     "baz",
	}

	cloneSet := labelSet.Clone()

	if len(labelSet) != len(cloneSet) {
		t.Errorf("expected the length of the cloned Label set to be %d, but got %d",
			len(labelSet), len(cloneSet))
	}

	for ln, lv := range labelSet {
		expected := cloneSet[ln]
		if expected != lv {
			t.Errorf("expected to get LabelValue %s, but got %s for LabelName %s", expected, lv, ln)
		}
	}
}

func TestLabelSetMerge(t *testing.T) {
	labelSet := LabelSet{
		"monitor": "codelab",
		"foo":     "bar",
		"bar":     "baz",
	}

	labelSet2 := LabelSet{
		"monitor": "codelab",
		"dolor":   "mi",
		"lorem":   "ipsum",
	}

	expectedSet := LabelSet{
		"monitor": "codelab",
		"foo":     "bar",
		"bar":     "baz",
		"dolor":   "mi",
		"lorem":   "ipsum",
	}

	mergedSet := labelSet.Merge(labelSet2)

	if len(mergedSet) != len(expectedSet) {
		t.Errorf("expected the length of the cloned Label set to be %d, but got %d",
			len(expectedSet), len(mergedSet))
	}

	for ln, lv := range mergedSet {
		expected := expectedSet[ln]
		if expected != lv {
			t.Errorf("expected to get LabelValue %s, but got %s for LabelName %s", expected, lv, ln)
		}
	}

}
