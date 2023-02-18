package store

import (
	"backend/models"

	"gorm.io/gorm/clause"
)

func (d Db) AddComponents(components []models.Component) ([]models.Component, error) {
	result := d.Clauses(clause.OnConflict{DoNothing: true}).Create(&components)
	return components, result.Error
}
