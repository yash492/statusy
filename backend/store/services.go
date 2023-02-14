package store

import (
	"backend/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Db struct {
	*gorm.DB
}

func (d Db) AddServices(services []types.Service) ([]types.Service, error) {
	result := d.Clauses(clause.OnConflict{DoNothing: true}).Create(&services)
	return services, result.Error
}

func (d Db) GetAllServices() ([]types.Service, error) {
	var services []types.Service
	result := d.Find(&services)
	return services, result.Error
}
