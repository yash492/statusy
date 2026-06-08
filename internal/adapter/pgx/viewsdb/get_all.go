package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/get_all_views.sql
var getAllViewsQuery string

//go:embed queries/count_views_with_filter.sql
var countViewsWithFilterQuery string

func (r *PostgresViewsRepository) GetAll(ctx context.Context, search string, limit int) ([]views.View, int64, error) {
	rows, err := r.readDB.Query(ctx, getAllViewsQuery, pgx.NamedArgs{
		"search": search,
		"limit":  limit,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying all views", slog.Any("err", err))
		return nil, 0, apperrors.InternalError("failed to query all views", err)
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting all views rows", slog.Any("err", err))
		return nil, 0, apperrors.InternalError("failed to collect all views rows", err)
	}

	var totalCount int64
	err = r.readDB.QueryRow(ctx, countViewsWithFilterQuery, pgx.NamedArgs{"search": search}).Scan(&totalCount)
	if err != nil {
		r.lg.ErrorContext(ctx, "error counting views", slog.Any("err", err))
		return nil, 0, apperrors.InternalError("failed to count views", err)
	}

	result := lo.Map(dtos, func(dto viewDto, _ int) views.View {
		return views.View{
			ID:          dto.ID,
			Name:        dto.Name,
			PublicID:    dto.PublicID,
			Description: dto.Description,
			IsDefault:   dto.IsDefault,
			CreatedAt:   dto.CreatedAt,
			UpdatedAt:   dto.UpdatedAt,
		}
	})
	return result, totalCount, nil
}
