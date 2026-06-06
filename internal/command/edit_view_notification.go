package command

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/yash492/statusy/internal/domain/notifications"
)

type EditViewNotification struct {
	ID     uint
	Name   string
	Type   notifications.NotificationType
	Config json.RawMessage
}

type EditViewNotificationHandler struct {
	lg   *slog.Logger
	repo notifications.NotificationsRepository
}

func NewEditViewNotificationHandler(lg *slog.Logger, repo notifications.NotificationsRepository) EditViewNotificationHandler {
	return EditViewNotificationHandler{
		lg:   lg,
		repo: repo,
	}
}

func (h EditViewNotificationHandler) Handle(ctx context.Context, cmd EditViewNotification) (notifications.ViewNotification, error) {
	h.lg.Info("Editing view notification destination", slog.Uint64("id", uint64(cmd.ID)))
	vn := notifications.ViewNotification{
		ID:     cmd.ID,
		Name:   cmd.Name,
		Type:   cmd.Type,
		Config: cmd.Config,
	}
	return h.repo.Update(ctx, vn)
}
