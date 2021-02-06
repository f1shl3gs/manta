package ingest

import "go.uber.org/zap"

type logger struct {
	*zap.Logger
}

func (l *logger) Log(keyvals ...interface{}) error {

	return nil
}
