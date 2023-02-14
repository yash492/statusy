package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Component struct {
	gorm.Model
	Name      string
	Uuid      uuid.UUID
	Slug      string
	ServiceId uint
}
