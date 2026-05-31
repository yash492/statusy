package command

import (
	"context"
	"errors"
	"log/slog"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

type GetOrCreateDefaultViewCmd struct {
	lg        *slog.Logger
	viewsRepo views.Repository
}

func NewGetOrCreateDefaultViewCmd(lg *slog.Logger, viewsRepo views.Repository) GetOrCreateDefaultViewCmd {
	return GetOrCreateDefaultViewCmd{
		lg:        lg,
		viewsRepo: viewsRepo,
	}
}

func (c GetOrCreateDefaultViewCmd) Execute(ctx context.Context) (views.View, error) {
	view, err := c.viewsRepo.GetDefault(ctx)
	if err != nil {
		if appErr, ok := errors.AsType[*apperrors.AppError](err); !ok || appErr.Type != apperrors.TypeNotFound {
			return views.View{}, err
		}

		// No default view exists — seed one.
		view, err = c.viewsRepo.Save(ctx, views.View{
			Name:        "Default View",
			Slug:        "default-view",
			Description: "Default monitoring dashboard",
			IsDefault:   true,
			Services:    []views.ViewServiceStatus{},
		})
		if err != nil {
			return views.View{}, err
		}
	}

	return view, nil
}
