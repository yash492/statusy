package command

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/views"
)

var ErrCannotChangeDefaultStatusOfOnlyView = errors.New("cannot change default status of the only view")

type EditViewParams struct {
	CurrentSlug string
	Name        string
	Slug        string
	Description string
	IsDefault   bool
}

type EditViewCmd struct {
	lg        *slog.Logger
	viewsRepo views.Repository
}

func NewEditViewCmd(lg *slog.Logger, viewsRepo views.Repository) EditViewCmd {
	return EditViewCmd{
		lg:        lg,
		viewsRepo: viewsRepo,
	}
}

func (c EditViewCmd) Execute(ctx context.Context, params EditViewParams) (views.View, error) {
	view, err := c.viewsRepo.GetBySlug(ctx, params.CurrentSlug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return views.View{}, ErrViewNotFound
		}
		c.lg.ErrorContext(ctx, "error fetching view for edit", slog.String("slug", params.CurrentSlug), slog.Any("err", err))
		return views.View{}, err
	}

	// Business rule check:
	// "if there's only one view the default status can't be changed"
	// So if changing is_default from true to false:
	if view.IsDefault && !params.IsDefault {
		count, err := c.viewsRepo.CountViews(ctx)
		if err != nil {
			return views.View{}, err
		}
		if count <= 1 {
			return views.View{}, ErrCannotChangeDefaultStatusOfOnlyView
		}
	}

	view.Name = params.Name
	view.Slug = params.Slug
	view.Description = params.Description
	view.IsDefault = params.IsDefault

	updatedView, err := c.viewsRepo.UpdateView(ctx, view)
	if err != nil {
		c.lg.ErrorContext(ctx, "error updating view", slog.Uint64("id", uint64(view.ID)), slog.Any("err", err))
		return views.View{}, err
	}

	return updatedView, nil
}
