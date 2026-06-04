package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"
)

//go:embed queries/count_views.sql
var countViewsQuery string

func (r *PostgresViewsRepository) CountViews(ctx context.Context) (int, error) {
	var count int
	err := r.readDB.QueryRow(ctx, countViewsQuery).Scan(&count)
	if err != nil {
		r.lg.ErrorContext(ctx, "error counting views", slog.Any("err", err))
		return 0, err
	}
	return count, nil
}
