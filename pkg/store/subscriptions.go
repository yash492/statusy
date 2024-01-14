package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/pkg/schema"
)

type subscriptionDBConn struct {
	db
}

func NewSubscriptionConn() subscriptionDBConn {
	return subscriptionDBConn{
		db: dbConn,
	}
}

func (db subscriptionDBConn) GetAllServicesForSubscriptions(serviceName string) ([]schema.ServicesForSubsciptions, error) {
	query := `SELECT services.id AS service_id, services.name AS service_name
				FROM services 
				WHERE services.name ILIKE '%' || $1 || '%'
				AND NOT EXISTS 
				(select * from subscriptions WHERE services.id = subscriptions.service_id AND subscriptions.deleted_at IS NULL)
				LIMIT 5`

	rows, err := db.pgConn.Query(context.Background(), query, serviceName)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.ServicesForSubsciptions])
}

func (db subscriptionDBConn) Create(serviceID uint, componentIDs []uint, isAllComponents bool) error {

	ctx := context.Background()
	tx, err := db.pgConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	insertSubscriptionQuery := `INSERT INTO subscriptions (service_id, is_all_components)
								VALUES($1, $2) RETURNING *`

	rows, err := tx.Query(context.Background(), insertSubscriptionQuery, serviceID, isAllComponents)
	if err != nil {
		return err
	}

	subscription, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Subscription])
	if err != nil {
		return err
	}

	if !isAllComponents {
		batch := &pgx.Batch{}
		for _, componentID := range componentIDs {
			batch.Queue(`INSERT INTO subscription_components (subscription_id, component_id)
							VALUES($1, $2)`, subscription.ID, componentID)
		}
		if err = tx.SendBatch(context.Background(), batch).Close(); err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (db subscriptionDBConn) GetWithComponents(subscriptionID uuid.UUID) ([]schema.SubscriptionWithComponents, error) {
	query := `SELECT
				subscriptions.uuid AS uuid,
				services.name AS service_name,
				components.id AS component_id,
				components.name AS component_name,
				services.id AS service_id,
				(
				CASE
					WHEN subscriptions.uuid IS NULL THEN 'false'::BOOLEAN
					WHEN subscriptions.uuid IS NOT NULL THEN 'true'::BOOLEAN
				END
				) AS is_configured
			FROM
				subscription_components
				JOIN subscriptions ON subscription_components.subscription_id = subscriptions.id
				RIGHT JOIN components ON subscription_components.component_id = components.id
				JOIN services ON components.service_id = services.id
			WHERE
				components.service_id = (
				SELECT
					service_id
				FROM
					subscriptions
				WHERE
					subscriptions.uuid = $1
				)
				AND subscriptions.deleted_at IS NULL
				AND subscription_components.deleted_at IS NULL`

	rows, err := db.pgConn.Query(context.Background(), query, subscriptionID.String())
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.SubscriptionWithComponents])
}

func (db subscriptionDBConn) Update(subscriptionID uuid.UUID, componentIDs []uint, isAllComponents bool) error {

	updateSubscriptionQuery := `UPDATE subscriptions SET is_all_components = $1, updated_at = $2 WHERE uuid = $3 RETURNING *`

	deleteSubscriptionComponentsQuery := `DELETE FROM subscription_components 
											WHERE subscription_id = $1`

	addSubscriptionComponentQuery := `INSERT INTO subscription_components (subscription_id, component_id) VALUES($1, $2)`

	ctx := context.Background()
	tx, err := db.pgConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, updateSubscriptionQuery, isAllComponents, time.Now(), subscriptionID.String())
	if err != nil {
		return err
	}

	subscription, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Subscription])
	if err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, deleteSubscriptionComponentsQuery, subscription.ID); err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, componentID := range componentIDs {
		batch.Queue(addSubscriptionComponentQuery, subscription.ID, componentID)
	}

	err = tx.SendBatch(ctx, batch).Close()
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (db subscriptionDBConn) GetForIncidentUpdates(incidentUpdateID uint) ([]schema.SubscriptionForIncidentUpdates, error) {
	query := `SELECT
				services.id AS service_id,
				services.name AS service_name,
				components.id AS component_id,
				components.name AS component_name,
				incidents.id AS incident_id,
				incidents.name AS incident_name,
				incidents.link AS incident_link,
				incidents.impact AS incident_impact,
				incident_updates.id AS incident_update_id,
				incident_updates.description AS incident_update,
				incident_updates.provider_status AS incident_update_provider_status,
				incident_updates.status AS incident_update_status,
				incident_updates.status_time::TIMESTAMPTZ AS incident_update_status_time,
				subscriptions.is_all_components AS is_all_components
			FROM
				incident_updates
				JOIN incidents ON incidents.id = incident_updates.incident_id
				JOIN incident_components ON incident_components.incident_id = incidents.id
				JOIN services ON services.id = incidents.service_id
				JOIN components ON incident_components.component_id = components.id
				JOIN subscriptions ON subscriptions.service_id = services.id
			WHERE
				(
				(subscriptions.is_all_components = true)
				OR (
					EXISTS (
					SELECT
						id
					FROM
						subscription_components
					WHERE
						subscription_components.subscription_id = subscriptions.id
						AND subscription_components.component_id IN (incident_components.component_id)
					)
				)
				)
				AND incident_updates.id = $1
				AND subscriptions.deleted_at IS NULL`

	rows, err := db.pgConn.Query(context.Background(), query, incidentUpdateID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.SubscriptionForIncidentUpdates])
}
