package postgres

import (
	"context"

	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"github.com/yash492/statusy/internal/repository/services"
	"github.com/yash492/statusy/sqlc/db"
)

func (s *PostgresServiceRepository) GetAll(ctx context.Context) ([]services.ServiceResult, error) {

	var resp []db.Service
	err := s.read(func(q *db.Queries) error {
		dbResp, err := q.GetAllServices(ctx)
		if err != nil {
			return err
		}
		resp = dbResp
		return nil
	})

	if err != nil {
		return nil, err
	}

	services := lo.Map(resp, func(serviceResp db.Service, _ int) services.ServiceResult {
		return services.ServiceResult{
			ID:                      uint(serviceResp.ID),
			Name:                    serviceResp.Name,
			Slug:                    serviceResp.Slug,
			IncidentsUrl:            serviceResp.IncidentsUrl,
			ScheduleMaintenancesUrl: serviceResp.ScheduleMaintenancesUrl,
			ComponentsUrl:           serviceResp.ComponentsUrl,
			ProviderType:            statuspage.ProviderType(serviceResp.ProviderType),
		}
	})

	return services, nil
}
