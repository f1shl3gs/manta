package kv

import (
	"context"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

var (
	templateBucket          = []byte("templates")
	templateNameIndexBucket = []byte("templatenameindex")
)

func (s *Service) FindTemplateByID(ctx context.Context, id manta.ID) (*manta.Template, error) {
	var (
		tmpl *manta.Template
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		tmpl, err = s.findTemplateByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (s *Service) findTemplateByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Template, error) {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	pk, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(templateBucket)
	if err != nil {
		return nil, err
	}

	value, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	tmpl := &manta.Template{}
	if err = tmpl.Unmarshal(value); err != nil {
		return nil, err
	}

	return tmpl, err
}

func (s *Service) FindTemplateByName(ctx context.Context, name string) (*manta.Template, error) {
	var (
		tmpl *manta.Template
		err  error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		tmpl, err = s.findTemplateByName(ctx, tx, name)
		return err
	})

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (s *Service) findTemplateByName(ctx context.Context, tx Tx, name string) (*manta.Template, error) {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	b, err := tx.Bucket(templateNameIndexBucket)
	if err != nil {
		return nil, err
	}

	key, err := b.Get([]byte(name))
	if err != nil {
		return nil, err
	}

	b, err = tx.Bucket(templateBucket)
	if err != nil {
		return nil, err
	}

	val, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	tmpl := &manta.Template{}
	if err = tmpl.Unmarshal(val); err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (s *Service) FindTemplates(ctx context.Context, filter manta.TemplateFilter, opt ...manta.FindOptions) ([]*manta.Template, int, error) {
	panic("implement me")
}

func (s *Service) CreateTemplate(ctx context.Context, template *manta.Template) error {
	panic("implement me")
}

func (s *Service) UpdateTemplate(ctx context.Context, id manta.ID, u manta.TemplateUpdate) (*manta.Template, error) {
	panic("implement me")
}

func (s *Service) DeleteTemplate(ctx context.Context, id manta.ID) error {
	panic("implement me")
}
