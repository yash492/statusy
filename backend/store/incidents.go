package store

import (
	"backend/models"
)

func (d Db) AddIncidents(incidents []models.Incident) ([]models.Incident, error) {
	result := d.Create(&incidents)
	return incidents, result.Error
}

// func (d *Db)
