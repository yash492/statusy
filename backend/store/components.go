package store

import (
	"backend/types"

	"gorm.io/gorm/clause"
)

func (d Db) AddComponents(components []types.Component) ([]types.Component, error) {
	result := d.Clauses(clause.OnConflict{DoNothing: true}).Create(&components)
	return components, result.Error
}
