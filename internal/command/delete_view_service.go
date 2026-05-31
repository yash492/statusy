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
	ViewSlug  string
	ServiceID int
}

func (c DeleteViewServiceCmd) Execute(ctx context.Context, params DeleteViewServiceParams) error {
	slug := strings.TrimSpace(params.ViewSlug)
	if slug == "" {
		return apperrors.InvalidInputError("slug cannot be empty", nil)
	}

	view, err := c.viewsRepo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}

	err = c.viewsRepo.DeleteViewService(ctx, view.ID, uint(params.ServiceID))
	if err != nil {
		return err
	}

	return nil
}
