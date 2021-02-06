package ingest

import (
	"github.com/thanos-io/thanos/pkg/receive"
)

type Ingester interface {
	receive.TenantStorage
}
