package command

import (
	"context"
	"log/slog"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

type DeleteViewParams struct {
	PublicID string
}

type DeleteViewCmd struct {
	lg        *slog.Logger
	viewsRepo views.Repository
}

func NewDeleteViewCmd(lg *slog.Logger, viewsRepo views.Repository) DeleteViewCmd {
	return DeleteViewCmd{
		lg:        lg,
		viewsRepo: viewsRepo,
	}
}

func (c DeleteViewCmd) Execute(ctx context.Context, params DeleteViewParams) error {
	view, err := c.viewsRepo.GetByPublicID(ctx, params.PublicID)
	if err != nil {
		return err
	}

	if view.IsDefault {
		return apperrors.ConflictError("default view cannot be deleted", nil)
	}

	err = c.viewsRepo.DeleteView(ctx, view.ID)
	if err != nil {
		return err
	}

	return nil
}
