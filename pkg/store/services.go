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

	query := `INSERT INTO services 
				(name, link, slug, provider_type, should_scrap_website, incident_url, schedule_maintenance_url, components_url)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				ON CONFLICT(slug) DO NOTHING`

	for _, service := range services {
		batch.Queue(
			query,
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

func (db serviceDBConn) GetAll() ([]schema.Service, error) {
	sqlQuery := `SELECT * from services WHERE deleted_at IS NULL`
	rows, err := db.pgConn.Query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	services, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.Service])
	if err != nil {
		return nil, err
	}

	return services, nil
}
