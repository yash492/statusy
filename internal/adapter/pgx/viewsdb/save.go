package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/common/snowflake"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/insert_views.sql
var insertViewQuery string

func (r *PostgresViewsRepository) Save(ctx context.Context, view views.View) (views.View, error) {
	if view.PublicID == "" {
		view.PublicID = snowflake.Generate()
	}

	rows, err := r.writeDB.Query(ctx, insertViewQuery, pgx.NamedArgs{
		"name":        view.Name,
		"public_id":   view.PublicID,
		"description": view.Description,
		"is_default":  view.IsDefault,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error inserting view", slog.String("public_id", view.PublicID), slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to insert view", err)
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting inserted view row", slog.String("public_id", view.PublicID), slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to collect inserted view row", err)
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
