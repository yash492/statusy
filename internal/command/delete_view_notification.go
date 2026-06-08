package command

import (
	"context"
	"log/slog"

	"github.com/yash492/statusy/internal/domain/notifications"
)

type DeleteViewNotification struct {
	ID uint
}

type DeleteViewNotificationHandler struct {
	lg   *slog.Logger
	repo notifications.NotificationsRepository
}

func NewDeleteViewNotificationHandler(lg *slog.Logger, repo notifications.NotificationsRepository) DeleteViewNotificationHandler {
	return DeleteViewNotificationHandler{
		lg:   lg,
		repo: repo,
	}
}

func (h DeleteViewNotificationHandler) Handle(ctx context.Context, cmd DeleteViewNotification) error {
	h.lg.Info("Deleting view notification destination", slog.Uint64("id", uint64(cmd.ID)))
	return h.repo.Delete(ctx, cmd.ID)
}
