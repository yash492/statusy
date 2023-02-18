package store

import (
	"backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Db struct {
	*gorm.DB
}

func (d Db) AddServices(services []models.Service) ([]models.Service, error) {
	result := d.Clauses(clause.OnConflict{DoNothing: true}).Create(&services)
	return services, result.Error
}

func (d Db) GetAllServices() ([]models.Service, error) {
	var services []models.Service
	result := d.Find(&services)
	return services, result.Error
}

func (d Db) GetServiceBySlug(slug string) (models.Service, error) {
	var service models.Service
	result := d.Where("slug = ?", slug).First(&service)
	return service, result.Error
}
