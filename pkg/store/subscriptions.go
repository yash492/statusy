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

func (db subscriptionDBConn) GetAllServicesForSubscriptions(serviceName string) ([]schema.ServicesForSubsciption, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.ServicesForSubsciption])
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

func (db subscriptionDBConn) GetWithComponents(subscriptionID uuid.UUID) ([]schema.SubscriptionWithComponent, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.SubscriptionWithComponent])
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

func (db subscriptionDBConn) GetForIncidentUpdates(incidentUpdateID uint) ([]schema.SubscriptionForIncidentUpdate, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.SubscriptionForIncidentUpdate])
}

func (db subscriptionDBConn) DashboardSubscription(serviceName string, offset, limit uint) ([]schema.DashboardSubscription, error) {
	query := `WITH
				incidents_cte AS (
				SELECT
					*
				FROM
					incidents
				WHERE
					NOT EXISTS (
					select
						*
					from
						incident_updates
					WHERE
						status = 'resolved'
						AND incidents.id = incident_updates.incident_id
					)
				),
				subscriptions_cte AS (
				SELECT
					subscriptions.id AS subscription_id,
					incidents_cte.id AS incident_id,
					subscriptions.uuid AS subscription_uuid,
					incidents_cte.name AS incident_name,
					incidents_cte.link AS incident_link,
					incidents_cte.impact AS incident_impact,
					RANK() OVER (
					PARTITION BY
						incidents_cte.service_id
					ORDER BY
						incidents_cte.provider_created_at DESC
					) AS rank
				FROM
					incidents_cte
					JOIN subscriptions ON subscriptions.service_id = incidents_cte.service_id
				WHERE
					subscriptions.is_all_components = true
					OR (
					EXISTS (
						SELECT
						subscription_components.id
						FROM
						subscription_components
						JOIN incident_components ON subscription_components.component_id = incident_components.component_id
						WHERE
						subscription_components.subscription_id = subscriptions.id
						AND incident_components.incident_id = incidents_cte.id
					)
					)
				),
				subscriptions_list_cte AS (
				SELECT
					*
				FROM
					subscriptions_cte
				WHERE
					rank < 2
				)
			SELECT
				COUNT(subscriptions.id) OVER () AS subscriptions_count,
				subscriptions_list_cte.incident_id,
				subscriptions.uuid AS subscription_uuid,
				subscriptions_list_cte.incident_name,
				subscriptions_list_cte.incident_link,
				subscriptions_list_cte.incident_impact,
				services.name AS service_name,
				services.id AS service_id,
				(
				CASE
					WHEN subscriptions_list_cte.incident_id IS NULL THEN 'false'::BOOLEAN
					WHEN subscriptions_list_cte.incident_id IS NOT NULL THEN 'true'::BOOLEAN
				END
				) AS is_down
			FROM
				subscriptions_list_cte
				RIGHT JOIN subscriptions ON subscriptions.id = subscriptions_list_cte.subscription_id
				JOIN services ON services.id = subscriptions.service_id
			WHERE
				subscriptions.deleted_at IS NULL
				AND services.name ILIKE '%' || $1 || '%'
			ORDER BY 
				is_down DESC
			OFFSET
				$2
			LIMIT
				$3`

	rows, err := db.pgConn.Query(context.Background(), query, serviceName, offset, limit)
	if err != nil {
		return nil, nil
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.DashboardSubscription])

}

func (db subscriptionDBConn) GetIncidentsForSubscription(subscriptionUUID uuid.UUID, offset, limit uint) ([]schema.SubscriptionIncident, error) {
	query := `WITH
				incidents_cte AS (
				SELECT
					incidents.id AS incident_id,
					incident_updates.status_time AS last_updated_status_time,
					incident_updates.provider_status AS incident_status,
					incident_updates.status AS incident_normalised_status,
					incidents.provider_created_at AS incident_created_at,
					incidents.name AS incident_name,
					incidents.link AS incident_link,
					incidents.service_id AS service_id,
					RANK() OVER (
					PARTITION BY
						incidents.id
					ORDER BY
						incident_updates.status_time DESC
					) AS rank
				FROM
					incidents
					JOIN incident_updates ON incidents.id = incident_updates.incident_id
				),
				incident_list_cte AS (
				SELECT
					COUNT(incident_id) OVER () AS total_count,
					incident_id,
					last_updated_status_time,
					incident_status,
					incident_normalised_status,
					incident_created_at,
					incident_name,
					incident_link,
					subscriptions.id AS subscription_id
				FROM
					incidents_cte
					JOIN subscriptions ON subscriptions.service_id = incidents_cte.service_id
				WHERE
					rank < 2
					AND (
					subscriptions.is_all_components = true
					OR (
						EXISTS (
						SELECT
							subscription_components.id
						FROM
							subscription_components
							JOIN incident_components ON subscription_components.component_id = incident_components.component_id
						WHERE
							subscription_components.subscription_id = subscriptions.id
							AND incident_components.incident_id = incidents_cte.incident_id
						)
					)
					)
				)
			SELECT
				COUNT(incident_id) OVER () AS total_count,
				incident_id,
				last_updated_status_time,
				incident_status,
				incident_normalised_status,
				incident_created_at,
				incident_name,
				incident_link,
				services.name AS service_name,
				services.id AS service_id,
				subscriptions.is_all_components AS is_all_components_configured
			FROM
				incident_list_cte
				RIGHT JOIN subscriptions ON subscriptions.id = incident_list_cte.subscription_id
				JOIN services ON services.id = subscriptions.service_id
			WHERE
				subscriptions.uuid = $1
			ORDER BY
				incident_created_at DESC
			OFFSET
				$2
			LIMIT
				$3`

	rows, err := db.pgConn.Query(context.Background(), query, subscriptionUUID.String(), offset, limit)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.SubscriptionIncident])
}

func (d subscriptionDBConn) GetByID(subscriptionUUID uuid.UUID) (schema.SubscriptionWithService, error) {
	query := `SELECT
				subscriptions.uuid AS uuid,
				services.name AS service_name,
				services.id AS service_id
			FROM
				subscriptions
				JOIN services ON services.id = subscriptions.service_id
			WHERE
				uuid = $1
			AND 
				subscriptions.deleted_at IS NULL`

	rows, err := d.pgConn.Query(context.Background(), query, subscriptionUUID.String())
	if err != nil {
		return schema.SubscriptionWithService{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.SubscriptionWithService])
}
