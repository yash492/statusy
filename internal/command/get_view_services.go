package command

import (
	"context"
	"log/slog"

	"github.com/yash492/statusy/internal/domain/views"
)

type GetViewServicesCmd struct {
	logger    *slog.Logger
	viewsRepo views.Repository
}

func NewGetViewServicesCmd(logger *slog.Logger, viewsRepo views.Repository) GetViewServicesCmd {
	return GetViewServicesCmd{
		logger:    logger,
		viewsRepo: viewsRepo,
	}
}

type GetViewServicesParams struct {
	PublicID   string
	Search     string
	PageNumber int
	PageSize   int
}

type GetViewServicesResult struct {
	Services   []views.ViewServiceStatus
	TotalCount int64
	UpCount    int64
	DownCount  int64
}

func (c GetViewServicesCmd) Execute(ctx context.Context, params GetViewServicesParams) (GetViewServicesResult, error) {
	view, err := c.viewsRepo.GetByPublicID(ctx, params.PublicID)
	if err != nil {
		return GetViewServicesResult{}, err
	}

	pageNumber := params.PageNumber
	if pageNumber <= 0 {
		pageNumber = 1
	}

	pageSize := params.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (pageNumber - 1) * pageSize

	services, totalCount, upCount, downCount, err := c.viewsRepo.GetViewServices(ctx, view.ID, params.Search, pageSize, offset)
	if err != nil {
		return GetViewServicesResult{}, err
	}

	return GetViewServicesResult{
		Services:   services,
		TotalCount: totalCount,
		UpCount:    upCount,
		DownCount:  downCount,
	}, nil
}
