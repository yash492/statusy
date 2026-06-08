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

//go:embed queries/get_view_by_public_id.sql
var getViewByPublicIDQuery string

func (r *PostgresViewsRepository) GetByPublicID(ctx context.Context, publicID string) (views.View, error) {
	rows, err := r.readDB.Query(ctx, getViewByPublicIDQuery, pgx.NamedArgs{"public_id": publicID})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying view by public_id", slog.String("public_id", publicID), slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to query view by public_id", err)
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return views.View{}, apperrors.NotFoundError("view not found", err)
		}
		r.lg.ErrorContext(ctx, "error collecting view row by public_id", slog.String("public_id", publicID), slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to collect view row by public_id", err)
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
