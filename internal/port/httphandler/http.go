package httphandler

import (
	"context"
	"log/slog"

	"github.com/yash492/statusy/internal/command"
	"github.com/yash492/statusy/internal/port/httphandler/generated/api"
)

var _ api.StrictServerInterface = Handler{}

type Handler struct {
	lg                                  *slog.Logger
	ListStatuspageCmd                   command.ListStatuspageCmd
	StatuspageBySlugCmd                 command.StatuspageBySlugCmd
	IncidentByStatuspageCmd             command.IncidentByStatuspageCmd
	ScheduledMaintenanceByStatuspageCmd command.ScheduledMaintenanceByStatuspageCmd
	GetOrCreateDefaultViewCmd           command.GetOrCreateDefaultViewCmd
	GetUnconfiguredServicesCmd           command.GetUnconfiguredServicesCmd
	GetServiceComponentsCmd              command.GetServiceComponentsCmd
	AddViewServiceCmd                    command.AddViewServiceCmd
	EditViewServiceCmd                   command.EditViewServiceCmd
	DeleteViewServiceCmd                 command.DeleteViewServiceCmd
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
			Url:  "",
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

	resp := api.IncidentByStatuspage200JSONResponse{
		Statuspage: api.Statuspage{
			Name: result.ServiceName,
			Slug: result.ServiceSlug,
			Id:   int(result.ServiceID),
			Url:  result.ServiceUrl,
		},
		Incidents:  make([]api.Incident, 0, len(result.Incidents)),
		TotalCount: result.TotalCount,
	}

	for _, incident := range result.Incidents {
		resp.Incidents = append(resp.Incidents, api.Incident{
			Id:                int(incident.ID),
			Title:             incident.Title,
			Status:            incident.Status,
			ProviderCreatedAt: incident.ProviderCreatedAt,
			IncidentUrl:       incident.Link,
		})
	}

	return resp, nil
}

// (GET /statuspages/{statuspageSlug}/schedule-maintenances)
func (h Handler) ScheduledMaintenanceByStatuspage(ctx context.Context, request api.ScheduledMaintenanceByStatuspageRequestObject) (api.ScheduledMaintenanceByStatuspageResponseObject, error) {

	pageNumber := 0
	pageSize := 0

	if request.Params.PageNumber != nil {
		pageNumber = *request.Params.PageNumber
	}

	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}

	result, err := h.ScheduledMaintenanceByStatuspageCmd.Execute(ctx, command.ScheduledMaintenanceByStatuspageParams{
		StatuspageSlug: request.StatuspageSlug,
		PageNumber:     pageNumber,
		PageSize:       pageSize,
	})

	if err != nil {
		return nil, err
	}

	resp := api.ScheduledMaintenanceByStatuspage200JSONResponse{
		Statuspage: api.Statuspage{
			Name: result.ServiceName,
			Slug: result.ServiceSlug,
			Id:   int(result.ServiceID),
			Url:  result.ServiceUrl,
		},
		ScheduledMaintenances: make([]api.ScheduledMaintenance, 0, len(result.ScheduledMaintenances)),
		TotalCount:            result.TotalCount,
	}

	for _, m := range result.ScheduledMaintenances {
		resp.ScheduledMaintenances = append(resp.ScheduledMaintenances, api.ScheduledMaintenance{
			Id:                      int(m.ID),
			Title:                   m.Title,
			Status:                  m.Status,
			StartsAt:                m.StartsAt,
			EndsAt:                  m.EndsAt,
			ProviderCreatedAt:       m.ProviderCreatedAt,
			ScheduledMaintenanceUrl: m.Link,
		})
	}

	return resp, nil
}

// (POST /api/views/default)
func (h Handler) GetDefaultView(ctx context.Context, request api.GetDefaultViewRequestObject) (api.GetDefaultViewResponseObject, error) {
	view, err := h.GetOrCreateDefaultViewCmd.Execute(ctx)
	if err != nil {
		return nil, err
	}

	services := make([]api.ViewServiceStatus, 0, len(view.Services))
	for _, s := range view.Services {
		services = append(services, api.ViewServiceStatus{
			Id:                   int(s.ID),
			Name:                 s.Name,
			Slug:                 s.Slug,
			Status:               s.Status,
			LastIncident:         s.LastIncident,
			IncludeAllComponents: s.IncludeAllComponents,
		})
	}

	return api.GetDefaultView200JSONResponse{
		Id:          int(view.ID),
		Name:        view.Name,
		Slug:        view.Slug,
		Description: view.Description,
		IsDefault:   view.IsDefault,
		Services:    services,
	}, nil
}

// (GET /statuspages/{statuspageSlug})
func (h Handler) StatuspageBySlug(ctx context.Context, request api.StatuspageBySlugRequestObject) (api.StatuspageBySlugResponseObject, error) {
	result, err := h.StatuspageBySlugCmd.Execute(ctx, command.StatuspageBySlugParams{
		Slug: request.StatuspageSlug,
	})
	if err != nil {
		return nil, err
	}

	return api.StatuspageBySlug200JSONResponse{
		Id:   int(result.ID),
		Name: result.Name,
		Slug: result.Slug,
		Url:  result.URL,
	}, nil
}

// (GET /views/{viewSlug}/unconfigured-services)
func (h Handler) GetUnconfiguredServices(ctx context.Context, request api.GetUnconfiguredServicesRequestObject) (api.GetUnconfiguredServicesResponseObject, error) {
	search := ""
	if request.Params.Search != nil {
		search = *request.Params.Search
	}

	result, err := h.GetUnconfiguredServicesCmd.Execute(ctx, command.GetUnconfiguredServicesParams{
		ViewSlug: request.ViewSlug,
		Search:   search,
	})
	if err != nil {
		return nil, err
	}

	servicesList := make(api.GetUnconfiguredServices200JSONResponse, 0, len(result))
	for _, s := range result {
		servicesList = append(servicesList, api.Statuspage{
			Id:   int(s.ID),
			Name: s.Name,
			Slug: s.Slug,
			Url:  s.URL,
		})
	}

	return servicesList, nil
}

// (GET /services/{serviceSlug}/components)
func (h Handler) GetServiceComponents(ctx context.Context, request api.GetServiceComponentsRequestObject) (api.GetServiceComponentsResponseObject, error) {
	result, err := h.GetServiceComponentsCmd.Execute(ctx, command.GetServiceComponentsParams{
		ServiceSlug: request.ServiceSlug,
	})
	if err != nil {
		return nil, err
	}

	grouped := make([]api.ComponentGroup, 0, len(result.GroupedComponents))
	for _, g := range result.GroupedComponents {
		comps := make([]api.Component, 0, len(g.Components))
		for _, c := range g.Components {
			comps = append(comps, api.Component{
				Id:         int(c.ID),
				Name:       c.Name,
				ProviderId: c.ProviderID,
			})
		}
		grouped = append(grouped, api.ComponentGroup{
			Id:         int(g.ID),
			Name:       g.Name,
			ProviderId: g.ProviderID,
			Components: comps,
		})
	}

	ungrouped := make([]api.Component, 0, len(result.UngroupedComponents))
	for _, c := range result.UngroupedComponents {
		ungrouped = append(ungrouped, api.Component{
			Id:         int(c.ID),
			Name:       c.Name,
			ProviderId: c.ProviderID,
		})
	}

	return api.GetServiceComponents200JSONResponse{
		ServiceId:           int(result.ServiceID),
		ServiceName:         result.ServiceName,
		ServiceSlug:         result.ServiceSlug,
		GroupedComponents:   grouped,
		UngroupedComponents: ungrouped,
	}, nil
}

// (POST /views/{viewSlug}/services)
func (h Handler) AddViewService(ctx context.Context, request api.AddViewServiceRequestObject) (api.AddViewServiceResponseObject, error) {
	componentIDs := []int{}
	if request.Body.ComponentIds != nil {
		componentIDs = *request.Body.ComponentIds
	}

	componentGroupIDs := []int{}
	if request.Body.ComponentGroupIds != nil {
		componentGroupIDs = *request.Body.ComponentGroupIds
	}

	result, err := h.AddViewServiceCmd.Execute(ctx, command.AddViewServiceParams{
		ViewSlug:             request.ViewSlug,
		ServiceID:            request.Body.ServiceId,
		IncludeAllComponents: request.Body.IncludeAllComponents,
		ComponentIDs:         componentIDs,
		ComponentGroupIDs:    componentGroupIDs,
	})
	if err != nil {
		return nil, err
	}

	return api.AddViewService200JSONResponse{
		Id:                   int(result.ID),
		ServiceId:            int(result.ServiceID),
		IncludeAllComponents: result.IncludeAllComponents,
	}, nil
}

// (PUT /views/{viewSlug}/services/{serviceId})
func (h Handler) EditViewService(ctx context.Context, request api.EditViewServiceRequestObject) (api.EditViewServiceResponseObject, error) {
	componentIDs := []int{}
	if request.Body.ComponentIds != nil {
		componentIDs = *request.Body.ComponentIds
	}

	componentGroupIDs := []int{}
	if request.Body.ComponentGroupIds != nil {
		componentGroupIDs = *request.Body.ComponentGroupIds
	}

	result, err := h.EditViewServiceCmd.Execute(ctx, command.EditViewServiceParams{
		ViewSlug:             request.ViewSlug,
		ServiceID:            request.ServiceId,
		IncludeAllComponents: request.Body.IncludeAllComponents,
		ComponentIDs:         componentIDs,
		ComponentGroupIDs:    componentGroupIDs,
	})
	if err != nil {
		return nil, err
	}

	return api.EditViewService200JSONResponse{
		Id:                   int(result.ID),
		ServiceId:            int(result.ServiceID),
		IncludeAllComponents: result.IncludeAllComponents,
	}, nil
}

// (DELETE /views/{viewSlug}/services/{serviceId})
func (h Handler) DeleteViewService(ctx context.Context, request api.DeleteViewServiceRequestObject) (api.DeleteViewServiceResponseObject, error) {
	err := h.DeleteViewServiceCmd.Execute(ctx, command.DeleteViewServiceParams{
		ViewSlug:  request.ViewSlug,
		ServiceID: request.ServiceId,
	})
	if err != nil {
		return nil, err
	}

	return api.DeleteViewService204Response{}, nil
}

