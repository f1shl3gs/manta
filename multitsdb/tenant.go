package multitsdb

import (
	"context"

	"github.com/prometheus/prometheus/scrape"
	"github.com/prometheus/prometheus/storage"
	promv1 "github.com/prometheus/prometheus/web/api/v1"

	"github.com/f1shl3gs/manta"
)

type TenantStorage interface {
	Queryable(ctx context.Context, id manta.ID) (storage.Queryable, error)

	Appendable(ctx context.Context, id manta.ID) (storage.Appendable, error)
}

type Noop struct{}

func (n *Noop) Queryable(ctx context.Context, id manta.ID) (storage.Queryable, error) {
	return nil, ErrNotReady
}

func (n *Noop) Appendable(ctx context.Context, id manta.ID) (storage.Appendable, error) {
	return nil, ErrNotReady
}

type TenantTargetRetriever interface {
	TargetsActive(id manta.ID) map[string][]*scrape.Target
	TargetsDropped(id manta.ID) map[string][]*scrape.Target
}

type TargetRetrievers struct {
	retrievers map[manta.ID]promv1.TargetRetriever
}

func (t *TargetRetrievers) TargetsActive(id manta.ID) map[string][]*scrape.Target {
	retriever, exist := t.retrievers[id]
	if !exist {
		return nil
	}

	return retriever.TargetsActive()
}

func (t *TargetRetrievers) TargetsDropped(id manta.ID) map[string][]*scrape.Target {
	retriever, exist := t.retrievers[id]
	if !exist {
		return nil
	}

	return retriever.TargetsDropped()
}

func (t *TargetRetrievers) Add(id manta.ID, r promv1.TargetRetriever) {
	if t.retrievers == nil {
		t.retrievers = make(map[manta.ID]promv1.TargetRetriever)
	}

	t.retrievers[id] = r
}
