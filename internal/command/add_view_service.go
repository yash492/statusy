package command

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/views"
)

type AddViewServiceCmd struct {
	logger    *slog.Logger
	viewsRepo views.Repository
}

func NewAddViewServiceCmd(logger *slog.Logger, viewsRepo views.Repository) AddViewServiceCmd {
	return AddViewServiceCmd{
		logger:    logger,
		viewsRepo: viewsRepo,
	}
}

type AddViewServiceParams struct {
	ViewSlug             string
	ServiceID            int
	IncludeAllComponents bool
	ComponentIDs         []int
	ComponentGroupIDs    []int
}

func (c AddViewServiceCmd) Execute(ctx context.Context, params AddViewServiceParams) (views.ViewService, error) {
	slug := strings.TrimSpace(params.ViewSlug)
	if slug == "" {
		return views.ViewService{}, ErrViewNotFound
	}

	view, err := c.viewsRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.logger.WarnContext(ctx, "view not found", slog.String("slug", slug))
			return views.ViewService{}, ErrViewNotFound
		}
		c.logger.ErrorContext(ctx, "failed to fetch view", slog.String("slug", slug), slog.Any("err", err))
		return views.ViewService{}, err
	}

	vs, err := c.viewsRepo.AddViewService(ctx, views.ViewService{
		ViewID:               view.ID,
		ServiceID:            uint(params.ServiceID),
		IncludeAllComponents: params.IncludeAllComponents,
	}, params.ComponentIDs, params.ComponentGroupIDs)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to add view service", slog.Uint64("view_id", uint64(view.ID)), slog.Int("service_id", params.ServiceID), slog.Any("err", err))
		return views.ViewService{}, err
	}

	return vs, nil
}
