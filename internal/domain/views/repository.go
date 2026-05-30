package views

import (
	"context"

	"github.com/yash492/statusy/internal/domain/services"
)

type Repository interface {
	GetDefault(ctx context.Context) (View, error)
	GetBySlug(ctx context.Context, slug string) (View, error)
	Save(ctx context.Context, view View) (View, error)
	GetServicesByViewID(ctx context.Context, viewID uint) ([]ViewServiceStatus, error)
	GetUnconfiguredServices(ctx context.Context, viewID uint) ([]services.ServiceResult, error)
}

