package command

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/views"
)

var ErrViewNotFound = errors.New("view not found")

type GetUnconfiguredServicesCmd struct {
	logger    *slog.Logger
	viewsRepo views.Repository
}

func NewGetUnconfiguredServicesCmd(logger *slog.Logger, viewsRepo views.Repository) GetUnconfiguredServicesCmd {
	return GetUnconfiguredServicesCmd{
		logger:    logger,
		viewsRepo: viewsRepo,
	}
}

type GetUnconfiguredServicesParams struct {
	ViewSlug string
}

func (c GetUnconfiguredServicesCmd) Execute(ctx context.Context, params GetUnconfiguredServicesParams) ([]services.ServiceResult, error) {
	slug := strings.TrimSpace(params.ViewSlug)
	if slug == "" {
		return nil, ErrViewNotFound
	}

	view, err := c.viewsRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.logger.WarnContext(ctx, "view not found", slog.String("slug", slug))
			return nil, ErrViewNotFound
		}
		c.logger.ErrorContext(ctx, "failed to fetch view", slog.String("slug", slug), slog.Any("err", err))
		return nil, err
	}

	unconfigured, err := c.viewsRepo.GetUnconfiguredServices(ctx, view.ID)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to fetch unconfigured services", slog.Uint64("view_id", uint64(view.ID)), slog.Any("err", err))
		return nil, err
	}

	return unconfigured, nil
}
