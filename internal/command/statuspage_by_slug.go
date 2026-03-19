package command

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/services"
)

type StatuspageBySlugCmd struct {
	logger       *slog.Logger
	servicesRepo services.Repository
}

func NewStatuspageBySlugCmd(
	logger *slog.Logger,
	repo services.Repository,
) StatuspageBySlugCmd {
	return StatuspageBySlugCmd{
		logger:       logger,
		servicesRepo: repo,
	}
}

type StatuspageBySlugParams struct {
	Slug string
}

type StatuspageBySlugResult struct {
	ID   uint
	Name string
	Slug string
}

func (c StatuspageBySlugCmd) Execute(ctx context.Context, params StatuspageBySlugParams) (StatuspageBySlugResult, error) {
	slug := strings.TrimSpace(params.Slug)
	if slug == "" {
		return StatuspageBySlugResult{}, ErrStatuspageNotFound
	}

	service, err := c.servicesRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.logger.WarnContext(ctx, "statuspage not found", slog.String("slug", slug))
			return StatuspageBySlugResult{}, ErrStatuspageNotFound
		}

		c.logger.ErrorContext(ctx, "failed to fetch statuspage by slug", slog.String("slug", slug), slog.Any("err", err))
		return StatuspageBySlugResult{}, err
	}

	return StatuspageBySlugResult{
		ID:   service.ID,
		Name: service.Title,
		Slug: service.Slug,
	}, nil
}
