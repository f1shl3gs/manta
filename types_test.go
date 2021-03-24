package manta_test

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestCell(t *testing.T) {
	var id manta.ID
	err := id.DecodeFromString("0000000000002b67")
	require.NoError(t, err)

	xyViewProperties := &manta.XYViewProperties{
		Type:       "xy",
		TimeFormat: "fff",
		Axes:       manta.Axes{},
	}

	cell := &manta.Cell{
		ID:             id,
		Name:           "xy",
		Desc:           "desc",
		W:              0,
		H:              2,
		X:              4,
		Y:              2,
		ViewProperties: xyViewProperties,

		/*ViewProperties: &manta.Cell_XY{
			XY: xyViewProperties,
		},*/
	}

	/*
		fmt.Println("i", i)
		size := m.ViewProperties.Size()
		fmt.Println("size", size)
		i -= size
		fmt.Println("i-size", i)
		fmt.Println("data size", len(dAtA[i:]))
	*/

	t.Run("proto marshal/unmarshal", func(t *testing.T) {
		data, err := cell.Marshal()
		require.NoError(t, err)
		// whole xy: 963b4410663a48680e506c853bbc5b51
		// original: 7c98a94789c0ef6144450350d350fc0a
		// size: 36
		// i 36
		// size 17
		// i-size 19
		// data size 17

		// whole xy: 963b4410663a48680e506c853bbc5b51
		// custom:   7c98a94789c0ef6144450350d350fc0a
		// size 36
		// i 36
		// size 17
		// i-size 19
		// data size 17
		fmt.Printf("%x\n", md5.Sum(data))

		unmarshal := true

		if unmarshal {
			var newCell manta.Cell
			err = newCell.Unmarshal(data)
			require.NoError(t, err)
		}
	})

	t.Run("marshal", func(t *testing.T) {
		txt, err := json.MarshalIndent(cell, "", "  ")
		require.NoError(t, err)
		fmt.Println(string(txt))
	})

	t.Run("unmarshal", func(t *testing.T) {
		txt := `{
"id": "0000000000002b67",
  "name": "xy",
  "desc": "desc",
  "x": 4,
  "h": 2,
"y": 2,
"w": 0,
  "viewProperties": {
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

		require.Equal(t, cell, p)
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

func TestCheckDecode(t *testing.T) {
	text := `{
  "name": "foo",
  "desc": "bar",
  "expr": "up",
  "cron": "@every 1m",
  "status": "active",
  "conditions": [
    {
      "status": "warn",
      "pending": "300s",
      "threshold": {
        "type": "lt",
        "value": 10
      }
    }
  ]
}`

	c := &manta.Check{}
	err := json.NewDecoder(bytes.NewBufferString(text)).Decode(c)
	require.NoError(t, err)
}

func TestConditionMarshal(t *testing.T) {
	c := &manta.Condition{
		Status:  "warn",
		Pending: 60 * time.Second,
		Threshold: manta.Threshold{
			Type:  "lt",
			Value: 10,
		},
	}

	err := json.NewEncoder(os.Stdout).Encode(c)
	require.NoError(t, err)
}

func TestDecodeCondition(t *testing.T) {
	text := `{
      "status": "CRIT",
      "threshold": {
        "type": "gt",
        "value": 0
      }
    }`

	c := &manta.Condition{}
	err := json.Unmarshal([]byte(text), c)
	require.NoError(t, err)
}

type content struct {
	Array []string `json:"array"`
}

func TestEmptySlice(t *testing.T) {
	type content struct {
		Array []string `json:"array"`
	}

	t.Run("marshal", func(t *testing.T) {
		c := &content{}
		data, err := json.Marshal(c)
		require.NoError(t, err)

		fmt.Println(string(data))
	})

	t.Run("unmarshal", func(t *testing.T) {
		text := `{"array":[]}`
		c := &content{}

		err := json.Unmarshal([]byte(text), c)
		require.NoError(t, err)

		data, err := json.Marshal(c)
		require.NoError(t, err)

		fmt.Println(string(data))
	})
}
