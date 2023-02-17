package store

import "backend/types"

func (d Db) AddIncidents(incidents []types.Incident) ([]types.Incident, error) {
	result := d.Create(&incidents)
	return incidents, result.Error
}

// func (d *Db)
