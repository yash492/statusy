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
	search := ""
	if request.Params.Search != nil {
		search = *request.Params.Search
	}

	result, err := h.ListStatuspageCmd.Execute(ctx, command.ListStatuspageParams{
		Search: search,
	})
	if err != nil {
		return nil, err
	}

	statuspages := make(api.ListStatuspages200JSONResponse, 0, len(result))
	for _, r := range result {
		statuspages = append(statuspages, api.Statuspage{
			Id:   int(r.ID),
			Name: r.Name,
			Slug: r.Slug,
		})
	}

	return statuspages, nil
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

	pageNumber := 0
	pageSize := 0

	if request.Params.PageNumber != nil {
		pageNumber = *request.Params.PageNumber
	}

	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}

	result, err := h.IncidentByStatuspageCmd.Execute(ctx, command.IncidentByStatuspageParams{
		StatuspageSlug: request.StatuspageSlug,
		PageNumber:     pageNumber,
		PageSize:       pageSize,
	})

	if err != nil {
		return nil, err
	}

	incidents := make(api.IncidentByStatuspage200JSONResponse, 0, len(result))
	for _, r := range result {
		incidents = append(incidents, api.Incident{
			Id:                int(r.ID),
			Title:             r.Title,
			Status:            r.Status,
			ProviderCreatedAt: r.ProviderCreatedAt,
		})
	}

	return incidents, nil
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
