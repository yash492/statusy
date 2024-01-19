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
	ServiceName               string                          `json:"service_name"`
	ServiceID                 uint                            `json:"service_id"`
	IsAllComponentsConfigured bool                            `json:"is_all_components_configured"`
	Components                []types.ComponentsWithNameAndID `json:"components"`
	Incidents                 []incidentRespHelper            `json:"incidents"`
}

type incidentRespHelper struct {
	ID                    uint      `json:"id"`
	LastUpdatedStatusTime time.Time `json:"last_updated_status_time"`
	Status                string    `json:"status"`
	CreatedAt             time.Time `json:"created_at"`
	Name                  string    `json:"name"`
	Link                  string    `json:"link"`
}

func SubscriptionIncidents(w http.ResponseWriter, r *http.Request) *api.Response {
	ctx := r.Context()
	subscriptionUUID := ctx.Value(types.SubscriptionIDCtx).(uuid.UUID)

	subscriptionIncidents, err := domain.Subscription.GetIncidentsForSubscription(subscriptionUUID, 0, 10)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot fetch incidents for the given subscription %v, err: %v", subscriptionUUID.String(), err.Error())
	}

	if len(subscriptionIncidents) < 1 {
		return api.Send(w, http.StatusOK, []any{}, types.JSON{
			"total_count": 0,
		})
	}

	totalCount := subscriptionIncidents[0].TotalCount
	isAllComponentsConfigured := subscriptionIncidents[0].IsAllComponentsConfigured
	serviceName := subscriptionIncidents[0].ServiceName
	serviceID := subscriptionIncidents[0].ServiceID
	var components []types.ComponentsWithNameAndID

	incidentsResp := lo.Map(subscriptionIncidents, func(incident schema.SubscriptionIncident, _ int) incidentRespHelper {
		return incidentRespHelper{
			ID:                    incident.IncidentID,
			LastUpdatedStatusTime: incident.LastUpdatedStatusTime,
			Status:                incident.IncidentStatus,
			CreatedAt:             incident.IncidentCreatedAt,
			Name:                  incident.IncidentName,
			Link:                  incident.IncidentLink,
		}
	})

	if !isAllComponentsConfigured {
		subscriptionComponents, err := domain.Subscription.GetWithComponents(subscriptionUUID)
		if err != nil {
			return api.Errorf(w, http.StatusInternalServerError, "cannot fetch components for the given subscription %v, err: %v", subscriptionUUID.String(), err.Error())
		}

		components = lo.Map(subscriptionComponents, func(component schema.SubscriptionWithComponent, _ int) types.ComponentsWithNameAndID {
			return types.ComponentsWithNameAndID{
				Name: component.ComponentName,
				ID:   component.ComponentID,
			}
		})
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
	}
	return api.Send(w, http.StatusOK, resp, meta)

}
