package incidentsdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/incidents"
)

//go:embed queries/get_incidents_by_service.sql
var getIncidentsByServiceQuery string

type incidentByServiceDto struct {
	ID                uint      `db:"id"`
	ServiceID         uint      `db:"service_id"`
	Title             string    `db:"title"`
	Status            string    `db:"status"`
	Link              string    `db:"link"`
	ProviderCreatedAt time.Time `db:"provider_created_at"`
	TotalCount        int64     `db:"total_count"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

func (c *PostgresIncidentRepository) GetByService(ctx context.Context, params incidents.IncidentByServiceParams) ([]incidents.IncidentByServiceResult, error) {
	compIDs := params.ComponentIDs
	if compIDs == nil {
		compIDs = []int{}
	}
	compGroupIDs := params.ComponentGroupIDs
	if compGroupIDs == nil {
		compGroupIDs = []int{}
	}
	hasFilter := len(compIDs) > 0 || len(compGroupIDs) > 0

	rows, err := c.readDB.Query(ctx, getIncidentsByServiceQuery, pgx.NamedArgs{
		"service_id":          params.ServiceID,
		"has_filter":          hasFilter,
		"component_ids":       compIDs,
		"component_group_ids": compGroupIDs,
		"limit":               params.Limit,
		"offset":              params.Offset,
	})
	if err != nil {
		c.lg.ErrorContext(ctx, "error querying incidents by service", slog.Any("service_id", params.ServiceID), slog.Any("err", err))

		return nil, apperrors.InternalError("failed to query incidents by service", err)
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[incidentByServiceDto])
	if err != nil {
		c.lg.ErrorContext(ctx, "error collecting incidents by service rows", slog.Any("service_id", params.ServiceID), slog.Any("err", err))
		return nil, apperrors.InternalError("failed to collect incidents by service rows", err)
	}

	results := make([]incidents.IncidentByServiceResult, 0, len(dtos))
	for _, item := range dtos {
		results = append(results, incidents.IncidentByServiceResult{
			ID:                item.ID,
			ServiceID:         item.ServiceID,
			Title:             item.Title,
			Status:            item.Status,
			ProviderCreatedAt: item.ProviderCreatedAt,
			Link:              item.Link,
			TotalCount:        item.TotalCount,
		})
	}

	return results, nil
}
