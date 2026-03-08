package servicesdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/services"
)

//go:embed queries/get_by_slug.sql
var getBySlugQuery string

func (s *PostgresServiceRepository) GetBySlug(ctx context.Context, slug string) (services.ServiceResult, error) {
	rows, err := s.readDB.Query(ctx, getBySlugQuery, pgx.NamedArgs{"slug": slug})
	if err != nil {
		s.lg.ErrorContext(ctx, "error querying service by slug", slog.String("slug", slug), slog.Any("err", err))
		return services.ServiceResult{}, err
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[serviceDto])
	if err != nil {
		s.lg.ErrorContext(ctx, "error collecting service by slug rows", slog.String("slug", slug), slog.Any("err", err))
		return services.ServiceResult{}, err
	}

	return services.ServiceResult{
		ID:   item.ID,
		Name: item.Name,
		Slug: item.Slug,
	}, nil
}
