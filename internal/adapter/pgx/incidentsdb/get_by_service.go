package incidentsdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/incidents"
)

//go:embed queries/get_incidents_by_service.sql
var getIncidentsByServiceQuery string

type incidentByServiceDto struct {
	ID                uint      `db:"id"`
	ServiceID         uint      `db:"service_id"`
	Title             string    `db:"title"`
	Status            string    `db:"status"`
	ProviderCreatedAt time.Time `db:"provider_created_at"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

func (c *PostgresIncidentRepository) GetByService(ctx context.Context, params incidents.IncidentByServiceParams) ([]incidents.IncidentByServiceResult, error) {
	rows, err := c.readDB.Query(ctx, getIncidentsByServiceQuery, pgx.NamedArgs{
		"service_id": params.ServiceID,
		"limit":      params.Limit,
		"offset":     params.Offset,
	})
	if err != nil {
		c.lg.ErrorContext(ctx, "error querying incidents by service", slog.Any("service_id", params.ServiceID), slog.Any("err", err))

		return nil, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[incidentByServiceDto])
	if err != nil {
		if c.lg != nil {
			c.lg.ErrorContext(ctx, "error collecting incidents by service rows", slog.Any("service_id", params.ServiceID), slog.Any("err", err))
		}
		return nil, err
	}

	results := make([]incidents.IncidentByServiceResult, 0, len(dtos))
	for _, item := range dtos {
		results = append(results, incidents.IncidentByServiceResult{
			ID:                item.ID,
			ServiceID:         item.ServiceID,
			Title:             item.Title,
			Status:            item.Status,
			ProviderCreatedAt: item.ProviderCreatedAt,
		})
	}

	return results, nil
}
