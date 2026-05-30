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
	GetUnconfiguredServices(ctx context.Context, viewID uint, search string) ([]services.ServiceResult, error)
	GetViewService(ctx context.Context, viewID uint, serviceID uint) (ViewService, error)
	AddViewService(ctx context.Context, vs ViewService, componentIDs []int, componentGroupIDs []int) (ViewService, error)
	UpdateViewService(ctx context.Context, vs ViewService, componentIDs []int, componentGroupIDs []int) (ViewService, error)
	DeleteViewService(ctx context.Context, viewID uint, serviceID uint) error

	UpdateView(ctx context.Context, view View) (View, error)
	DeleteView(ctx context.Context, viewID uint) error
	CountViews(ctx context.Context) (int, error)
}


