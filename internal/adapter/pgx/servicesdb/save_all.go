package servicesdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common"
	"github.com/yash492/statusy/internal/domain/services"
)

//go:embed queries/insert_services.sql
var insertServiceQuery string

type serviceDto struct {
	ID                      uint
	Name                    string
	Slug                    string
	IncidentsUrl            string
	ScheduleMaintenancesUrl string
	ComponentsUrl           string
	ProviderType            string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               pgtype.Timestamptz
}

func (s *PostgresServiceRepository) SaveAll(ctx context.Context, servicesYaml []services.ServiceParams) ([]services.ServiceResult, error) {

	batchInserts := &pgx.Batch{}
	servicesResponse := []serviceDto{}

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
			service, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[serviceDto])
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

	result := lo.Map(servicesResponse, func(item serviceDto, _ int) services.ServiceResult {
		return services.ServiceResult{
			ID:                      item.ID,
			Name:                    item.Name,
			Slug:                    item.Slug,
			IncidentsUrl:            item.IncidentsUrl,
			ScheduleMaintenancesUrl: item.ScheduleMaintenancesUrl,
			ComponentsUrl:           item.ComponentsUrl,
			ProviderType:            services.ProviderType(item.ProviderType),
			CreatedAt:               item.CreatedAt,
			UpdatedAt:               item.UpdatedAt,
			DeletedAt: common.Nullable[time.Time]{
				Value: item.DeletedAt.Time,
				Valid: item.DeletedAt.Valid,
			},
		}
	})

	return result, nil
}
