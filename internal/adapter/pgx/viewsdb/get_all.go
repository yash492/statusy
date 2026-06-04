package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/get_all_views.sql
var getAllViewsQuery string

func (r *PostgresViewsRepository) GetAll(ctx context.Context, search string) ([]views.View, error) {
	rows, err := r.readDB.Query(ctx, getAllViewsQuery, pgx.NamedArgs{"search": search})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying all views", slog.Any("err", err))
		return nil, apperrors.InternalError("failed to query all views", err)
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting all views rows", slog.Any("err", err))
		return nil, apperrors.InternalError("failed to collect all views rows", err)
	}

	result := make([]views.View, len(dtos))
	for i, dto := range dtos {
		result[i] = views.View{
			ID:          dto.ID,
			Name:        dto.Name,
			PublicID:    dto.PublicID,
			Description: dto.Description,
			IsDefault:   dto.IsDefault,
			CreatedAt:   dto.CreatedAt,
			UpdatedAt:   dto.UpdatedAt,
		}
	}
	return result, nil
}
