package manta

import (
	"context"
	"encoding/json"
	"time"
)

var (
    ErrNoResource = &Error{
        Code: EUnprocessableEntity,
		Msg:  "no resource found in template",
	}
)

type ResourceItem struct {
	ID   ID           `json:"id"`
	Type ResourceType `json:"type"`
	Name string       `json:"name"`
}

type Template struct {
	ID      ID        `json:"ID"`
	Name    string    `json:"name"`
	Desc    string    `json:"desc,omitempty"`
	OrgID   ID        `json:"orgID"`
	Created time.Time `json:"created"`

	Resources []ResourceItem `json:"resources"`
}

type TemplateCreate struct {
	Name      string `json:"name"`
	Desc      string `json:"desc,omitempty"`
	OrgID     ID     `json:"orgID"`
	Resources []struct {
		Type ResourceType    `json:"type"`
		Spec json.RawMessage `json:"spec"`
	} `json:"resources"`
}

type TemplateService interface {
	Install(ctx context.Context, tmpl TemplateCreate) (*Template, error)
	Uninstall(ctx context.Context, orgID, id ID) error

	ListTemplate(ctx context.Context, orgID ID) ([]*Template, error)
}
