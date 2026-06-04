package command

import (
	"context"
	"log/slog"
	"strings"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

type AddViewServiceCmd struct {
	logger    *slog.Logger
	viewsRepo views.Repository
}

func NewAddViewServiceCmd(logger *slog.Logger, viewsRepo views.Repository) AddViewServiceCmd {
	return AddViewServiceCmd{
		logger:    logger,
		viewsRepo: viewsRepo,
	}
}

type AddViewServiceParams struct {
	ViewPublicID                 string
	ServiceID                    int
	IncludeAllComponents         bool
	MonitorIncidents             bool
	MonitorScheduledMaintenances bool
	ComponentIDs                 []int
	ComponentGroupIDs            []int
}

func (c AddViewServiceCmd) Execute(ctx context.Context, params AddViewServiceParams) (views.ViewService, error) {
	publicID := strings.TrimSpace(params.ViewPublicID)
	if publicID == "" {
		return views.ViewService{}, apperrors.InvalidInputError("public_id cannot be empty", nil)
	}

	if !params.MonitorIncidents && !params.MonitorScheduledMaintenances {
		return views.ViewService{}, apperrors.InvalidInputError("at least one alert type (Incidents or Scheduled Maintenances) must be monitored", nil)
	}

	view, err := c.viewsRepo.GetByPublicID(ctx, publicID)
	if err != nil {
		return views.ViewService{}, err
	}

	vs, err := c.viewsRepo.AddViewService(ctx, views.ViewService{
		ViewID:                       view.ID,
		ServiceID:                    uint(params.ServiceID),
		IncludeAllComponents:         params.IncludeAllComponents,
		MonitorIncidents:             params.MonitorIncidents,
		MonitorScheduledMaintenances: params.MonitorScheduledMaintenances,
	}, params.ComponentIDs, params.ComponentGroupIDs)
	if err != nil {
		return views.ViewService{}, err
	}

	return vs, nil
}
