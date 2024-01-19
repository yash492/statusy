package handlers

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

type dashboardListResp struct {
	IncidentID       uint   `json:"incident_id"`
	ServiceID        uint   `json:"service_id"`
	ServiceName      string `json:"service_name"`
	SubscriptionUUID string `json:"subscription_uuid"`
	IncidentName     string `json:"incident_name"`
	IncidentLink     string `json:"incident_link"`
	IncidentImpact   string `json:"incident_impact"`
	IsDown           bool   `json:"is_down"`
}

// All subscriptions in a table as of now, accrues into a single dashboard
func DashboardList(w http.ResponseWriter, r *http.Request) *api.Response {

	queryParams := r.URL.Query()

	serviceName := queryParams.Get("service_name")
	pageNumberStr := queryParams.Get("page_number")
	pageLimitStr := queryParams.Get("page_limit")

	pageNumber := parsePaginationParams(pageNumberStr, 0)
	pageLimit := parsePaginationParams(pageLimitStr, 5)

	dashboardList, err := domain.Subscription.DashboardSubscription(serviceName, pageNumber*pageLimit, pageLimit)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, err.Error())
	}

	meta := types.JSON{
		"page_number": pageNumber,
		"page_limit":  pageLimit,
	}

	if len(dashboardList) < 1 {
		meta["total_count"] = 0
		return api.Send(w, http.StatusOK, []any{}, meta)
	}

	resp := lo.Map(dashboardList, func(item schema.DashboardSubscription, _ int) dashboardListResp {
		return dashboardListResp{
			IncidentID:       uint(item.IncidentID.Int64),
			ServiceID:        item.ServiceID,
			ServiceName:      item.ServiceName,
			SubscriptionUUID: item.SubscriptionUUID.String(),
			IncidentName:     item.IncidentName.String,
			IncidentLink:     item.IncidentLink.String,
			IncidentImpact:   item.IncidentImpact.String,
			IsDown:           item.IsDown,
		}
	})

	meta["total_count"] = dashboardList[0].SubscriptionsCount
	return api.Send(w, http.StatusOK, resp, meta)
}
