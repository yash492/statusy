package command

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/views"
)

var ErrViewServiceNotFound = errors.New("view service not found")

type EditViewServiceCmd struct {
	logger    *slog.Logger
	viewsRepo views.Repository
}

func NewEditViewServiceCmd(logger *slog.Logger, viewsRepo views.Repository) EditViewServiceCmd {
	return EditViewServiceCmd{
		logger:    logger,
		viewsRepo: viewsRepo,
	}
}

type EditViewServiceParams struct {
	ViewSlug             string
	ServiceID            int
	IncludeAllComponents bool
	ComponentIDs         []int
	ComponentGroupIDs    []int
}

func (c EditViewServiceCmd) Execute(ctx context.Context, params EditViewServiceParams) (views.ViewService, error) {
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

	existingVS, err := c.viewsRepo.GetViewService(ctx, view.ID, uint(params.ServiceID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.logger.WarnContext(ctx, "view service not found", slog.Uint64("view_id", uint64(view.ID)), slog.Int("service_id", params.ServiceID))
			return views.ViewService{}, ErrViewServiceNotFound
		}
		c.logger.ErrorContext(ctx, "failed to fetch view service", slog.Any("err", err))
		return views.ViewService{}, err
	}

	vs, err := c.viewsRepo.UpdateViewService(ctx, views.ViewService{
		ID:                   existingVS.ID,
		ViewID:               existingVS.ViewID,
		ServiceID:            existingVS.ServiceID,
		IncludeAllComponents: params.IncludeAllComponents,
	}, params.ComponentIDs, params.ComponentGroupIDs)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to update view service", slog.Uint64("id", uint64(existingVS.ID)), slog.Any("err", err))
		return views.ViewService{}, err
	}

	return vs, nil
}
