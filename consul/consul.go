package consul

import (
	"context"

	consulapi "github.com/armon/consul-api"
	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"go.uber.org/zap"
)

type CatalogService struct {
	logger *zap.Logger

	organizationService manta.OrganizationService
	nodeService         manta.NodeService
	collectionService   manta.OtclService
}

// todo: the performance would not be good, cache layer might need,
// or find nodes and collections parallel
func (cs *CatalogService) Service(ctx context.Context, service string, tags []string) ([]*consulapi.ServiceEntry, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	orgs, _, err := cs.organizationService.FindOrganizations(ctx, manta.OrganizationFilter{})
	if err != nil {
		return nil, err
	}

	entries := make([]*consulapi.ServiceEntry, 0)

	for _, org := range orgs {
		nodes, _, err := cs.nodeService.FindNodes(ctx, manta.NodeFilter{OrgID: &org.ID})
		if err != nil {
			return nil, err
		}

		collections, err := cs.collectionService.FindCollections(ctx, manta.OtclFilter{})
		if err != nil {
			return nil, err
		}

		for _, node := range nodes {
			for _, c := range collections {
				entries = append(entries, newServiceEntry(org, node, c))
			}
		}
	}

	return entries, nil
}

func newServiceEntry(org *manta.Organization, node *manta.Node, c *manta.Collection) *consulapi.ServiceEntry {
	return nil
}
