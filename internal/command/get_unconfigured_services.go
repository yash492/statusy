package command

import (
	"context"
	"log/slog"
	"strings"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/views"
)

type GetUnconfiguredServicesCmd struct {
	logger    *slog.Logger
	viewsRepo views.Repository
}

func NewGetUnconfiguredServicesCmd(logger *slog.Logger, viewsRepo views.Repository) GetUnconfiguredServicesCmd {
	return GetUnconfiguredServicesCmd{
		logger:    logger,
		viewsRepo: viewsRepo,
	}
}

type GetUnconfiguredServicesParams struct {
	ViewPublicID string
	Search       string
}

func (c GetUnconfiguredServicesCmd) Execute(ctx context.Context, params GetUnconfiguredServicesParams) ([]services.ServiceResult, error) {
	publicID := strings.TrimSpace(params.ViewPublicID)
	if publicID == "" {
		return nil, apperrors.InvalidInputError("public_id cannot be empty", nil)
	}

	view, err := c.viewsRepo.GetByPublicID(ctx, publicID)
	if err != nil {
		return nil, err
	}

	unconfigured, err := c.viewsRepo.GetUnconfiguredServices(ctx, view.ID, params.Search)
	if err != nil {
		return nil, err
	}

	return unconfigured, nil
}
