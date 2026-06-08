package viewsdb

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/get_default_view.sql
var getDefaultViewQuery string

func (r *PostgresViewsRepository) GetDefault(ctx context.Context) (views.View, error) {
	rows, err := r.readDB.Query(ctx, getDefaultViewQuery)
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying default view", slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to query default view", err)
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return views.View{}, apperrors.NotFoundError("default view not found", err)
		}
		r.lg.ErrorContext(ctx, "error collecting default view row", slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to collect default view row", err)
	}

	return views.View{
		ID:          dto.ID,
		Name:        dto.Name,
		PublicID:    dto.PublicID,
		Description: dto.Description,
		IsDefault:   dto.IsDefault,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}, nil
}
