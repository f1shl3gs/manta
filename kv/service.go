package kv

import (
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/snowflake"
	"github.com/f1shl3gs/manta/token"
)

type Service struct {
	kv Store

	logger   *zap.Logger
	idGen    manta.IDGenerator
	tokenGen token.Generator
}

type Option func(service *Service)

func WithIDGenerator(idGen manta.IDGenerator) Option {
	return func(svc *Service) {
		svc.idGen = idGen
	}
}

func NewService(logger *zap.Logger, kv Store, opts ...Option) *Service {
	svc := &Service{
		kv:       kv,
		logger:   logger.With(zap.String("service", "kv")),
		idGen:    snowflake.NewIDGenerator(),
		tokenGen: token.NewGenerator(0),
	}

	for _, fn := range opts {
		fn(svc)
	}

	return svc
}
