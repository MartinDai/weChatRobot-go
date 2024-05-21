package util

import "encoding/json"

func ToJsonString(v any) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
