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
	ViewSlug             string
	ServiceID            int
	IncludeAllComponents bool
	ComponentIDs         []int
	ComponentGroupIDs    []int
}

func (c AddViewServiceCmd) Execute(ctx context.Context, params AddViewServiceParams) (views.ViewService, error) {
	slug := strings.TrimSpace(params.ViewSlug)
	if slug == "" {
		return views.ViewService{}, apperrors.InvalidInputError("slug cannot be empty", nil)
	}

	view, err := c.viewsRepo.GetBySlug(ctx, slug)
	if err != nil {
		return views.ViewService{}, err
	}

	vs, err := c.viewsRepo.AddViewService(ctx, views.ViewService{
		ViewID:               view.ID,
		ServiceID:            uint(params.ServiceID),
		IncludeAllComponents: params.IncludeAllComponents,
	}, params.ComponentIDs, params.ComponentGroupIDs)
	if err != nil {
		return views.ViewService{}, err
	}

	return vs, nil
}
