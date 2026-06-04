package command

import (
	"context"
	"log/slog"
	"strings"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

type EditViewServiceCmd struct {
	logger    *slog.Logger
	viewsRepo views.Repository
}

func NewEditViewServiceCmd(logger *slog.Logger, viewsRepo views.Repository) EditViewServiceCmd {
	return EditViewServiceCmd{
		logger:    logger,
		viewsRepo: viewsRepo,
	}
}

type EditViewServiceParams struct {
	ViewPublicID                 string
	ServiceID                    int
	IncludeAllComponents         bool
	MonitorIncidents             bool
	MonitorScheduledMaintenances bool
	ComponentIDs                 []int
	ComponentGroupIDs            []int
}

func (c EditViewServiceCmd) Execute(ctx context.Context, params EditViewServiceParams) (views.ViewService, error) {
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

	existingVS, err := c.viewsRepo.GetViewService(ctx, view.ID, uint(params.ServiceID))
	if err != nil {
		return views.ViewService{}, err
	}

	vs, err := c.viewsRepo.UpdateViewService(ctx, views.ViewService{
		ID:                           existingVS.ID,
		ViewID:                       existingVS.ViewID,
		ServiceID:                    existingVS.ServiceID,
		IncludeAllComponents:         params.IncludeAllComponents,
		MonitorIncidents:             params.MonitorIncidents,
		MonitorScheduledMaintenances: params.MonitorScheduledMaintenances,
	}, params.ComponentIDs, params.ComponentGroupIDs)
	if err != nil {
		return views.ViewService{}, err
	}

	return vs, nil
}
