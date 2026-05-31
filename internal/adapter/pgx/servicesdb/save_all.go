package servicesdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/services"
)

//go:embed queries/insert_services.sql
var insertServiceQuery string

type serviceDto struct {
	ID   uint
	Name string
	Slug string
	URL  string
}

func (s *PostgresServiceRepository) SaveAll(ctx context.Context, servicesYaml []services.ServiceParams) ([]services.ServiceResult, error) {

	batchInserts := &pgx.Batch{}
	servicesResponse := []serviceDto{}

	for _, service := range servicesYaml {
		queryArgs := pgx.NamedArgs{
			"name": service.Name,
			"slug": service.Slug,
			"url":  service.URL,
		}

		preparedQuery := batchInserts.Queue(
			insertServiceQuery,
			queryArgs,
		)

		preparedQuery.Query(func(rows pgx.Rows) error {
			serviceRow, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[serviceDto])
			if err != nil {
				s.lg.ErrorContext(ctx, "error collecting service from batch", slog.String("slug", service.Slug), slog.Any("err", err))
				return apperrors.InternalError("failed to collect service from batch", err)
			}

			servicesResponse = append(servicesResponse, serviceRow)
			return nil
		})

	}

	err := s.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		s.lg.ErrorContext(ctx, "error while bulk inserting services", slog.Any("err", err))
		return nil, apperrors.InternalError("failed to bulk insert services", err)
	}

	result := lo.Map(servicesResponse, func(item serviceDto, _ int) services.ServiceResult {
		return services.ServiceResult{
			ID:   item.ID,
			Name: item.Name,
			Slug: item.Slug,
			URL:  item.URL,
		}
	})

	return result, nil
}
