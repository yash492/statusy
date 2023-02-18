package store

import (
	"backend/models"
)

func (d Db) AddIncidentUpdates([]models.IncidentUpdate) ([]models.IncidentUpdate, error) {
	return nil, nil
}

func (d Db) GetLastIncidentCreatedAtForSlug(slug string) (models.LastUpdatedIncidentForSlug, error) {

	var response models.LastUpdatedIncidentForSlug
	result := d.Table("incident_updates").
		Select("services.slug", "incident_updates.created_at").
		Joins("JOIN incidents ON incidents.id = incident_updates.incident_id").
		Joins("JOIN services ON incidents.service_id = services.id").
		Where("services.slug = ?", slug).
		Order("incident_updates.created_at DESC").
		Limit(1).
		Scan(&response)

	return response, result.Error
}
