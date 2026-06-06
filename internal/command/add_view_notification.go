package command

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/yash492/statusy/internal/domain/notifications"
)

type AddViewNotification struct {
	ViewID uint
	Type   notifications.NotificationType
	Config json.RawMessage
}

type AddViewNotificationHandler struct {
	lg   *slog.Logger
	repo notifications.NotificationsRepository
}

func NewAddViewNotificationHandler(lg *slog.Logger, repo notifications.NotificationsRepository) AddViewNotificationHandler {
	return AddViewNotificationHandler{
		lg:   lg,
		repo: repo,
	}
}

func (h AddViewNotificationHandler) Handle(ctx context.Context, cmd AddViewNotification) (notifications.ViewNotification, error) {
	h.lg.Info("Adding notification destination to view", slog.Uint64("view_id", uint64(cmd.ViewID)))
	vn := notifications.ViewNotification{
		ViewID: cmd.ViewID,
		Type:   cmd.Type,
		Config: cmd.Config,
	}
	return h.repo.Save(ctx, vn)
}
