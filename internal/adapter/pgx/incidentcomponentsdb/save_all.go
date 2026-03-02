package incidentcomponentsdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/incidents"
)

//go:embed queries/insert_incident_components.sql
var insertIncidentComponentsQuery string

type incidentComponentDto struct {
	ID          uint             `db:"id"`
	IncidentID  uint             `db:"incident_id"`
	ComponentID uint             `db:"component_id"`
	CreatedAt   time.Time        `db:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at"`
	DeletedAt   pgtype.Timestamp `db:"deleted_at"`
}

func (r *PostgresIncidentComponentsRepository) SaveAll(ctx context.Context, params []incidents.IncidentComponentParams) ([]incidents.IncidentComponentResult, error) {
	batchInserts := &pgx.Batch{}
	componentsResponse := []incidentComponentDto{}

	for _, param := range params {
		queryArgs := pgx.NamedArgs{
			"incident_id":  param.IncidentID,
			"component_id": param.ComponentID,
		}

		preparedQuery := batchInserts.Queue(insertIncidentComponentsQuery, queryArgs)

		preparedQuery.Query(func(rows pgx.Rows) error {
			component, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[incidentComponentDto])
			if err != nil {
				r.lg.ErrorContext(ctx, "error collecting incident component from batch", slog.Any("err", err))
				return err
			}

			componentsResponse = append(componentsResponse, *component)
			return nil
		})
	}

	err := r.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		r.lg.ErrorContext(ctx, "error while bulk inserting incident components", slog.Any("err", err))
		return nil, err
	}

	result := lo.Map(componentsResponse, func(item incidentComponentDto, _ int) incidents.IncidentComponentResult {
		return incidents.IncidentComponentResult{
			ID:          item.ID,
			IncidentID:  item.IncidentID,
			ComponentID: item.ComponentID,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			DeletedAt: nullable.Nullable[time.Time]{
				Value: item.DeletedAt.Time,
				Valid: item.DeletedAt.Valid,
			},
		}
	})

	return result, nil
}
