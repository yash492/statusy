package services

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/repository/services"
)

//go:embed queries/insert_services.sql
var insertServiceQuery string

func (s *PostgresServiceRepository) SaveAll(ctx context.Context, servicesYaml []services.ServiceParams) ([]services.ServiceResult, error) {

	batchInserts := &pgx.Batch{}
	servicesResponse := []services.ServiceResult{}

	for _, service := range servicesYaml {
		queryArgs := pgx.NamedArgs{
			"name":                      service.Name,
			"slug":                      service.Slug,
			"components_url":            service.ComponentsUrl,
			"incidents_url":             service.IncidentsUrl,
			"schedule_maintenances_url": service.ScheduleMaintenancesUrl,
			"provider_type":             service.ProviderType,
		}

		preparedQuery := batchInserts.Queue(
			insertServiceQuery,
			queryArgs,
		)

		preparedQuery.Query(func(rows pgx.Rows) error {
			service, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[services.ServiceResult])
			if err != nil {
				s.lg.ErrorContext(ctx, "error collecting service %s from batch", service.Slug, slog.Any("err", err))
				return err
			}

			servicesResponse = append(servicesResponse, *service)
			return nil
		})

	}

	err := s.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		s.lg.ErrorContext(ctx, "error while bulk inserting services", slog.Any("err", err))
		return nil, err
	}
	return servicesResponse, nil
}
