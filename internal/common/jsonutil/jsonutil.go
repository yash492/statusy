package jsonutil

import (
	"encoding/json"
	"fmt"
)

// UnmarshalWithType parses raw JSON into a typed value T.
func UnmarshalWithType[T any](raw json.RawMessage) (T, error) {
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		return v, fmt.Errorf("failed to parse %T: %w", v, err)
	}
	return v, nil
}
