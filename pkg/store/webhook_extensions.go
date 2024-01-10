package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/pkg/schema"
)

type webhookExtensionDBConn struct {
	db
}

func NewWebhookExtensionConn() webhookExtensionDBConn {
	return webhookExtensionDBConn{
		db: dbConn,
	}
}

func (db webhookExtensionDBConn) Save(webhookURL string, secret *string, uuid uuid.UUID) error {
	query := `INSERT INTO chatops_extensions(secret, webhook_url, uuid) 
				VALUES($1, $2, $3) ON CONFLICT(uuid) 
				DO UPDATE SET type=EXCLUDED.secret, webhook_url=EXCLUDED.webhook_url`

	_, err := db.pgConn.Exec(context.Background(), query, secret, webhookURL, uuid.String())
	return err
}

func (db webhookExtensionDBConn) Delete(uuid uuid.UUID) error {
	query := `UPDATE webhook_extensions SET deleled_at = $1 WHERE uuid = $2`
	_, err := db.pgConn.Exec(context.Background(), query, time.Now(), uuid.String())
	return err
}

func (db webhookExtensionDBConn) Get() (schema.WebhookExtension, error) {
	query := `SELECT * FROM webhook_extensions WHERE deleted_at IS NULL`
	rows, err := db.pgConn.Query(context.Background(), query)
	if err != nil {
		return schema.WebhookExtension{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.WebhookExtension])
}
