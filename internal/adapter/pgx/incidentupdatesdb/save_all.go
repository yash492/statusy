package incidentupdatesdb

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

//go:embed queries/insert_incident_updates.sql
var insertIncidentUpdatesQuery string

type incidentUpdateDto struct {
	ID             uint             `db:"id"`
	IncidentID     uint             `db:"incident_id"`
	Description    string           `db:"description"`
	ProviderID     string           `db:"provider_id"`
	ProviderStatus string           `db:"provider_status"`
	Status         string           `db:"status"`
	StatusTime     time.Time        `db:"status_time"`
	CreatedAt      time.Time        `db:"created_at"`
	UpdatedAt      time.Time        `db:"updated_at"`
	DeletedAt      pgtype.Timestamp `db:"deleted_at"`
}

func (r *PostgresIncidentUpdatesRepository) SaveAll(ctx context.Context, params []incidents.IncidentUpdateParams) ([]incidents.IncidentUpdateResult, error) {
	batchInserts := &pgx.Batch{}
	updatesResponse := []incidentUpdateDto{}

	for _, param := range params {
		queryArgs := pgx.NamedArgs{
			"incident_id":     param.IncidentID,
			"description":     param.Description,
			"provider_id":     param.ProviderID,
			"provider_status": param.ProviderStatus,
			"status":          param.Status,
			"status_time":     param.StatusTime,
		}

		preparedQuery := batchInserts.Queue(insertIncidentUpdatesQuery, queryArgs)

		preparedQuery.Query(func(rows pgx.Rows) error {
			update, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[incidentUpdateDto])
			if err != nil {
				r.lg.ErrorContext(ctx, "error collecting incident update from batch", slog.Any("err", err))
				return err
			}

			updatesResponse = append(updatesResponse, *update)
			return nil
		})
	}

	err := r.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		r.lg.ErrorContext(ctx, "error while bulk inserting incident updates", slog.Any("err", err))
		return nil, err
	}

	result := lo.Map(updatesResponse, func(item incidentUpdateDto, _ int) incidents.IncidentUpdateResult {
		return incidents.IncidentUpdateResult{
			ID:             item.ID,
			IncidentID:     item.IncidentID,
			Description:    item.Description,
			ProviderID:     item.ProviderID,
			ProviderStatus: item.ProviderStatus,
			Status:         item.Status,
			StatusTime:     item.StatusTime,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
			DeletedAt:      nullable.SetValue(item.DeletedAt.Time, item.DeletedAt.Valid),
		}
	})

	return result, nil
}
