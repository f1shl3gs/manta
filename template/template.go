package template

import (
	"time"

	"github.com/f1shl3gs/manta"
)

type ResourceType string

const (
	ResourceCheck     ResourceType = "check"
	ResourceConfig    ResourceType = "config"
	ResourceDashboard ResourceType = "dashboard"
	ResourceScrape    ResourceType = "scrape"
)

type ResourceItem struct {
	ID   manta.ID     `json:"id"`
	Type ResourceType `json:"type"`
	Name string       `json:"name"`
}

type Template struct {
	ID      manta.ID  `json:"ID"`
	Name    string    `json:"name"`
	Desc    string    `json:"desc,omitempty"`
	OrgID   manta.ID  `json:"orgID"`
	Created time.Time `json:"created"`

	Resources []ResourceItem `json:"resources"`
}
