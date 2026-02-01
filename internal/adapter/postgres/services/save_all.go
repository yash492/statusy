package postgres

import (
	"context"

	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/repository/services"
	"github.com/yash492/statusy/sqlc/db"
)

func (s *PostgresServiceRepository) SaveAll(ctx context.Context, servicesYaml []services.ServiceParams) error {
	err := s.write(func(q *db.Queries) error {
		_, err := q.CreateServices(ctx, lo.Map(servicesYaml, func(service services.ServiceParams, _ int) db.CreateServicesParams {
			return db.CreateServicesParams{
				Name:                    service.Name,
				Slug:                    service.Slug,
				ProviderType:            service.ProviderType.String(),
				IncidentsUrl:            service.IncidentsUrl,
				ScheduleMaintenancesUrl: service.ScheduleMaintenancesUrl,
				ComponentsUrl:           service.ComponentsUrl,
			}
		}))
		return err
	})
	return err
}
