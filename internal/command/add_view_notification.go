package command

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/notifications"
	"github.com/yash492/statusy/internal/domain/views"
)

type AddViewNotification struct {
	ViewPublicID string
	Name         string
	Type         notifications.NotificationType
	Config       json.RawMessage
}

type AddViewNotificationHandler struct {
	lg        *slog.Logger
	repo      notifications.NotificationsRepository
	viewsRepo views.Repository
}

func NewAddViewNotificationHandler(lg *slog.Logger, repo notifications.NotificationsRepository, viewsRepo views.Repository) AddViewNotificationHandler {
	return AddViewNotificationHandler{
		lg:        lg,
		repo:      repo,
		viewsRepo: viewsRepo,
	}
}

func (h AddViewNotificationHandler) Handle(ctx context.Context, cmd AddViewNotification) (notifications.ViewNotification, error) {
	publicID := strings.TrimSpace(cmd.ViewPublicID)
	if publicID == "" {
		return notifications.ViewNotification{}, apperrors.InvalidInputError("public_id cannot be empty", nil)
	}

	view, err := h.viewsRepo.GetByPublicID(ctx, publicID)
	if err != nil {
		return notifications.ViewNotification{}, err
	}

	h.lg.Info("Adding notification destination to view", slog.Uint64("view_id", uint64(view.ID)))
	vn := notifications.ViewNotification{
		ViewID: view.ID,
		Name:   cmd.Name,
		Type:   cmd.Type,
		Config: cmd.Config,
	}
	return h.repo.Save(ctx, vn)
}
