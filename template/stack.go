package template

type Resource map[string]interface{}

type Object struct {
	Kind string
	Name string
	Spec Resource
}
