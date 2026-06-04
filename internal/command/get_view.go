package command

import (
	"context"
	"log/slog"

	"github.com/yash492/statusy/internal/domain/views"
)

type GetViewParams struct {
	PublicID string
}

type GetViewCmd struct {
	lg        *slog.Logger
	viewsRepo views.Repository
}

func NewGetViewCmd(lg *slog.Logger, viewsRepo views.Repository) GetViewCmd {
	return GetViewCmd{
		lg:        lg,
		viewsRepo: viewsRepo,
	}
}

func (c GetViewCmd) Execute(ctx context.Context, params GetViewParams) (views.View, error) {
	return c.viewsRepo.GetByPublicID(ctx, params.PublicID)
}
