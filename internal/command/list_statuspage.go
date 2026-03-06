package command

import (
	"context"
	"log/slog"
	"strings"

	"github.com/yash492/statusy/internal/domain/services"
)

type ListStatuspageCmd struct {
	logger       *slog.Logger
	ServicesRepo services.Repository
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
	servicesList, err := s.ServicesRepo.SearchBySlug(ctx, search)
	if err != nil {
		if s.logger != nil {
			s.logger.ErrorContext(ctx, "failed to list status pages", slog.Any("err", err))
		}
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
