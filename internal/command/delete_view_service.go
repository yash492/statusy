package command

import (
	"context"
	"log/slog"
	"strings"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

type DeleteViewServiceCmd struct {
	logger    *slog.Logger
	viewsRepo views.Repository
}

func NewDeleteViewServiceCmd(logger *slog.Logger, viewsRepo views.Repository) DeleteViewServiceCmd {
	return DeleteViewServiceCmd{
		logger:    logger,
		viewsRepo: viewsRepo,
	}
}

type DeleteViewServiceParams struct {
	ViewPublicID string
	ServiceID    int
}

func (c DeleteViewServiceCmd) Execute(ctx context.Context, params DeleteViewServiceParams) error {
	publicID := strings.TrimSpace(params.ViewPublicID)
	if publicID == "" {
		return apperrors.InvalidInputError("public_id cannot be empty", nil)
	}

	view, err := c.viewsRepo.GetByPublicID(ctx, publicID)
	if err != nil {
		return err
	}

	err = c.viewsRepo.DeleteViewService(ctx, view.ID, uint(params.ServiceID))
	if err != nil {
		return err
	}

	return nil
}
