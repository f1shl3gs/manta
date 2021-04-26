package template

var (
	SupportedResources = []string{
		"check",
		"dashboard",
		"notification_endpoint",
		"otcl",
		"secret",
		"variables",
	}
)

type Resource map[string]interface{}

type Object struct {
	Kind string
	Name string
	Spec Resource
}
