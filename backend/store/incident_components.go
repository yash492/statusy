package store

import (
	"backend/models"
)

func (d Db) AddIncidentComponents(incidentComponents []models.IncidentComponent) ([]models.IncidentComponent, error) {
	result := d.Create(&incidentComponents)
	return incidentComponents, result.Error
}
