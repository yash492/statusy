package store

import (
	"backend/models"
)

func (d Db) AddIncident(incident models.Incident) (models.Incident, error) {
	result := d.Create(&incident)
	return incident, result.Error
}

// func (d *Db)
