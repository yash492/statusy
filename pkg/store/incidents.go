package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yash492/statusy/pkg/schema"
)

type incidentDBConn struct {
	db
}

func NewIncidentDBConn() incidentDBConn {
	return incidentDBConn{
		db: dbConn,
	}
}

func (db incidentDBConn) Create(incidents []schema.Incident) ([]schema.Incident, error) {
	batch := pgx.Batch{}

	var returningIncidents []schema.Incident

	query := `INSERT INTO incidents
				(name, link, service_id, provider_id) 
				VALUES ($1, $2, $3, $4)
				ON CONFLICT DO NOTHING
				RETURNING *`

	for _, incident := range incidents {
		batch.
			Queue(query,
				incident.Name,
				incident.Link,
				incident.ServiceID,
				incident.ProviderID,
			).
			Query(func(rows pgx.Rows) error {
				queriedIncident, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Incident])
				if err != nil {
					return err
				}
				returningIncidents = append(returningIncidents, queriedIncident)
				return nil
			})
	}

	err := db.pgConn.SendBatch(context.Background(), &batch).Close()
	return returningIncidents, err

}

func (db incidentDBConn) GetByProviderID(providerID string) (schema.Incident, error) {
	query := "SELECT * FROM incidents WHERE provider_id = $1 AND deleted_at IS NULL"
	rows, err := db.pgConn.Query(context.Background(), query, providerID)
	if err != nil {
		return schema.Incident{}, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Incident])
}

func (db incidentDBConn) CreateIncidentUpdates(incidentUpdates []schema.IncidentUpdate) error {
	query := `INSERT INTO incident_updates
				(incident_id, description, status_time, status, provider_id) 
				VALUES ($1, $2, $3, $4, $5)`

	batch := pgx.Batch{}

	for _, update := range incidentUpdates {
		batch.Queue(
			query,
			update.IncidentID,
			update.Description,
			update.StatusTime,
			update.Status,
			update.ProviderID,
		)
	}

	err := dbConn.pgConn.SendBatch(context.Background(), &batch).Close()
	return err
}

func (db incidentDBConn) GetLastIncidentUpdatesTimeByService(serviceID uint, incidentIDs []uint) ([]schema.LastIncidentUpdateForIncident, error) {
	query := `SELECT MAX(status_time) AS last_incident_updates_time, incidents.id AS incident_id 
				FROM incident_updates
				RIGHT JOIN incidents ON incidents.id = incident_updates.incident_id
				WHERE incidents.service_id = $1
				AND incidents.id = ANY($2)
				GROUP BY incidents.id`

	incidentIDsPgType := pgtype.FlatArray[uint](incidentIDs)
	rows, err := db.pgConn.Query(context.Background(), query, serviceID, incidentIDsPgType)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[schema.LastIncidentUpdateForIncident])
}

func (d incidentDBConn) CreateIncidentComponents(incidentComponents []schema.IncidentComponent) error {
	query := `INSERT INTO incident_components
				(incident_id, component_id) 
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING`

	batch := pgx.Batch{}

	for _, components := range incidentComponents {
		batch.Queue(
			query,
			components.IncidentID,
			components.ComponentID,
		)
	}

	err := dbConn.pgConn.SendBatch(context.Background(), &batch).Close()
	return err
}
