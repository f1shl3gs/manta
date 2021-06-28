package template

import (
	"encoding/json"
	"io"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/slices"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	SupportedResources = []string{
		"Check",
		"Dashboard",
		"NotificationEndpoint",
		"Otcl",
		"Secret",
		"Variables",
		"Scrape",
	}
)

type Object struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Name       string `json:"name" yaml:"name"`

	Spec interface{}
}

func (obj *Object) UnmarshalJSON(bytes []byte) error {
	tempObj := struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Name       string `json:"name"`
	}{}

	err := json.Unmarshal(bytes, &tempObj)
	if err != nil {
		return err
	}

	switch tempObj.Kind {
	case "Check":
		obj.Spec = &manta.Check{}
	case "Dashboard":
		obj.Spec = &manta.Dashboard{}
	case "Scraper":
		obj.Spec = &manta.ScrapeTarget{}
	default:
		return errors.Errorf("unsupported kind %s", obj.Kind)
	}

	return json.Unmarshal(bytes, obj)
}

func (obj *Object) validateKind() error {
	if slices.Contain(SupportedResources, obj.Kind) {
		return nil
	}

	return nil
}

func (obj *Object) UnmarshalYAML(node *yaml.Node) error {
	var (
		specNode *yaml.Node
		err      error
	)

	for i := 0; i < len(node.Content); i++ {
		if node.Content[i].Value == "kind" {
			obj.Kind = node.Content[i+1].Value
			i += 1
			continue
		}

		if node.Content[i].Value == "apiVersion" {
			obj.APIVersion = node.Content[i+1].Value
			i += 1
			continue
		}

		if node.Content[i].Value == "name" {
			obj.Name = node.Content[i+1].Value
			i += 1
			continue
		}

		if node.Content[i].Value == "spec" {
			specNode = node.Content[i+1]
		}
	}

	if specNode == nil {
		return errors.New("spec node is expected")
	}

	switch obj.Kind {
	case "Check":
		obj.Spec = &manta.Check{}
	case "Dashboard":
		obj.Spec = &manta.Dashboard{}
	case "Scraper":
		obj.Spec = &manta.ScrapeTarget{}
	case "Otcl":
		obj.Spec = &manta.Otcl{}
	case "Secret":
		obj.Spec = &manta.Secret{}
	case "Variable":
		obj.Spec = &manta.Variable{}
	default:
		return errors.Errorf("unknown kind %s", obj.Kind)
	}

	err = specNode.Decode(obj.Spec)
	if err != nil {
		return err
	}

	return nil
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

func ApplyYAML(backend Backend, r io.Reader, orgID manta.ID) error {
	objs := make([]*Object, 0)

	dec := yaml.NewDecoder(r)
	for {
		obj := &Object{}
		err := dec.Decode(obj)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		objs = append(objs, obj)
	}

	return applyObjects(backend, orgID, objs)
}

func ApplyJSON(backend Backend, r io.Reader) error {
	return nil
}

// applyObjects apply the objs to the backend one by one
// 1. if the resource is missing, then create it
// 2. if the resource exist, then overwrite it ?
// 3. if the associate resource is deleted, create it again?
func applyObjects(backend Backend, orgID manta.ID, objs []*Object) error {
	return nil
}
