package command

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/views"
)

type GetOrCreateDefaultViewCmd struct {
	lg        *slog.Logger
	viewsRepo views.Repository
}

func NewGetOrCreateDefaultViewCmd(lg *slog.Logger, viewsRepo views.Repository) GetOrCreateDefaultViewCmd {
	return GetOrCreateDefaultViewCmd{
		lg:        lg,
		viewsRepo: viewsRepo,
	}
}

func (c GetOrCreateDefaultViewCmd) Execute(ctx context.Context) (views.View, error) {
	view, err := c.viewsRepo.GetDefault(ctx)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			c.lg.ErrorContext(ctx, "error fetching default view", slog.Any("err", err))
			return views.View{}, err
		}

		// No default view exists — seed one.
		view, err = c.viewsRepo.Save(ctx, views.View{
			Name:        "Default View",
			Slug:        "default-view",
			Description: "Default monitoring dashboard",
			IsDefault:   true,
			Services:    []views.ViewServiceStatus{},
		})
		if err != nil {
			c.lg.ErrorContext(ctx, "error seeding default view", slog.Any("err", err))
			return views.View{}, err
		}
	}

	return view, nil
}
