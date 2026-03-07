package httphandler

import (
	"context"
	"log/slog"

	"github.com/yash492/statusy/internal/command"
	"github.com/yash492/statusy/internal/port/httphandler/generated/api"
)

var _ api.StrictServerInterface = Handler{}

type Handler struct {
	lg                      *slog.Logger
	ListStatuspageCmd       command.ListStatuspageCmd
	IncidentByStatuspageCmd command.IncidentByStatuspageCmd
}

// (GET /statuspages)
func (h Handler) ListStatuspages(ctx context.Context, request api.ListStatuspagesRequestObject) (api.ListStatuspagesResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/feed.atom)
func (h Handler) GetAtomFeed(ctx context.Context, request api.GetAtomFeedRequestObject) (api.GetAtomFeedResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/feed.rss)
func (h Handler) GetRssFeed(ctx context.Context, request api.GetRssFeedRequestObject) (api.GetRssFeedResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/incidents)
func (h Handler) IncidentByStatuspage(ctx context.Context, request api.IncidentByStatuspageRequestObject) (api.IncidentByStatuspageResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/incidents/{incidentID})
func (h Handler) IncidentInfo(ctx context.Context, request api.IncidentInfoRequestObject) (api.IncidentInfoResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/schedule-maintenances)
func (h Handler) ScheduleMaintenanceByStatuspage(ctx context.Context, request api.ScheduleMaintenanceByStatuspageRequestObject) (api.ScheduleMaintenanceByStatuspageResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/schedule-maintenances/{scheduleMaintenanceID})
func (h Handler) ScheduleMaintenanceInfo(ctx context.Context, request api.ScheduleMaintenanceInfoRequestObject) (api.ScheduleMaintenanceInfoResponseObject, error) {
	return nil, nil
}
