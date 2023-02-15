package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	Uuid uuid.UUID
}

func (b *BaseModel) BeforeCreate(db *gorm.DB) error {
	b.Uuid = uuid.New()
	return nil
}
