package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

type saveSubscriptionReq struct {
	ServiceID        uint   `json:"service_id"`
	IsAllComponents  bool   `json:"is_all_components"`
	CustomComponents []uint `json:"custom_components"`
}

func (s *saveSubscriptionReq) Validate() error {
	if !s.IsAllComponents && len(s.CustomComponents) == 0 {
		return fmt.Errorf("%w: please choose components for the chosen option", api.ErrValidation)
	}

	return nil
}

func AddSubscription(w http.ResponseWriter, r *http.Request) *api.Response {
	var req saveSubscriptionReq
	if err := api.Decode(r, &req); err != nil {
		if errors.Is(err, api.ErrValidation) {
			return api.Errorf(w, http.StatusBadRequest, err.Error())
		}
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
