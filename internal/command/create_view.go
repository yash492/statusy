package command

import (
	"context"
	"log/slog"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

type CreateViewParams struct {
	Name        string
	Description string
}

type CreateViewCmd struct {
	lg        *slog.Logger
	viewsRepo views.Repository
}

func NewCreateViewCmd(lg *slog.Logger, viewsRepo views.Repository) CreateViewCmd {
	return CreateViewCmd{
		lg:        lg,
		viewsRepo: viewsRepo,
	}
}

func (c CreateViewCmd) Execute(ctx context.Context, params CreateViewParams) (views.View, error) {
	if params.Name == "" {
		return views.View{}, apperrors.ConflictError("view name cannot be empty", nil)
	}

	view := views.View{
		Name:        params.Name,
		Description: params.Description,
		IsDefault:   false,
	}

	createdView, err := c.viewsRepo.Save(ctx, view)
	if err != nil {
		return views.View{}, err
	}

	return createdView, nil
}
