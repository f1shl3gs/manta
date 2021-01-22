package kmsgreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
)

type kmsgReceiver struct {
}

func (k *kmsgReceiver) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (k *kmsgReceiver) Shutdown(ctx context.Context) error {
	return nil
}

func newKmsgReceiver() *kmsgReceiver {
	return &kmsgReceiver{}
}
