package marshaling

import (
	"encoding/json"
)

func JSON() Marshaling {
	return New(MarshalerFunc(json.Marshal), UnmarshalerFunc(json.Unmarshal))
}
