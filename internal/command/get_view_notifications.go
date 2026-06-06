package command

import (
	"context"
	"log/slog"

	"github.com/yash492/statusy/internal/domain/notifications"
)

type GetViewNotifications struct {
	ViewID uint
}

type GetViewNotificationsHandler struct {
	lg   *slog.Logger
	repo notifications.NotificationsRepository
}

func NewGetViewNotificationsHandler(lg *slog.Logger, repo notifications.NotificationsRepository) GetViewNotificationsHandler {
	return GetViewNotificationsHandler{
		lg:   lg,
		repo: repo,
	}
}

func (h GetViewNotificationsHandler) Handle(ctx context.Context, cmd GetViewNotifications) ([]notifications.ViewNotification, error) {
	h.lg.Info("Fetching notification destinations for view", slog.Uint64("view_id", uint64(cmd.ViewID)))
	return h.repo.GetByViewID(ctx, cmd.ViewID)
}
