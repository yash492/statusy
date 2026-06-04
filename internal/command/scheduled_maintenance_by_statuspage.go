package command

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/yash492/statusy/internal/common/apperrors"
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
	StatuspageSlug    string
	ComponentIDs      []int
	ComponentGroupIDs []int
	PageNumber        int
	PageSize          int
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
	ServiceUrl            string
	ServiceID             uint
	TotalCount            int
}

func (c ScheduledMaintenanceByStatuspageCmd) Execute(ctx context.Context, params ScheduledMaintenanceByStatuspageParams) (ScheduledMaintenanceByStatuspageResult, error) {
	slug := strings.TrimSpace(params.StatuspageSlug)
	if slug == "" {
		return ScheduledMaintenanceByStatuspageResult{}, apperrors.InvalidInputError("slug cannot be empty", fmt.Errorf("slug cannot be empty"))
	}

	service, err := c.ServicesRepo.GetBySlug(ctx, slug)
	if err != nil {
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
		ServiceID:         service.ID,
		ComponentIDs:      params.ComponentIDs,
		ComponentGroupIDs: params.ComponentGroupIDs,
		Limit:             pageSize,
		Offset:            offset,
	})
	if err != nil {
		return ScheduledMaintenanceByStatuspageResult{}, err
	}

	totalCount := 0
	if len(maintenanceRows) > 0 {
		totalCount = int(maintenanceRows[0].TotalCount)
	}

	maintenances := make([]ScheduledMaintenanceByStatuspageIncident, 0, len(maintenanceRows))
	for _, m := range maintenanceRows {
		maintenances = append(maintenances, ScheduledMaintenanceByStatuspageIncident{
			ID:                m.ID,
			Title:             m.Title,
			Status:            m.Status,
			StartsAt:          m.StartsAt,
			EndsAt:            m.EndsAt,
			ProviderCreatedAt: m.ProviderCreatedAt,
			Link:              m.Link,
		})
	}

	result := ScheduledMaintenanceByStatuspageResult{
		ScheduledMaintenances: maintenances,
		ServiceName:           service.Name,
		ServiceSlug:           service.Slug,
		ServiceUrl:            service.URL,
		TotalCount:            totalCount,
		ServiceID:             service.ID,
	}

	return result, nil
}
