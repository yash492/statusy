package handlers

import (
	"net/http"

	"github.com/yash492/statusy/pkg/api"
)

// All subscriptions in a table as of now, accrues into a single dashboard
func DashboardList(w http.ResponseWriter, r *http.Request) *api.Response {

	
	return api.Send(w, http.StatusOK, nil, nil)
}
