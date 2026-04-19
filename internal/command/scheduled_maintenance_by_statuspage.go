package command

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/scheduledmaintenance"
	"github.com/yash492/statusy/internal/domain/services"
)

type ScheduledMaintenanceByStatuspageCmd struct {
	logger                    *slog.Logger
	ServicesRepo              services.Repository
	scheduledMaintenancesRepo scheduledmaintenance.Repository
}

func NewScheduledMaintenanceByStatuspageCmd(
	logger *slog.Logger,
	servicesRepo services.Repository,
	scheduledMaintenancesRepo scheduledmaintenance.Repository,
) ScheduledMaintenanceByStatuspageCmd {
	return ScheduledMaintenanceByStatuspageCmd{
		logger:                    logger,
		ServicesRepo:              servicesRepo,
		scheduledMaintenancesRepo: scheduledMaintenancesRepo,
	}
}

type ScheduledMaintenanceByStatuspageParams struct {
	StatuspageSlug string
	PageNumber     int
	PageSize       int
}

type ScheduledMaintenanceByStatuspageIncident struct {
	ID                uint
	Title             string
	Status            string
	StartsAt          time.Time
	EndsAt            time.Time
	ProviderCreatedAt time.Time
	Link              string
}

type ScheduledMaintenanceByStatuspageResult struct {
	ScheduledMaintenances []ScheduledMaintenanceByStatuspageIncident
	ServiceName           string
	ServiceSlug           string
}

func (c ScheduledMaintenanceByStatuspageCmd) Execute(ctx context.Context, params ScheduledMaintenanceByStatuspageParams) (ScheduledMaintenanceByStatuspageResult, error) {
	slug := strings.TrimSpace(params.StatuspageSlug)
	if slug == "" {
		return ScheduledMaintenanceByStatuspageResult{}, ErrStatuspageNotFound
	}

	service, err := c.ServicesRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.logger.WarnContext(ctx, "statuspage not found", slog.String("slug", slug))
			return ScheduledMaintenanceByStatuspageResult{}, ErrStatuspageNotFound
		}

		c.logger.ErrorContext(ctx, "failed to fetch statuspage service", slog.String("slug", slug), slog.Any("err", err))
		return ScheduledMaintenanceByStatuspageResult{}, err
	}

	pageNumber := params.PageNumber
	if pageNumber <= 0 {
		pageNumber = 1
	}

	pageSize := params.PageSize

	//Keeping it hardcoded for now
	if pageSize != 10 {
		pageSize = 10
	}

	offset := (pageNumber - 1) * pageSize

	maintenanceRows, err := c.scheduledMaintenancesRepo.GetByService(ctx, scheduledmaintenance.ScheduledMaintenanceByServiceParams{
		ServiceID: service.ID,
		Limit:     pageSize,
		Offset:    offset,
	})
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to fetch scheduled maintenances by statuspage", slog.String("slug", slug), slog.Any("service_id", service.ID), slog.Any("err", err))
		return ScheduledMaintenanceByStatuspageResult{}, err
	}

	maintenances := make([]ScheduledMaintenanceByStatuspageIncident, 0, len(maintenanceRows))
	for _, m := range maintenanceRows {
		maintenances = append(maintenances, ScheduledMaintenanceByStatuspageIncident{
			ID:                m.ID,
			Title:             m.Title,
			Status:            m.Status,
			StartsAt:          m.StartsAt,
			EndsAt:           m.EndsAt,
			ProviderCreatedAt: m.ProviderCreatedAt,
			Link:              m.Link,
		})
	}

	result := ScheduledMaintenanceByStatuspageResult{
		ScheduledMaintenances: maintenances,
		ServiceName:           service.Name,
		ServiceSlug:           service.Slug,
	}

	return result, nil
}
