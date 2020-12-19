package manta_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestPanel(t *testing.T) {
	panel := &manta.Cell{
		Name:        "xy",
		Description: "desc",
		W:           0,
		H:           2,
		X:           4,
		Properties: &manta.XYView{
			Queries: []manta.Query{
				{},
			},
			Type:       "xy",
			TimeFormat: "fff",
			Axes:       manta.Axes{},
		},
		/*Properties: &manta.Panel_XY{
			XY: &manta.XYView{
				Type:       "xy",
				TimeFormat: "fff",
				Axes:       manta.Axes{},
			},
		},*/
	}

	t.Run("marshal", func(t *testing.T) {
		txt, err := json.MarshalIndent(panel, "", "  ")
		require.NoError(t, err)
		fmt.Println(string(txt))
	})

	t.Run("unmarshal", func(t *testing.T) {
		txt := `{
  "name": "xy",
  "description": "desc",
  "x": 4,
  "h": 2,
  "properties": {
    "type": "xy",
    "axes": {
      "x": {},
      "y": {}
    },
    "queries": [
      {}
    ],
    "timeFormat": "fff"
  }
}`
		p := &manta.Cell{}
		err := json.Unmarshal([]byte(txt), p)
		require.NoError(t, err)

		require.Equal(t, panel, p)
	})
}

func BenchmarkEventMarshal(b *testing.B) {
	idGen := snowflake.NewIDGenerator()
	now := time.Now()
	total := 0

	ev := &manta.Event{
		ID:    idGen.ID(),
		Name:  "foo",
		Start: now.Add(-time.Minute),
		End:   now,
		OrgID: idGen.ID(),
		Labels: map[string]string{
			"foo":  "bar",
			"foo1": "bar",
			"foo2": "bar",
			"foo3": "bar",
			"foo4": "bar",
			"foo5": "bar",
			"foo6": "bar",
		},
		Annotations: map[string]string{
			"foo":  "bar",
			"foo1": "bar",
			"foo2": "bar",
			"foo3": "bar",
			"foo4": "bar",
			"foo5": "bar",
			"foo6": "barrrrrrrrrrrrrrr",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := ev.Marshal()
		if err != nil {
			panic(err)
		}

		total += len(data)
	}

	b.SetBytes(int64(total / b.N))
}
