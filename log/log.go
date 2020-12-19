package log

import (
	"fmt"
	"io"
	"time"

	zaplogfmt "github.com/jsternberg/zap-logfmt"
	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const TimeFormat = "2006-01-02T15:04:05.000000Z07:00"

func New(w io.Writer) (*zap.Logger, error) {
	cf := NewConfig()
	return cf.New(w)
}

func (c *Config) New(w io.Writer) (*zap.Logger, error) {
	format := c.Format
	if format == "" {
		if IsTerminal(w) {
			format = Console
		} else {
			format = Logfmt
		}
	}

	encoder, err := newEncoder(format)
	if err != nil {
		return nil, err
	}

	return zap.New(zapcore.NewCore(
		encoder,
		zapcore.Lock(zapcore.AddSync(w)),
		c.Level,
	)), nil
}

func newEncoder(format string) (zapcore.Encoder, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.Format(TimeFormat))
	}
	config.EncodeDuration = func(d time.Duration, encoder zapcore.PrimitiveArrayEncoder) {
		val := float64(d) / float64(time.Millisecond)
		encoder.AppendString(fmt.Sprintf("%.3fms", val))
	}
	config.LevelKey = "lvl"

	switch format {
	case JSON:
		return zapcore.NewJSONEncoder(config), nil
	case Console:
		return zapcore.NewConsoleEncoder(config), nil
	case Logfmt:
		return zaplogfmt.NewEncoder(config), nil
	default:
		return nil, fmt.Errorf("unknown logging format: %s", format)
	}
}

// IsTerminal checks if w is a file and whether it is an interactive terminal session.
func IsTerminal(w io.Writer) bool {
	if f, ok := w.(interface {
		Fd() uintptr
	}); ok {
		return isatty.IsTerminal(f.Fd())
	}
	return false
}
