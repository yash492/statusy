package handlers

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

type componentsByServiceResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ComponentsByService(w http.ResponseWriter, r *http.Request) *api.Response {
	serviceID := r.Context().Value(types.ServiceIDCtx).(uint)
	components, err := domain.Component.GetAllByServiceID(uint(serviceID))
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "could not find components for service %v", serviceID)
	}

	componentsResp := lo.Map(components, func(component schema.Component, _ int) componentsByServiceResp {
		return componentsByServiceResp{
			ID:   component.ID,
			Name: component.Name,
		}
	})

	return api.Send(w, http.StatusOK, componentsResp, nil)
}
