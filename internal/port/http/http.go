package http

import (
	"context"

	"github.com/yash492/statusy/internal/port/generated/api"
)

var _ api.StrictServerInterface = HttpHandler{}

type HttpHandler struct {
}

// (GET /statuspages)
func (h HttpHandler) ListStatuspages(ctx context.Context, request api.ListStatuspagesRequestObject) (api.ListStatuspagesResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/feed.atom)
func (h HttpHandler) GetAtomFeed(ctx context.Context, request api.GetAtomFeedRequestObject) (api.GetAtomFeedResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/feed.rss)
func (h HttpHandler) GetRssFeed(ctx context.Context, request api.GetRssFeedRequestObject) (api.GetRssFeedResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/incidents)
func (h HttpHandler) IncidentByStatuspage(ctx context.Context, request api.IncidentByStatuspageRequestObject) (api.IncidentByStatuspageResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/incidents/{incidentID})
func (h HttpHandler) IncidentInfo(ctx context.Context, request api.IncidentInfoRequestObject) (api.IncidentInfoResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/schedule-maintenances)
func (h HttpHandler) ScheduleMaintenanceByStatuspage(ctx context.Context, request api.ScheduleMaintenanceByStatuspageRequestObject) (api.ScheduleMaintenanceByStatuspageResponseObject, error) {
	return nil, nil
}

// (GET /statuspages/{statuspageSlug}/schedule-maintenances/{scheduleMaintenanceID})
func (h HttpHandler) ScheduleMaintenanceInfo(ctx context.Context, request api.ScheduleMaintenanceInfoRequestObject) (api.ScheduleMaintenanceInfoResponseObject, error) {
	return nil, nil
}
