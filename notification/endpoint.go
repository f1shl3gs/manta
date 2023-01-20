package notification

import (
	"encoding/json"
	"fmt"

	"github.com/f1shl3gs/manta"
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
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Msg:  fmt.Sprintf("invliad notification endpoint type %s", raw.Type),
		}
	}

	if err := json.Unmarshal(b, endpoint); err != nil {
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Err:  err,
		}
	}

	return endpoint, nil
}
