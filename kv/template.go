package kv

import (
    "bytes"
    "context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/template"
)

var (
	// key is orgID + id
	templatesBucket = []byte("templates")
)

func decodeResourceErr(index int, typ template.ResourceType, err error) *manta.Error {
	return &manta.Error{
		Code: manta.EInvalid,
		Msg:  fmt.Sprintf("resource[%d] cannot be decoded into %s", index, typ),
		Err:  err,
	}
}

func (s *Service) Install(ctx context.Context, create template.TemplateCreate) (*template.Template, error) {
	if len(create.Resources) == 0 {
		return nil, template.ErrNoResource
	}

	var tmpl = &template.Template{
		ID:      s.idGen.ID(),
		Name:    create.Name,
		Desc:    create.Desc,
		Created: time.Now(),
		OrgID:   create.OrgID,
	}

	err := s.kv.Update(ctx, func(tx Tx) error {
		for i, res := range create.Resources {
			switch res.Type {
			case template.ResourceCheck:
				check := &manta.Check{}
				err := json.Unmarshal(res.Spec, check)
				if err != nil {
					return decodeResourceErr(i, res.Type, err)
				}

				err = s.createCheck(ctx, tx, check)
				if err != nil {
					return err
				}

				tmpl.Resources = append(tmpl.Resources, template.ResourceItem{
					ID:   check.ID,
					Type: res.Type,
					Name: check.Name,
				})

			case template.ResourceConfig:
				cf := &manta.Configuration{}
				err := json.Unmarshal(res.Spec, cf)
				if err != nil {
					return decodeResourceErr(i, res.Type, err)
				}

				err = s.createConfiguration(ctx, tx, cf)
				if err != nil {
					return err
				}

				tmpl.Resources = append(tmpl.Resources, template.ResourceItem{
					ID:   cf.ID,
                    Type: res.Type,
					Name: cf.Name,
				})

			case template.ResourceDashboard:
				dashboard := &manta.Dashboard{}
				err := json.Unmarshal(res.Spec, dashboard)
				if err != nil {
					return decodeResourceErr(i, res.Type, err)
				}

				err = s.createDashboard(ctx, tx, dashboard)
				if err != nil {
					return err
				}

				tmpl.Resources = append(tmpl.Resources, template.ResourceItem{
					ID:   dashboard.ID,
                    Type: res.Type,
					Name: dashboard.Name,
				})

			case template.ResourceScrape:
				scrape := &manta.ScrapeTarget{}
				err := json.Unmarshal(res.Spec, scrape)
				if err != nil {
					return decodeResourceErr(i, res.Type, err)
				}

				err = s.createScrapeTarget(ctx, tx, scrape)
				if err != nil {
					return err
				}

				tmpl.Resources = append(tmpl.Resources, template.ResourceItem{
					ID:   scrape.ID,
                    Type: res.Type,
					Name: scrape.Name,
				})

			default:
				return &manta.Error{
					Code: manta.EInvalid,
					Msg:  fmt.Sprintf("unknown resource type %s at resources[%d]", res.Type, i),
				}
			}
		}

		data, err := json.Marshal(tmpl)
		if err != nil {
			return err
		}

		b, err := tx.Bucket(templatesBucket)
		if err != nil {
			return err
		}

		key, err := indexIDKey(tmpl.ID, tmpl.OrgID)
		if err != nil {
			return err
		}

		return b.Put(key, data)
	})

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (s *Service) Uninstall(ctx context.Context, orgID, id manta.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		key, err := indexIDKey(id, orgID)
		if err != nil {
			return err
		}

		b, err := tx.Bucket(templatesBucket)
		if err != nil {
			return err
		}

		value, err := b.Get(key)
		if err != nil {
			return err
		}

		var tmpl = &template.Template{}
		err = json.Unmarshal(value, tmpl)
		if err != nil {
			return err
		}

		for _, res := range tmpl.Resources {
			switch res.Type {
            case template.ResourceCheck:
                err = s.deleteCheck(tx, id)
                if err == ErrKeyNotFound {
                    // resource already deleted
                    continue
                }

                if err != nil {
                    return err
                }

            case template.ResourceConfig:
                err = s.deleteConfig(tx, id)
                if err == ErrKeyNotFound {
                    // resource already deleted
                    continue
                }

                if err != nil {
                    return err
                }

            case template.ResourceDashboard:
                err = s.deleteDashboard(ctx, tx, id)
                if err == ErrKeyNotFound {
                    // resource already deleted
                    continue
                }

                if err != nil {
                    return err
                }

            case template.ResourceScrape:
                err = s.deleteScrapeTarget(tx, id)
                if err == ErrKeyNotFound {
                    // resource already deleted
                    continue
                }

                if err != nil {
                    return err
                }

            default:
                return &manta.Error{
                    Code: manta.EInvalid,
                    Msg: fmt.Sprintf("unsupport resource type %s, name: %s", res.Type, res.Name),
                }
			}
		}

        return nil
	})
}

func (s *Service) ListTemplate(ctx context.Context, orgID manta.ID) ([]*template.Template, error) {
    var templates []*template.Template

	err := s.kv.View(ctx, func(tx Tx) error {
        prefix, err := orgID.Encode()
        if err != nil {
            return err
        }

        b, err := tx.Bucket(templatesBucket)
        if err != nil {
            return err
        }

        cursor, err := b.Cursor()
        if err != nil {
            return err
        }

        for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
            tmpl := &template.Template{}
            err = json.Unmarshal(v, tmpl)
            if err != nil {
                return err
            }

            templates = append(templates, tmpl)
        }

        return nil
	})

    if err != nil {
        return nil, err
    }

    return templates, nil
}
