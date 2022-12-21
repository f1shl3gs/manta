package manta

import (
	"context"
	"strconv"
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" {
		*d = 0
		return nil
	}

	s, err := strconv.Unquote(s)
	if err != nil {
		return err
	}

	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = Duration(parsed)
	return nil
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Duration(*d).String() + "\""), nil
}

type ComponentResource struct {
	Name          string `json:"name"`
	ComponentType string `json:"component_type"`
	Address       string `json:"address"`
	Port          uint16 `json:"port"`
}

type Instance struct {
	Uuid    string    `json:"uuid"`
	Created time.Time `json:"created"`
	Uptime  time.Time `json:"uptime"`

	Hostname  string              `json:"hostname"`
	Address   string              `json:"address"`
	Version   string              `json:"version"`
	Lease     Duration            `json:"lease"`
	Os        string              `json:"os"`
	Kernel    string              `json:"kernel,omitempty"`
	Tags      map[string]string   `json:"tags"`
	Resources []ComponentResource `json:"resources"`
}

func (ins *Instance) ExpiredAt(ts time.Time) bool {
	return ts.Sub(ins.Created) > time.Duration(ins.Lease)
}

type RegistryService interface {
	Register(ctx context.Context, ins *Instance) error

	Catalog(ctx context.Context) ([]*Instance, error)
}
