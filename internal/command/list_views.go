package command

import (
	"context"
	"log/slog"

	"github.com/yash492/statusy/internal/domain/views"
)

type ListViewsCmd struct {
	lg        *slog.Logger
	viewsRepo views.Repository
}

func NewListViewsCmd(lg *slog.Logger, viewsRepo views.Repository) ListViewsCmd {
	return ListViewsCmd{
		lg:        lg,
		viewsRepo: viewsRepo,
	}
}

func (c ListViewsCmd) Execute(ctx context.Context, search string) ([]views.View, int64, error) {
	return c.viewsRepo.GetAll(ctx, search, 5)
}
