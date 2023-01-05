package template

import (
	"context"
	"encoding/json"

	"github.com/f1shl3gs/manta"
)

var (
	ErrNoResource = &manta.Error{
		Code: manta.EUnprocessableEntity,
		Msg:  "no resource found in template",
	}
)

type Resource struct {
	Type ResourceType    `json:"type"`
	Spec json.RawMessage `json:"spec"`
}

type TemplateCreate struct {
	Name      string     `json:"name"`
	Desc      string     `json:"desc,omitempty"`
	OrgID     manta.ID   `json:"orgID"`
	Resources []Resource `json:"resources"`
}

type TemplateService interface {
	Install(ctx context.Context, tmpl TemplateCreate) (*Template, error)
	Uninstall(ctx context.Context, orgID, id manta.ID) error

	ListTemplate(ctx context.Context, orgID manta.ID) ([]*Template, error)
}
