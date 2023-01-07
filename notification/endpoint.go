package notification

import (
	"encoding/json"
	"fmt"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/errors"
)

func UnmarshalJSON(b []byte) (manta.NotificationEndpoint, error) {
	var raw struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}

	var endpoint manta.NotificationEndpoint
	switch raw.Type {
	case "http":
		endpoint = &HTTP{}
	default:
		return nil, &errors.Error{
			Code: errors.EInvalid,
			Msg:  fmt.Sprintf("invliad notification endpoint type %s", raw.Type),
		}
	}

	if err := json.Unmarshal(b, endpoint); err != nil {
		return nil, &errors.Error{
			Code: errors.EInvalid,
			Err:  err,
		}
	}

	return endpoint, nil
}
