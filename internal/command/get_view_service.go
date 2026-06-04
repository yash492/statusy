package command

import (
	"context"
	"log/slog"
	"strings"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

type GetViewServiceCmd struct {
	logger    *slog.Logger
	viewsRepo views.Repository
}

func NewGetViewServiceCmd(logger *slog.Logger, viewsRepo views.Repository) GetViewServiceCmd {
	return GetViewServiceCmd{
		logger:    logger,
		viewsRepo: viewsRepo,
	}
}

type GetViewServiceParams struct {
	ViewPublicID string
	ServiceID    int
}

func (c GetViewServiceCmd) Execute(ctx context.Context, params GetViewServiceParams) (views.ViewService, error) {
	publicID := strings.TrimSpace(params.ViewPublicID)
	if publicID == "" {
		return views.ViewService{}, apperrors.InvalidInputError("public_id cannot be empty", nil)
	}

	view, err := c.viewsRepo.GetByPublicID(ctx, publicID)
	if err != nil {
		return views.ViewService{}, err
	}

	vs, err := c.viewsRepo.GetViewService(ctx, view.ID, uint(params.ServiceID))
	if err != nil {
		return views.ViewService{}, err
	}

	return vs, nil
}
