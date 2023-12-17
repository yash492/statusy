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

func (db componentDBConn) Create(components []schema.Component) error {
	batch := pgx.Batch{}

	query := `INSERT INTO components 
			(name, service_id, provider_component_id)
			VALUES ($1, $2, $3)
			ON CONFLICT(name, service_id) DO NOTHING`

	for _, component := range components {
		batch.Queue(query, component.Name, component.ServiceId, component.ProviderComponentId)
	}

	err := db.pgConn.SendBatch(context.Background(), &batch).Close()
	return err
}
