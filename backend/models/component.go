package models

type Component struct {
	BaseModel
	Name      string
	Slug      string
	ServiceId uint
}
