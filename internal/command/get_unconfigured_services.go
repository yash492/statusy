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
	ViewSlug string
	Search   string
}

func (c GetUnconfiguredServicesCmd) Execute(ctx context.Context, params GetUnconfiguredServicesParams) ([]services.ServiceResult, error) {
	slug := strings.TrimSpace(params.ViewSlug)
	if slug == "" {
		return nil, apperrors.InvalidInputError("slug cannot be empty", nil)
	}

	view, err := c.viewsRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	unconfigured, err := c.viewsRepo.GetUnconfiguredServices(ctx, view.ID, params.Search)
	if err != nil {
		return nil, err
	}

	return unconfigured, nil
}
