package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

type subscriptionByIDResponse struct {
	ServiceID       uint                        `json:"service_id"`
	ServiceName     string                      `json:"service_name"`
	UUID            string                      `json:"uuid"`
	IsAllComponents bool                        `json:"is_all_components"`
	Components      []componentsForSubscription `json:"components"`
}

type componentsForSubscription struct {
	IsConfigured bool   `json:"is_configured"`
	Name         string `json:"name"`
	ID           uint   `json:"id"`
}

func SubscriptionByID(w http.ResponseWriter, r *http.Request) *api.Response {
	ctx := r.Context()
	subscriptionID := ctx.Value(types.SubscriptionIDCtx).(uuid.UUID)

	// The result of this query will contain both null values and actual values for
	// subscription table fields such as is_all_components and uuid
	// since joins have been employed to fetch all the components configured and available
	// (not configured) for that subscription belonging to the particular service
	// The result is an array because it has many different components
	subscription, err := domain.Subscription.GetWithComponents(subscriptionID)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot retrieve subscription ID for %v, err: %v", subscriptionID, err.Error())
	}

	isAllComponents := true

	components := lo.Map(subscription, func(sub schema.SubscriptionWithComponent, _ int) componentsForSubscription {
		if sub.IsConfigured {
			isAllComponents = false
		}
		return componentsForSubscription{
			IsConfigured: sub.IsConfigured,
			Name:         sub.ComponentName,
			ID:           sub.ComponentID,
		}
	})

	subscriptionResp := subscriptionByIDResponse{
		ServiceID:       subscription[0].ServiceID,
		ServiceName:     subscription[0].ServiceName,
		UUID:            subscriptionID.String(),
		IsAllComponents: isAllComponents,
		Components:      components,
	}

	return api.Send(w, http.StatusOK, subscriptionResp, nil)
}
