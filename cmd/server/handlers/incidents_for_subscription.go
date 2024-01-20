package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

type subscriptionIncidentsResp struct {
	ServiceName               string                           `json:"service_name"`
	ServiceID                 uint                             `json:"service_id"`
	IsAllComponentsConfigured bool                             `json:"is_all_components_configured"`
	Components                []subscriptionIncidentComponents `json:"components"`
	Incidents                 []incidentRespHelper             `json:"incidents"`
}

type incidentRespHelper struct {
	ID                    uint      `json:"id"`
	LastUpdatedStatusTime time.Time `json:"last_updated_status_time"`
	NormalisedStatus      string    `json:"normalised_status"`
	Status                string    `json:"status"`
	CreatedAt             time.Time `json:"created_at"`
	Name                  string    `json:"name"`
	Link                  string    `json:"link"`
}

type subscriptionIncidentComponents struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
}

func SubscriptionIncidents(w http.ResponseWriter, r *http.Request) *api.Response {
	ctx := r.Context()
	subscriptionUUID := ctx.Value(types.SubscriptionIDCtx).(uuid.UUID)
	queryParams := r.URL.Query()
	pageNumberStr := queryParams.Get("page_number")
	pageLimitStr := queryParams.Get("page_limit")

	pageNumber := parsePaginationParams(pageNumberStr, 0)
	pageLimit := parsePaginationParams(pageLimitStr, 5)

	subscriptionIncidents, err := domain.Subscription.GetIncidentsForSubscription(subscriptionUUID, pageNumber*pageLimit, pageLimit)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot fetch incidents for the given subscription %v, err: %v", subscriptionUUID.String(), err.Error())
	}

	totalCount := subscriptionIncidents[0].TotalCount
	isAllComponentsConfigured := subscriptionIncidents[0].IsAllComponentsConfigured
	serviceName := subscriptionIncidents[0].ServiceName
	serviceID := subscriptionIncidents[0].ServiceID
	components := []subscriptionIncidentComponents{}
	incidentsResp := []incidentRespHelper{}

	if totalCount > 0 {
		incidentsResp = lo.Map(subscriptionIncidents, func(incident schema.SubscriptionIncident, _ int) incidentRespHelper {
			return incidentRespHelper{
				ID:                    uint(incident.IncidentID.Int64),
				LastUpdatedStatusTime: incident.LastUpdatedStatusTime.Time,
				Status:                incident.IncidentStatus.String,
				CreatedAt:             incident.IncidentCreatedAt.Time,
				Name:                  incident.IncidentName.String,
				Link:                  incident.IncidentLink.String,
				NormalisedStatus:      incident.IncidentNormalisedStatus.String,
			}
		})
	}

	if !isAllComponentsConfigured {
		subscriptionComponents, err := domain.Subscription.GetWithComponents(subscriptionUUID)
		if err != nil {
			return api.Errorf(w, http.StatusInternalServerError, "cannot fetch components for the given subscription %v, err: %v", subscriptionUUID.String(), err.Error())
		}

		for _, component := range subscriptionComponents {
			if component.IsConfigured {
				components = append(components, subscriptionIncidentComponents{
					Name: component.ComponentName,
					ID:   component.ComponentID,
				})
			}
		}
	}

	resp := subscriptionIncidentsResp{
		ServiceName:               serviceName,
		ServiceID:                 serviceID,
		Components:                components,
		IsAllComponentsConfigured: isAllComponentsConfigured,
		Incidents:                 incidentsResp,
	}

	meta := types.JSON{
		"total_count": totalCount,
		"page_number": pageNumber,
		"page_limit":  pageLimit,
	}
	return api.Send(w, http.StatusOK, resp, meta)

}
