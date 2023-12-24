package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/pkg/schema"
)

type componentDBConn struct {
	db
}

func NewComponentDBConn() componentDBConn {
	return componentDBConn{
		db: dbConn,
	}
}

func (db componentDBConn) Create(components []schema.Component) ([]schema.Component, error) {
	batch := pgx.Batch{}
	var returnedComponents []schema.Component

	query := `INSERT INTO components 
			(name, service_id, provider_id)
			VALUES ($1, $2, $3)
			ON CONFLICT(name, service_id) DO NOTHING
			RETURNING *`

	for _, component := range components {
		batch.
			Queue(query, component.Name, component.ServiceID, component.ProviderID).
			Query(func(rows pgx.Rows) error {
				returnedComponent, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Component])
				if err != nil {
					return err
				}
				returnedComponents = append(returnedComponents, returnedComponent)
				return nil
			})
	}

	err := db.pgConn.SendBatch(context.Background(), &batch).Close()
	return returnedComponents, err
}

func (db componentDBConn) GetAllByServiceID(serviceID uint) ([]schema.Component, error) {
	query := "SELECT * FROM components WHERE service_id = $1"
	rows, err := db.pgConn.Query(context.Background(), query, serviceID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.Component])
}
