package store

import (
	"backend/models"

	"gorm.io/gorm/clause"
)

func (d Db) AddComponents(components []models.Component) ([]models.Component, error) {
	result := d.Clauses(clause.OnConflict{DoNothing: true}).Create(&components)
	return components, result.Error
}

func (d Db) GetComponentsByCodeAndService(code string, serviceId uint) (models.Component, error) {
	var component models.Component
	result := d.Joins("LEFT JOIN services ON components.service_id = services.id").
		Where("components.metadata ->> 'component_id' = ? AND services.id = ?", code, serviceId).
		First(&component)
	return component, result.Error
}
