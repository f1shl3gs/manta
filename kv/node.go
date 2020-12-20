package kv

import (
	"context"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	nodeBucket             = []byte("node")
	nodeAddressIndexBucket = []byte("nodeaddrindex")
	nodeOrgIndexBucket     = []byte("nodeorgindex")
)

func (s *Service) FindNodeByID(ctx context.Context, id manta.ID) (*manta.Node, error) {
	var (
		node *manta.Node
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		node, err = s.findNodeByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return node, nil
}

func (s *Service) findNodeByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Node, error) {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	key, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(nodeBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	node := &manta.Node{}
	if err = node.Unmarshal(data); err != nil {
		return nil, err
	}

	return node, nil
}

func (s *Service) FindNodes(ctx context.Context, filter manta.NodeFilter, opt ...manta.FindOptions) ([]*manta.Node, int, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		nodes []*manta.Node
		err   error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		if filter.Address != nil {
			nodes, err = s.findNodeByAddress(ctx, tx, *filter.Address)
			return err
		}

		return nil
	})

	return nodes, len(nodes), nil
}

func (s *Service) findNodeByOrgID(ctx context.Context, tx Tx, orgID manta.ID) ([]*manta.Node, error) {
	prefix, err := orgID.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(nodeOrgIndexBucket)
	if err != nil {
		return nil, err
	}

	cur, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix))
	if err != nil {
		return nil, err
	}

	keys := make([][]byte, 0, 8)
	err = WalkCursor(ctx, cur, func(k, v []byte) error {
		keys = append(keys, v)
		return nil
	})
	if err != nil {
		return nil, err
	}

	values, err := b.GetBatch(keys...)
	if err != nil {
		return nil, err
	}

	nodes := make([]*manta.Node, 0, len(values))
	for i := 0; i < len(values); i++ {
		s := &manta.Node{}
		if err = s.Unmarshal(values[i]); err != nil {
			return nil, err
		}

		nodes[i] = s
	}

	return nodes, nil
}

func (s *Service) findNodeByAddress(ctx context.Context, tx Tx, addr string) ([]*manta.Node, error) {
	b, err := tx.Bucket(nodeAddressIndexBucket)
	if err != nil {
		return nil, err
	}

	prefix := []byte(addr)
	cursor, err := b.ForwardCursor(prefix, WithCursorPrefix(prefix))
	if err != nil {
		return nil, err
	}

	keys := make([][]byte, 0, 10)
	err = WalkCursor(ctx, cursor, func(k, v []byte) error {
		keys = append(keys, v)
		return nil
	})
	if err != nil {
		return nil, err
	}

	values, err := b.GetBatch(keys...)
	if err != nil {
		return nil, err
	}

	list := make([]*manta.Node, 0, len(values))
	for i := 0; i < len(values); i++ {
		s := &manta.Node{}
		err = s.Unmarshal(values[i])
		if err != nil {
			return nil, err
		}

		list = append(list, s)
	}

	return list, nil
}

func (s *Service) CreateNode(ctx context.Context, node *manta.Node) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.kv.Update(ctx, func(tx Tx) error {
		if _, err := s.findOrganizationByID(ctx, tx, node.OrgID); err != nil {
			return err
		}

		return s.createNode(ctx, tx, node)
	})
}

func (s *Service) createNode(ctx context.Context, tx Tx, node *manta.Node) error {
	node.ID = s.idGen.ID()
	node.Created = time.Now()
	node.Updated = time.Now()

	return s.putNode(ctx, tx, node)
}

func (s *Service) UpdateNode(ctx context.Context, id manta.ID, u manta.NodeUpdate) (*manta.Node, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		node *manta.Node
		err  error
	)

	err = s.kv.Update(ctx, func(tx Tx) error {
		node, err = s.findNodeByID(ctx, tx, id)
		if err != nil {
			return err
		}

		err = s.deleteNode(ctx, tx, id)
		if err != nil {
			return err
		}

		if u.Env != nil {
			node.Env = *u.Env
		}

		node.Updated = time.Now()
		return s.putNode(ctx, tx, node)
	})

	if err != nil {
		return nil, err
	}

	return node, nil
}

func (s *Service) DeleteNode(ctx context.Context, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.deleteNode(ctx, tx, id)
	})
}

func (s *Service) deleteNode(ctx context.Context, tx Tx, id manta.ID) error {
	node, err := s.findNodeByID(ctx, tx, id)
	if err != nil {
		return err
	}

	pk, err := node.ID.Encode()
	if err != nil {
		return err
	}

	// address index
	addrIdx := IndexKey([]byte(node.Address), pk)
	b, err := tx.Bucket(nodeAddressIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Delete(addrIdx); err != nil {
		return err
	}

	// organization index
	fk, err := node.OrgID.Encode()
	if err != nil {
		return err
	}

	refIdx := IndexKey(fk, pk)
	b, err = tx.Bucket(nodeOrgIndexBucket)
	if err != nil {
		return err
	}

	if err = b.Delete(refIdx); err != nil {
		return err
	}

	// delete node
	b, err = tx.Bucket(nodeBucket)
	if err != nil {
		return err
	}

	return b.Delete(pk)
}

func (s *Service) putNode(ctx context.Context, tx Tx, node *manta.Node) error {
	pk, err := node.ID.Encode()
	if err != nil {
		return err
	}

	// address index
	addrIdx := IndexKey([]byte(node.Address), pk)
	b, err := tx.Bucket(nodeAddressIndexBucket)
	if err != nil {
		return err
	}

	if err := b.Put(addrIdx, pk); err != nil {
		return err
	}

	// organization index
	fk, err := node.OrgID.Encode()
	if err != nil {
		return err
	}
	refIdx := IndexKey(fk, pk)
	b, err = tx.Bucket(nodeOrgIndexBucket)
	if err != nil {
		return err
	}

	err = b.Put(refIdx, pk)
	if err != nil {
		return err
	}

	// save node
	data, err := node.Marshal()
	if err != nil {
		return err
	}

	b, err = tx.Bucket(nodeBucket)
	if err != nil {
		return err
	}

	return b.Put(pk, data)
}
