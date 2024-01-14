package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/pkg/schema"
)

type squadcastExtensionDBConn struct {
	db
}

func NewSquadcastExtensionConn() squadcastExtensionDBConn {
	return squadcastExtensionDBConn{
		db: dbConn,
	}
}

func (db squadcastExtensionDBConn) Save(webhookURL string, uuid uuid.UUID) error {
	query := `INSERT INTO squadcast_extensions(webhook_url, uuid) 
	VALUES($1, $2) ON CONFLICT(uuid) 
	DO UPDATE SET webhook_url=EXCLUDED.webhook_url, updated_at=$3`

	_, err := db.pgConn.Exec(context.Background(), query, webhookURL, uuid.String(), time.Now())
	return err
}

func (db squadcastExtensionDBConn) Delete(uuid uuid.UUID) error {
	query := `UPDATE squadcast_extensions SET deleled_at = $1 WHERE uuid = $2`
	_, err := db.pgConn.Exec(context.Background(), query, time.Now(), uuid.String())
	return err
}

func (db squadcastExtensionDBConn) Get() (schema.SquadcastExtension, error) {
	query := `SELECT * FROM squadcast_extensions WHERE deleted_at IS NULL`
	rows, err := db.pgConn.Query(context.Background(), query)
	if err != nil {
		return schema.SquadcastExtension{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.SquadcastExtension])
}
