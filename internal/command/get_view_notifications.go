package command

import (
	"context"
	"log/slog"
	"strings"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/notifications"
	"github.com/yash492/statusy/internal/domain/views"
)

type GetViewNotifications struct {
	ViewPublicID string
}

type GetViewNotificationsHandler struct {
	lg        *slog.Logger
	repo      notifications.NotificationsRepository
	viewsRepo views.Repository
}

func NewGetViewNotificationsHandler(lg *slog.Logger, repo notifications.NotificationsRepository, viewsRepo views.Repository) GetViewNotificationsHandler {
	return GetViewNotificationsHandler{
		lg:        lg,
		repo:      repo,
		viewsRepo: viewsRepo,
	}
}

func (h GetViewNotificationsHandler) Handle(ctx context.Context, cmd GetViewNotifications) ([]notifications.ViewNotification, error) {
	publicID := strings.TrimSpace(cmd.ViewPublicID)
	if publicID == "" {
		return nil, apperrors.InvalidInputError("public_id cannot be empty", nil)
	}

	view, err := h.viewsRepo.GetByPublicID(ctx, publicID)
	if err != nil {
		return nil, err
	}

	h.lg.Info("Fetching notification destinations for view", slog.Uint64("view_id", uint64(view.ID)))
	return h.repo.GetByViewID(ctx, view.ID)
}
