package models

import "encoding/json"

type Component struct {
	BaseModel
	Name      string
	Slug      string
	ServiceId uint
	Metadata  json.RawMessage
}
