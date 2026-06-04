package incidentsdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/incidents"
)

//go:embed queries/insert_incidents.sql
var insertIncidentsQuery string

type incidentDto struct {
	ID                uint             `db:"id"`
	Title             string           `db:"title"`
	Link              string           `db:"link"`
	ProviderImpact    pgtype.Text      `db:"provider_impact"`
	Impact            pgtype.Text      `db:"impact"`
	ServiceID         uint             `db:"service_id"`
	ProviderID        string           `db:"provider_id"`
	ProviderCreatedAt time.Time        `db:"provider_created_at"`
	IsResolved        bool             `db:"is_resolved"`
	CreatedAt         time.Time        `db:"created_at"`
	UpdatedAt         time.Time        `db:"updated_at"`
	DeletedAt         pgtype.Timestamp `db:"deleted_at"`
}

func (c *PostgresIncidentRepository) SaveAll(ctx context.Context, params []incidents.IncidentParams) ([]incidents.IncidentResult, error) {
	batchInserts := &pgx.Batch{}
	incidentResponse := []incidentDto{}

	for _, param := range params {
		provImpStr, provImpOk := param.ProviderImpact.Get()
		impactStr, impactOk := param.Impact.Get()

		queryArgs := pgx.NamedArgs{
			"title": param.Title,
			"link":  param.Link,
			"provider_impact": pgtype.Text{
				String: provImpStr,
				Valid:  provImpOk,
			},
			"impact": pgtype.Text{
				String: impactStr,
				Valid:  impactOk,
			},
			"service_id":          param.ServiceID,
			"provider_id":         param.ProviderID,
			"provider_created_at": param.ProviderCreatedAt,
		}

		preparedQuery := batchInserts.Queue(
			insertIncidentsQuery,
			queryArgs,
		)

		preparedQuery.Query(func(rows pgx.Rows) error {
			incidentRow, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[incidentDto])
			if err != nil {
				c.lg.ErrorContext(ctx, "error collecting incident from batch", slog.String("provider_id", param.ProviderID), slog.Uint64("service_id", uint64(param.ServiceID)), slog.Any("err", err))
				return apperrors.InternalError("failed to collect incident from batch", err)
			}

			incidentResponse = append(incidentResponse, incidentRow)
			return nil
		})

	}

	err := c.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		c.lg.ErrorContext(ctx, "error while bulk inserting incidents", slog.Any("err", err))
		return nil, apperrors.InternalError("failed to bulk insert incidents", err)
	}

	response := lo.Map(incidentResponse, func(item incidentDto, _ int) incidents.IncidentResult {
		return incidents.IncidentResult{
			ID:                item.ID,
			Title:             item.Title,
			Link:              item.Link,
			ProviderImpact:    nullable.SetValue(item.ProviderImpact.String, item.ProviderImpact.Valid),
			Impact:            nullable.SetValue(item.Impact.String, item.Impact.Valid),
			ServiceID:         item.ServiceID,
			ProviderID:        item.ProviderID,
			ProviderCreatedAt: item.ProviderCreatedAt,
			CreatedAt:         item.CreatedAt,
			UpdatedAt:         item.UpdatedAt,
			DeletedAt:         nullable.SetValue(item.DeletedAt.Time, item.DeletedAt.Valid),
		}
	})

	return response, nil
}
