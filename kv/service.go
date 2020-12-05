package kv

import (
	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/snowflake"
	"github.com/f1shl3gs/manta/resource"
	"go.uber.org/zap"
)

type Service struct {
	kv Store

	logger   *zap.Logger
	idGen    manta.IDGenerator
	tokenGen manta.TokenGenerator
	audit    resource.Logger
}

func NewService(logger *zap.Logger, kv Store) *Service {
	return &Service{
		kv:     kv,
		logger: logger.With(zap.String("service", "kv")),
		idGen:  snowflake.NewIDGenerator(),
	}
}
