package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	year = 365 * 24 * time.Hour
	week = 7 * 24 * time.Hour
	day  = 24 * time.Hour
)

func DurationLiteral(key string, val time.Duration) zapcore.Field {
	if val == 0 {
		return zap.String(key, "0s")
	}

	var (
		value int
		unit  string
	)
	switch {
	case val%year == 0:
		value = int(val / year)
		unit = "y"
	case val%week == 0:
		value = int(val / week)
		unit = "w"
	case val%day == 0:
		value = int(val / day)
		unit = "d"
	case val%time.Hour == 0:
		value = int(val / time.Hour)
		unit = "h"
	case val%time.Minute == 0:
		value = int(val / time.Minute)
		unit = "m"
	case val%time.Second == 0:
		value = int(val / time.Second)
		unit = "s"
	default:
		value = int(val / time.Millisecond)
		unit = "ms"
	}
	return zap.String(key, fmt.Sprintf("%d%s", value, unit))
}
