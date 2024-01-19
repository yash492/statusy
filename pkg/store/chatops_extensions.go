package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/pkg/schema"
)

type chatOpsExtensionDBConn struct {
	db
}

func NewChatOpsExtensionConn() chatOpsExtensionDBConn {
	return chatOpsExtensionDBConn{
		db: dbConn,
	}
}

func (db chatOpsExtensionDBConn) Save(chatOpsType string, webhookURL string, uuid uuid.UUID) error {
	query := `INSERT INTO chatops_extensions(type, webhook_url, uuid) 
				VALUES($1, $2, $3) ON CONFLICT(uuid) 
				DO UPDATE SET type=EXCLUDED.type, webhook_url=EXCLUDED.webhook_url, updated_at=$4`

	_, err := db.pgConn.Exec(context.Background(), query, chatOpsType, webhookURL, uuid.String(), time.Now())
	return err
}

func (db chatOpsExtensionDBConn) Delete(uuid uuid.UUID) error {
	query := `UPDATE chatops_extensions SET deleted_at = $1 WHERE uuid = $2`
	_, err := db.pgConn.Exec(context.Background(), query, time.Now(), uuid.String())
	return err
}

func (db chatOpsExtensionDBConn) Get() ([]schema.ChatopsExtension, error) {
	query := `SELECT * FROM chatops_extensions WHERE deleted_at IS NULL`
	rows, err := db.pgConn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.ChatopsExtension])
}

func (db chatOpsExtensionDBConn) GetByType(chatopType string) (schema.ChatopsExtension, error) {
	query := `SELECT * FROM chatops_extensions WHERE type = $1 AND deleted_at IS NULL`
	rows, err := db.pgConn.Query(context.Background(), query, chatopType)
	if err != nil {
		return schema.ChatopsExtension{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.ChatopsExtension])
}
