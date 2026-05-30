package command

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
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
		return ErrViewNotFound
	}

	view, err := c.viewsRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.logger.WarnContext(ctx, "view not found", slog.String("slug", slug))
			return ErrViewNotFound
		}
		c.logger.ErrorContext(ctx, "failed to fetch view", slog.String("slug", slug), slog.Any("err", err))
		return err
	}

	err = c.viewsRepo.DeleteViewService(ctx, view.ID, uint(params.ServiceID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.logger.WarnContext(ctx, "view service not found for delete", slog.Uint64("view_id", uint64(view.ID)), slog.Int("service_id", params.ServiceID))
			return ErrViewServiceNotFound
		}
		c.logger.ErrorContext(ctx, "failed to delete view service", slog.Uint64("view_id", uint64(view.ID)), slog.Int("service_id", params.ServiceID), slog.Any("err", err))
		return err
	}

	return nil
}
