package notification

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/f1shl3gs/manta"
)

const (
	httpTokenSuffix    = "-token"
	httpUsernameSuffix = "-username"
	httpPasswordSuffix = "-password"
)

type HTTP struct {
	Base

	URL        string            `json:"url"`
	Headers    map[string]string `json:"headers"`
	Username   manta.SecretField `json:"username"`
	Password   manta.SecretField `json:"password"`
	Token      manta.SecretField `json:"token"`
	Method     string            `json:"method"`
	AuthMethod string            `json:"authMethod"`

	ContentTemplate string `json:"contentTemplate"`
}

// BackfillSecretKeys fill back fill the secret field key during the unmarshaling
// if value of that secret field is not nil
func (h *HTTP) BackfillSecretKeys() {
	if h.Token.Key == "" && h.Token.Value != nil {
		h.Token.Key = h.ID.String() + httpTokenSuffix
	}

	if h.Username.Key == "" && h.Username.Value != nil {
		h.Username.Key = h.ID.String() + httpUsernameSuffix
	}

	if h.Password.Key == "" && h.Password.Value != nil {
		h.Password.Key = h.ID.String() + httpPasswordSuffix
	}
}

// SecretFields return available secret fields
func (h HTTP) SecretFields() []manta.SecretField {
	arr := make([]manta.SecretField, 0)
	if h.Token.Key != "" {
		arr = append(arr, h.Token)
	}

	if h.Username.Key != "" {
		arr = append(arr, h.Username)
	}

	if h.Password.Key != "" {
		arr = append(arr, h.Password)
	}

	return arr
}

var (
	validMethods = map[string]bool{
		http.MethodGet:  true,
		http.MethodPost: true,
		http.MethodPut:  true,
	}

	validAuthMethods = map[string]bool{
		"none":   true,
		"basic":  true,
		"bearer": true,
	}
)

func (h HTTP) Valid() error {
	if err := h.Base.Valid(); err != nil {
		return err
	}

	if h.URL == "" {
		return &manta.Error{
			Code: manta.EInvalid,
			Msg:  "http endpoint URL is empty",
		}
	}

	if _, err := url.Parse(h.URL); err != nil {
		return &manta.Error{
			Code: manta.EInvalid,
			Msg:  "http endpoint URL is invalid",
			Err:  err,
		}
	}

	if !validMethods[h.Method] {
		return &manta.Error{
			Code: manta.EInvalid,
			Msg:  "invalid http method",
		}
	}

	if !validAuthMethods[h.AuthMethod] {
		return &manta.Error{
			Code: manta.EInvalid,
			Msg:  "invalid http auth method",
		}
	}

	if h.AuthMethod == "basic" && (h.Username.Key == "" || h.Password.Key == "") {
		return &manta.Error{
			Code: manta.EInvalid,
			Msg:  "invalid http username/password for basic auth",
		}
	}

	if h.AuthMethod == "baerer" && h.Token.Key == "" {
		return &manta.Error{
			Code: manta.EInvalid,
			Msg:  "invalid http token for bearer auth",
		}
	}

	return nil
}

// MarshalJSON implent json.Marshaler
func (h HTTP) MarshalJSON() ([]byte, error) {
	type httpAlias HTTP

	return json.Marshal(
		struct {
			httpAlias
			Type string `json:"type"`
		}{
			httpAlias: httpAlias(h),
			Type:      h.Type(),
		})
}

// Type implement manta.NotificationEndpoint
func (h HTTP) Type() string {
	return "http"
}

// ParseResponse will parse the http response from http
func (h HTTP) ParseResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return &manta.Error{
			Msg: string(body),
		}
	}

	return nil
}
