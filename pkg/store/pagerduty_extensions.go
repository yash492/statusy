package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/pkg/schema"
)

type pagedutyExtensionDBConn struct {
	db
}

func NewPagerdutyExtensionConn() pagedutyExtensionDBConn {
	return pagedutyExtensionDBConn{
		db: dbConn,
	}
}

func (db pagedutyExtensionDBConn) Save(routingKey string, uuid uuid.UUID) error {
	query := `INSERT INTO pagerduty_extensions(routing_key, uuid) 
				VALUES($1, $2) ON CONFLICT(uuid) 
				DO UPDATE SET routing_key=EXCLUDED.routing_key, updated_at=$3`

	_, err := db.pgConn.Exec(context.Background(), query, routingKey, uuid.String(), time.Now())
	return err
}

func (db pagedutyExtensionDBConn) Delete(uuid uuid.UUID) error {
	query := `UPDATE pagerduty_extensions SET deleled_at = $1 WHERE uuid = $2`
	_, err := db.pgConn.Exec(context.Background(), query, time.Now(), uuid.String())
	return err
}

func (db pagedutyExtensionDBConn) Get() (schema.PagerdutyExtension, error) {
	query := `SELECT * FROM pagerduty_extensions WHERE deleted_at IS NULL`
	rows, err := db.pgConn.Query(context.Background(), query)
	if err != nil {
		return schema.PagerdutyExtension{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.PagerdutyExtension])
}
