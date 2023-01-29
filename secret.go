package manta

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

type SecretField struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
}

func (s SecretField) String() string {
	if s.Key == "" {
		return ""
	}

	return "secret: " + s.Key
}

func (s SecretField) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *SecretField) UnmarshalJSON(b []byte) error {
	var ss string
	if err := json.Unmarshal(b, &ss); err != nil {
		return err
	}

	if ss == "" {
		s.Key = ""
		return nil
	}

	if strings.HasPrefix(ss, "secret: ") {
		s.Key = ss[len("secret: "):]
	} else {
		s.Value = strPtr(ss)
	}

	return nil
}

func strPtr(s string) *string {
	ss := new(string)
	*ss = s
	return ss
}

type Secret struct {
	Key     string    `json:"key"`
	Updated time.Time `json:"updated"`
	OrgID   ID        `json:"orgID"`
	Value   string    `json:"value,omitempty"`
}

func (s *Secret) Desensitize() {
	s.Value = ""
}

type SecretService interface {
	// LoadSecret retrieves the secret value v found at key k for organization orgID
	LoadSecret(ctx context.Context, orgID ID, k string) (*Secret, error)

	// GetSecrets retrieves desensitized secrets of 'orgID'
	GetSecrets(ctx context.Context, orgID ID) ([]Secret, error)

	// PutSecret creates or updates a secret and return the desensitized secret
	PutSecret(ctx context.Context, secret *Secret) (*Secret, error)

	// DeleteSecret deletes secrets by keys
	DeleteSecret(ctx context.Context, orgID ID, keys ...string) error
}
