package template

import "github.com/f1shl3gs/manta"

var (
	SupportedResources = []string{
		"check",
		"dashboard",
		"notification_endpoint",
		"otcl",
		"secret",
		"variables",
		"scrape",
	}
)

type Resource map[string]interface{}

type Object struct {
	Kind string
	Name string
	Spec Resource
}

type Backend interface {
	manta.CheckService
	manta.DashboardService
	manta.NotificationEndpointService
	manta.OtclService
	manta.SecretService
	manta.VariableService
	manta.ScraperTargetService
}

func Apply(backend Backend, source string) error {

}
