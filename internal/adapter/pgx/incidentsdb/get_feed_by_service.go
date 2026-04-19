package incidentsdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/incidents"
)

//go:embed queries/get_feed_incidents_by_service.sql
var getFeedIncidentsByServiceQuery string

type feedIncidentByServiceDto struct {
	ID                 uint      `db:"id"`
	ServiceID          uint      `db:"service_id"`
	Title              string    `db:"title"`
	Status             string    `db:"status"`
	Link               string    `db:"link"`
	ProviderCreatedAt  time.Time `db:"provider_created_at"`
	AffectedComponents string    `db:"affected_components"`
}

func (c *PostgresIncidentRepository) GetFeedByService(ctx context.Context, params incidents.IncidentByServiceParams) ([]incidents.FeedIncidentByServiceResult, error) {
	rows, err := c.readDB.Query(ctx, getFeedIncidentsByServiceQuery, pgx.NamedArgs{
		"service_id": params.ServiceID,
		"limit":      params.Limit,
		"offset":     params.Offset,
	})
	if err != nil {
		c.lg.ErrorContext(ctx, "error querying feed incidents by service", slog.Any("service_id", params.ServiceID), slog.Any("err", err))

		return nil, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[feedIncidentByServiceDto])
	if err != nil {
		if c.lg != nil {
			c.lg.ErrorContext(ctx, "error collecting feed incidents by service rows", slog.Any("service_id", params.ServiceID), slog.Any("err", err))
		}
		return nil, err
	}

	results := make([]incidents.FeedIncidentByServiceResult, 0, len(dtos))
	for _, item := range dtos {
		results = append(results, incidents.FeedIncidentByServiceResult{
			ID:                 item.ID,
			ServiceID:          item.ServiceID,
			Title:              item.Title,
			Status:             item.Status,
			ProviderCreatedAt:  item.ProviderCreatedAt,
			Link:               item.Link,
			AffectedComponents: item.AffectedComponents,
		})
	}

	return results, nil
}
