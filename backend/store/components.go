package store

import (
	"backend/models"

	"gorm.io/gorm/clause"
)

func (d Db) AddComponents(components []models.Component) ([]models.Component, error) {
	result := d.Clauses(clause.OnConflict{DoNothing: true}).Create(&components)
	return components, result.Error
}

func (d Db) GetComponentsBySlugAndService(slug string, serviceId uint) (models.Component, error) {
	var component models.Component
	result := d.Joins("LEFT JOIN services ON components.service_id = services.id").
		Where("components.slug = ? AND services.id = ?", slug, serviceId).
		First(&component)
	return component, result.Error
}
