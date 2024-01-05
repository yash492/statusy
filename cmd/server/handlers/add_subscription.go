package handlers

import (
	"net/http"

	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

type addSubscriptionReq struct {
	ServiceID        uint   `json:"service_id"`
	IsAllComponents  bool   `json:"is_all_components"`
	CustomComponents []uint `json:"custom_components"`
}

func AddSubscription(w http.ResponseWriter, r *http.Request) *api.Response {
	var req addSubscriptionReq
	if err := api.Decode(r, &req); err != nil {
		return api.Errorf(w, http.StatusUnprocessableEntity, "cannot process the given json req")
	}

	err := domain.Subscription.Create(req.ServiceID, req.CustomComponents, req.IsAllComponents)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot create subscription for %v service", req.ServiceID)
	}

	return api.Send(w, http.StatusCreated, types.JSON{
		"msg": "subscription successfully created",
	}, nil)
}
