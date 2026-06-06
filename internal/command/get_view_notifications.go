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
	Search       string
	PageNumber   int
	PageSize     int
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

func (h GetViewNotificationsHandler) Handle(ctx context.Context, cmd GetViewNotifications) ([]notifications.ViewNotification, int64, error) {
	publicID := strings.TrimSpace(cmd.ViewPublicID)
	if publicID == "" {
		return nil, 0, apperrors.InvalidInputError("public_id cannot be empty", nil)
	}

	view, err := h.viewsRepo.GetByPublicID(ctx, publicID)
	if err != nil {
		return nil, 0, err
	}

	h.lg.Info("Fetching notification destinations for view", slog.Uint64("view_id", uint64(view.ID)))

	page := cmd.PageNumber
	if page <= 0 {
		page = 1
	}
	pageSize := cmd.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	limit := pageSize
	offset := (page - 1) * pageSize

	return h.repo.GetByViewID(ctx, view.ID, cmd.Search, limit, offset)
}
