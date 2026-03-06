package servicesdb

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/services"
)

//go:embed queries/search_by_slug_services.sql
var searchBySlugServiceQuery string

func (s *PostgresServiceRepository) SearchBySlug(ctx context.Context, slug string) ([]services.ServiceResult, error) {
	// Query uses a named parameter `@slug`; match all by default
	args := pgx.NamedArgs{"slug": fmt.Sprintf("%%%s%%", slug)}

	rows, err := s.readDB.Query(ctx, searchBySlugServiceQuery, args)
	if err != nil {
		s.lg.ErrorContext(ctx, "error querying services", slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[serviceDto])
	if err != nil {
		if s.lg != nil {
			s.lg.ErrorContext(ctx, "error collecting service rows", slog.Any("err", err))
		}
		return nil, err
	}

	results := make([]services.ServiceResult, 0, len(dtos))
	for _, service := range dtos {
		results = append(results, services.ServiceResult{
			ID:                      service.ID,
			Name:                    service.Name,
			Slug:                    service.Slug,
			IncidentsUrl:            service.IncidentsUrl,
			ScheduleMaintenancesUrl: service.ScheduleMaintenancesUrl,
			ComponentsUrl:           service.ComponentsUrl,
			ProviderType:            services.ProviderType(service.ProviderType),
			CreatedAt:               service.CreatedAt,
			UpdatedAt:               service.UpdatedAt,
			DeletedAt:               nullable.SetValue(service.DeletedAt.Time, service.DeletedAt.Valid),
		})
	}

	return results, nil
}
