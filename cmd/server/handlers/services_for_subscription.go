package handlers

import (
	"net/http"
	"strings"

	"github.com/samber/lo"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/schema"
)

type servicesForSubsciptionsResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ServicesForSubsciptions(w http.ResponseWriter, r *http.Request) *api.Response {
	query := r.URL.Query().Get("query")
	serviceNameQuery := strings.TrimSpace(query)

	services, err := domain.Subscription.GetAllServicesForSubscriptions(serviceNameQuery)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "could not fetch services for %v query %v", serviceNameQuery, err.Error())
	}

	formattedService := lo.Map(services, func(service schema.ServicesForSubsciption, _ int) servicesForSubsciptionsResponse {
		return servicesForSubsciptionsResponse{
			ID:   service.ServiceID,
			Name: service.ServiceName,
		}
	})

	return api.Send(w, http.StatusOK, formattedService, nil)
}
