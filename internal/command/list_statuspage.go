package command

import (
	"context"
	"log/slog"
	"strings"

	"github.com/yash492/statusy/internal/domain/services"
)

type ListStatuspageCmd struct {
	logger       *slog.Logger
	servicesRepo services.Repository
}

func NewListStatuspageCmd(
	logger *slog.Logger,
	repo services.Repository,
) ListStatuspageCmd {
	return ListStatuspageCmd{
		logger:       logger,
		servicesRepo: repo,
	}
}

type ListStatuspageParams struct {
	Search string
}

type ListStatuspageResult struct {
	ID   uint
	Name string
	Slug string
}

func (s ListStatuspageCmd) Execute(ctx context.Context, params ListStatuspageParams) ([]ListStatuspageResult, error) {

	search := strings.TrimSpace(params.Search)
	servicesList, err := s.servicesRepo.SearchByName(ctx, search)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to list status pages", slog.Any("err", err))
		return nil, err
	}

	results := make([]ListStatuspageResult, 0, len(servicesList))

	for _, svc := range servicesList {
		results = append(results, ListStatuspageResult{
			ID:   svc.ID,
			Name: svc.Name,
			Slug: svc.Slug,
		})
	}

	return results, nil
}
