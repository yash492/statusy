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
	ViewSlug             string
	ServiceID            int
	IncludeAllComponents bool
	ComponentIDs         []int
	ComponentGroupIDs    []int
}

func (c EditViewServiceCmd) Execute(ctx context.Context, params EditViewServiceParams) (views.ViewService, error) {
	slug := strings.TrimSpace(params.ViewSlug)
	if slug == "" {
		return views.ViewService{}, apperrors.InvalidInputError("slug cannot be empty", nil)
	}

	view, err := c.viewsRepo.GetBySlug(ctx, slug)
	if err != nil {
		return views.ViewService{}, err
	}

	existingVS, err := c.viewsRepo.GetViewService(ctx, view.ID, uint(params.ServiceID))
	if err != nil {
		return views.ViewService{}, err
	}

	vs, err := c.viewsRepo.UpdateViewService(ctx, views.ViewService{
		ID:                   existingVS.ID,
		ViewID:               existingVS.ViewID,
		ServiceID:            existingVS.ServiceID,
		IncludeAllComponents: params.IncludeAllComponents,
	}, params.ComponentIDs, params.ComponentGroupIDs)
	if err != nil {
		return views.ViewService{}, err
	}

	return vs, nil
}
