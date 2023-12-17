package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/pkg/schema"
)

type serviceDBConn struct {
	db
}

func NewServiceDBConn() serviceDBConn {
	return serviceDBConn{
		db: dbConn,
	}
}

func (db serviceDBConn) Create(services []schema.Service) error {
	batch := pgx.Batch{}

	sqlQuery := `INSERT INTO services 
				(name, link, slug, provider_type, should_scrap_website, incident_url, schedule_maintenance_url, components_url)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				ON CONFLICT(slug) DO NOTHING`

	for _, service := range services {
		batch.Queue(
			sqlQuery,
			service.Name,
			service.Link,
			service.Slug,
			service.ProviderType,
			service.ShouldScrapWebsite,
			service.IncidentURL,
			service.ScheduleMaintenanceURL,
			service.ComponentsURL,
		)
	}

	err := db.pgConn.SendBatch(context.Background(), &batch).Close()
	return err
}
