package httphandler

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/yash492/statusy/internal/command"
	"github.com/yash492/statusy/internal/domain/notifications"
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
	GetUnconfiguredServicesCmd          command.GetUnconfiguredServicesCmd
	GetServiceComponentsCmd             command.GetServiceComponentsCmd
	AddViewServiceCmd                   command.AddViewServiceCmd
	EditViewServiceCmd                  command.EditViewServiceCmd
	GetViewServiceCmd                   command.GetViewServiceCmd
	DeleteViewServiceCmd                command.DeleteViewServiceCmd
	EditViewCmd                         command.EditViewCmd
	DeleteViewCmd                       command.DeleteViewCmd
	GetViewServicesCmd                  command.GetViewServicesCmd
	ListViewsCmd                        command.ListViewsCmd
	CreateViewCmd                       command.CreateViewCmd
	GetViewCmd                          command.GetViewCmd
	AddViewNotificationCmd              command.AddViewNotificationHandler
	GetViewNotificationsCmd             command.GetViewNotificationsHandler
	EditViewNotificationCmd             command.EditViewNotificationHandler
	DeleteViewNotificationCmd           command.DeleteViewNotificationHandler
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

	var componentIDs []int
	if request.Params.ComponentIds != nil {
		componentIDs = *request.Params.ComponentIds
	}

	var componentGroupIDs []int
	if request.Params.ComponentGroupIds != nil {
		componentGroupIDs = *request.Params.ComponentGroupIds
	}

	result, err := h.IncidentByStatuspageCmd.Execute(ctx, command.IncidentByStatuspageParams{
		StatuspageSlug:    request.StatuspageSlug,
		ComponentIDs:      componentIDs,
		ComponentGroupIDs: componentGroupIDs,
		PageNumber:        pageNumber,
		PageSize:          pageSize,
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

	var componentIDs []int
	if request.Params.ComponentIds != nil {
		componentIDs = *request.Params.ComponentIds
	}

	var componentGroupIDs []int
	if request.Params.ComponentGroupIds != nil {
		componentGroupIDs = *request.Params.ComponentGroupIds
	}

	result, err := h.ScheduledMaintenanceByStatuspageCmd.Execute(ctx, command.ScheduledMaintenanceByStatuspageParams{
		StatuspageSlug:    request.StatuspageSlug,
		ComponentIDs:      componentIDs,
		ComponentGroupIDs: componentGroupIDs,
		PageNumber:        pageNumber,
		PageSize:          pageSize,
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
func (h Handler) CreateOrGetDefaultView(ctx context.Context, request api.CreateOrGetDefaultViewRequestObject) (api.CreateOrGetDefaultViewResponseObject, error) {
	view, err := h.GetOrCreateDefaultViewCmd.Execute(ctx)
	if err != nil {
		return nil, err
	}

	return api.CreateOrGetDefaultView200JSONResponse{
		Name:        view.Name,
		PublicId:    view.PublicID,
		Description: view.Description,
		IsDefault:   view.IsDefault,
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

// (GET /views/{public_id}/unconfigured-services)
func (h Handler) GetUnconfiguredServices(ctx context.Context, request api.GetUnconfiguredServicesRequestObject) (api.GetUnconfiguredServicesResponseObject, error) {
	search := ""
	if request.Params.Search != nil {
		search = *request.Params.Search
	}

	result, err := h.GetUnconfiguredServicesCmd.Execute(ctx, command.GetUnconfiguredServicesParams{
		ViewPublicID: request.PublicId,
		Search:       search,
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

// (POST /views/{public_id}/services)
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
		ViewPublicID:                 request.PublicId,
		ServiceID:                    request.Body.ServiceId,
		IncludeAllComponents:         request.Body.IncludeAllComponents,
		MonitorIncidents:             request.Body.MonitorIncidents,
		MonitorScheduledMaintenances: request.Body.MonitorScheduledMaintenances,
		ComponentIDs:                 componentIDs,
		ComponentGroupIDs:            componentGroupIDs,
	})
	if err != nil {
		return nil, err
	}

	compIDs := result.ComponentIDs
	if compIDs == nil {
		compIDs = []int{}
	}

	compGrpIDs := result.ComponentGroupIDs
	if compGrpIDs == nil {
		compGrpIDs = []int{}
	}

	return api.AddViewService200JSONResponse{
		Id:                           int(result.ID),
		ServiceId:                    int(result.ServiceID),
		IncludeAllComponents:         result.IncludeAllComponents,
		MonitorIncidents:             result.MonitorIncidents,
		MonitorScheduledMaintenances: result.MonitorScheduledMaintenances,
		ComponentIds:                 &compIDs,
		ComponentGroupIds:            &compGrpIDs,
	}, nil
}

// (GET /views/{public_id}/services/{serviceId})
func (h Handler) GetViewService(ctx context.Context, request api.GetViewServiceRequestObject) (api.GetViewServiceResponseObject, error) {
	result, err := h.GetViewServiceCmd.Execute(ctx, command.GetViewServiceParams{
		ViewPublicID: request.PublicId,
		ServiceID:    request.ServiceId,
	})
	if err != nil {
		return nil, err
	}

	compIDs := result.ComponentIDs
	if compIDs == nil {
		compIDs = []int{}
	}

	compGrpIDs := result.ComponentGroupIDs
	if compGrpIDs == nil {
		compGrpIDs = []int{}
	}

	return api.GetViewService200JSONResponse{
		Id:                           int(result.ID),
		ServiceId:                    int(result.ServiceID),
		IncludeAllComponents:         result.IncludeAllComponents,
		MonitorIncidents:             result.MonitorIncidents,
		MonitorScheduledMaintenances: result.MonitorScheduledMaintenances,
		ComponentIds:                 &compIDs,
		ComponentGroupIds:            &compGrpIDs,
	}, nil
}

// (PUT /views/{public_id}/services/{serviceId})
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
		ViewPublicID:                 request.PublicId,
		ServiceID:                    request.ServiceId,
		IncludeAllComponents:         request.Body.IncludeAllComponents,
		MonitorIncidents:             request.Body.MonitorIncidents,
		MonitorScheduledMaintenances: request.Body.MonitorScheduledMaintenances,
		ComponentIDs:                 componentIDs,
		ComponentGroupIDs:            componentGroupIDs,
	})
	if err != nil {
		return nil, err
	}

	compIDs := result.ComponentIDs
	if compIDs == nil {
		compIDs = []int{}
	}

	compGrpIDs := result.ComponentGroupIDs
	if compGrpIDs == nil {
		compGrpIDs = []int{}
	}

	return api.EditViewService200JSONResponse{
		Id:                           int(result.ID),
		ServiceId:                    int(result.ServiceID),
		IncludeAllComponents:         result.IncludeAllComponents,
		MonitorIncidents:             result.MonitorIncidents,
		MonitorScheduledMaintenances: result.MonitorScheduledMaintenances,
		ComponentIds:                 &compIDs,
		ComponentGroupIds:            &compGrpIDs,
	}, nil
}

// (DELETE /views/{public_id}/services/{serviceId})
func (h Handler) DeleteViewService(ctx context.Context, request api.DeleteViewServiceRequestObject) (api.DeleteViewServiceResponseObject, error) {
	err := h.DeleteViewServiceCmd.Execute(ctx, command.DeleteViewServiceParams{
		ViewPublicID: request.PublicId,
		ServiceID:    request.ServiceId,
	})
	if err != nil {
		return nil, err
	}

	return api.DeleteViewService204Response{}, nil
}

// (PUT /views/{public_id})
func (h Handler) EditView(ctx context.Context, request api.EditViewRequestObject) (api.EditViewResponseObject, error) {
	result, err := h.EditViewCmd.Execute(ctx, command.EditViewParams{
		PublicID:    request.PublicId,
		Name:        request.Body.Name,
		Description: request.Body.Description,
		IsDefault:   request.Body.IsDefault,
	})
	if err != nil {
		return nil, err
	}

	return api.EditView200JSONResponse{
		Name:        result.Name,
		PublicId:    result.PublicID,
		Description: result.Description,
		IsDefault:   result.IsDefault,
	}, nil
}

// (DELETE /views/{public_id})
func (h Handler) DeleteView(ctx context.Context, request api.DeleteViewRequestObject) (api.DeleteViewResponseObject, error) {
	err := h.DeleteViewCmd.Execute(ctx, command.DeleteViewParams{
		PublicID: request.PublicId,
	})
	if err != nil {
		return nil, err
	}

	return api.DeleteView204Response{}, nil
}

// (GET /views/{public_id})
func (h Handler) GetView(ctx context.Context, request api.GetViewRequestObject) (api.GetViewResponseObject, error) {
	result, err := h.GetViewCmd.Execute(ctx, command.GetViewParams{
		PublicID: request.PublicId,
	})
	if err != nil {
		return nil, err
	}

	return api.GetView200JSONResponse{
		Name:        result.Name,
		PublicId:    result.PublicID,
		Description: result.Description,
		IsDefault:   result.IsDefault,
	}, nil
}

// (GET /views/{publicId}/view-services)
func (h Handler) GetViewServices(ctx context.Context, request api.GetViewServicesRequestObject) (api.GetViewServicesResponseObject, error) {
	pageNumber := 0
	pageSize := 0
	search := ""

	if request.Params.PageNumber != nil {
		pageNumber = *request.Params.PageNumber
	}

	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}

	if request.Params.Search != nil {
		search = *request.Params.Search
	}

	result, err := h.GetViewServicesCmd.Execute(ctx, command.GetViewServicesParams{
		PublicID:   request.PublicId,
		Search:     search,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	})
	if err != nil {
		return nil, err
	}

	services := make([]api.ViewServiceStatus, 0, len(result.Services))
	for _, s := range result.Services {
		compIDs := s.ComponentIDs
		compGroupIDs := s.ComponentGroupIDs
		services = append(services, api.ViewServiceStatus{
			Id:                           int(s.ID),
			Name:                         s.Name,
			Slug:                         s.Slug,
			Status:                       s.Status,
			LastIncident:                 s.LastIncident,
			LastIncidentLink:             s.LastIncidentLink,
			IncludeAllComponents:         s.IncludeAllComponents,
			MonitorIncidents:             s.MonitorIncidents,
			MonitorScheduledMaintenances: s.MonitorScheduledMaintenances,
			UpcomingMaintenance:          s.UpcomingMaintenance,
			UpcomingMaintenanceLink:      s.UpcomingMaintenanceLink,
			ComponentIds:                 &compIDs,
			ComponentGroupIds:            &compGroupIDs,
		})
	}

	return api.GetViewServices200JSONResponse{
		Services:   services,
		TotalCount: int(result.TotalCount),
		UpCount:    int(result.UpCount),
		DownCount:  int(result.DownCount),
	}, nil
}

// (GET /views)
func (h Handler) ListViews(ctx context.Context, request api.ListViewsRequestObject) (api.ListViewsResponseObject, error) {
	search := ""
	if request.Params.Search != nil {
		search = *request.Params.Search
	}

	result, totalCount, err := h.ListViewsCmd.Execute(ctx, search)
	if err != nil {
		return nil, err
	}

	viewsList := make([]api.View, 0, len(result))
	for _, v := range result {
		viewsList = append(viewsList, api.View{
			Name:        v.Name,
			PublicId:    v.PublicID,
			Description: v.Description,
			IsDefault:   v.IsDefault,
		})
	}

	return api.ListViews200JSONResponse{
		Views:      viewsList,
		TotalCount: int(totalCount),
	}, nil
}

// (POST /views)
func (h Handler) CreateView(ctx context.Context, request api.CreateViewRequestObject) (api.CreateViewResponseObject, error) {
	result, err := h.CreateViewCmd.Execute(ctx, command.CreateViewParams{
		Name:        request.Body.Name,
		Description: request.Body.Description,
	})
	if err != nil {
		return nil, err
	}

	return api.CreateView200JSONResponse{
		Name:        result.Name,
		PublicId:    result.PublicID,
		Description: result.Description,
		IsDefault:   result.IsDefault,
	}, nil
}

// (GET /views/{publicId}/notifications)
func (h Handler) ListViewNotifications(ctx context.Context, request api.ListViewNotificationsRequestObject) (api.ListViewNotificationsResponseObject, error) {
	pageNumber := 1
	pageSize := 20
	search := ""

	if request.Params.PageNumber != nil {
		pageNumber = *request.Params.PageNumber
	}

	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}

	if request.Params.Search != nil {
		search = *request.Params.Search
	}

	result, totalCount, err := h.GetViewNotificationsCmd.Handle(ctx, command.GetViewNotifications{
		ViewPublicID: request.PublicId,
		Search:       search,
		PageNumber:   pageNumber,
		PageSize:     pageSize,
	})
	if err != nil {
		return nil, err
	}

	list := make([]api.ViewNotificationResponse, 0, len(result))
	for _, vn := range result {
		configMap, err := rawJSONToMap(vn.Config)
		if err != nil {
			return nil, err
		}
		list = append(list, api.ViewNotificationResponse{
			Id:        int(vn.ID),
			Name:      vn.Name,
			Type:      string(vn.Type),
			Config:    configMap,
			CreatedAt: vn.CreatedAt,
		})
	}

	return api.ListViewNotifications200JSONResponse{
		Notifications: list,
		TotalCount:    int(totalCount),
	}, nil
}

// (POST /views/{publicId}/notifications)
func (h Handler) AddViewNotification(ctx context.Context, request api.AddViewNotificationRequestObject) (api.AddViewNotificationResponseObject, error) {
	configJSON, err := json.Marshal(request.Body.Config)
	if err != nil {
		return nil, err
	}

	result, err := h.AddViewNotificationCmd.Handle(ctx, command.AddViewNotification{
		ViewPublicID: request.PublicId,
		Name:         request.Body.Name,
		Type:         notifications.NotificationType(request.Body.Type),
		Config:       configJSON,
	})
	if err != nil {
		return nil, err
	}

	configMap, err := rawJSONToMap(result.Config)
	if err != nil {
		return nil, err
	}

	return api.AddViewNotification200JSONResponse{
		Id:        int(result.ID),
		Name:      result.Name,
		Type:      string(result.Type),
		Config:    configMap,
		CreatedAt: result.CreatedAt,
	}, nil
}

// (PUT /views/{publicId}/notifications/{notificationId})
func (h Handler) EditViewNotification(ctx context.Context, request api.EditViewNotificationRequestObject) (api.EditViewNotificationResponseObject, error) {
	configJSON, err := json.Marshal(request.Body.Config)
	if err != nil {
		return nil, err
	}

	result, err := h.EditViewNotificationCmd.Handle(ctx, command.EditViewNotification{
		ID:     uint(request.NotificationId),
		Name:   request.Body.Name,
		Type:   notifications.NotificationType(request.Body.Type),
		Config: configJSON,
	})
	if err != nil {
		return nil, err
	}

	configMap, err := rawJSONToMap(result.Config)
	if err != nil {
		return nil, err
	}

	return api.EditViewNotification200JSONResponse{
		Id:        int(result.ID),
		Name:      result.Name,
		Type:      string(result.Type),
		Config:    configMap,
		CreatedAt: result.CreatedAt,
	}, nil
}

// (DELETE /views/{publicId}/notifications/{notificationId})
func (h Handler) DeleteViewNotification(ctx context.Context, request api.DeleteViewNotificationRequestObject) (api.DeleteViewNotificationResponseObject, error) {
	err := h.DeleteViewNotificationCmd.Handle(ctx, command.DeleteViewNotification{
		ID: uint(request.NotificationId),
	})
	if err != nil {
		return nil, err
	}

	return api.DeleteViewNotification204Response{}, nil
}

func rawJSONToMap(raw json.RawMessage) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}
