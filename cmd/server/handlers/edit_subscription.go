package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

func EditSubscription(w http.ResponseWriter, r *http.Request) *api.Response {
	var req saveSubscriptionReq
	ctx := r.Context()

	subscriptionID := ctx.Value(types.SubscriptionIDCtx).(uuid.UUID)

	if err := api.Decode(r, &req); err != nil {
		if errors.Is(err, api.ErrValidation) {
			return api.Errorf(w, http.StatusBadRequest, err.Error())
		}
		return api.Errorf(w, http.StatusUnprocessableEntity, "cannot process the given json req")
	}

	if err := domain.Subscription.Update(subscriptionID, req.CustomComponents, req.IsAllComponents); err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot update subscription for %v, err: %v", subscriptionID.String(), err.Error())
	}

	return api.Send(w, http.StatusOK, types.JSON{
		"msg": "subscription successfully updated",
	}, nil)

}
