package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

func DeleteSubscription(w http.ResponseWriter, r *http.Request) *api.Response {
	ctx := r.Context()
	subscriptionID := ctx.Value(types.SubscriptionIDCtx).(uuid.UUID)

	if err := domain.Subscription.Delete(subscriptionID); err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot delete subscription for %v, err: %v", subscriptionID.String(), err.Error())
	}

	return api.Send(w, http.StatusOK, types.JSON{
		"msg": "subscription successfully deleted",
	}, nil)

}
