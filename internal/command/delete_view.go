package command

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/views"
)

var ErrDefaultViewCannotBeDeleted = errors.New("default view cannot be deleted")

type DeleteViewParams struct {
	Slug string
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
	view, err := c.viewsRepo.GetBySlug(ctx, params.Slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrViewNotFound
		}
		c.lg.ErrorContext(ctx, "error fetching view for delete", slog.String("slug", params.Slug), slog.Any("err", err))
		return err
	}

	if view.IsDefault {
		return ErrDefaultViewCannotBeDeleted
	}

	err = c.viewsRepo.DeleteView(ctx, view.ID)
	if err != nil {
		c.lg.ErrorContext(ctx, "error deleting view", slog.Uint64("id", uint64(view.ID)), slog.Any("err", err))
		return err
	}

	return nil
}
