package command

import (
	"context"
	"log/slog"
	"strings"

	"github.com/yash492/statusy/internal/common/apperrors"
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
	URL  string
}

func (c StatuspageBySlugCmd) Execute(ctx context.Context, params StatuspageBySlugParams) (StatuspageBySlugResult, error) {
	slug := strings.TrimSpace(params.Slug)
	if slug == "" {
		return StatuspageBySlugResult{}, apperrors.InvalidInputError("slug cannot be empty", nil)
	}

	service, err := c.servicesRepo.GetBySlug(ctx, slug)
	if err != nil {
		return StatuspageBySlugResult{}, err
	}

	return StatuspageBySlugResult{
		ID:   service.ID,
		Name: service.Name,
		Slug: service.Slug,
		URL:  service.URL,
	}, nil
}
