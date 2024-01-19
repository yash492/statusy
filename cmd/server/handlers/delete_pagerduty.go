package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

func DeletePagerdutyExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	pagerdutyUUIDStr := chi.URLParam(r, "uuid")

	pagerdutyUUID, err := uuid.Parse(pagerdutyUUIDStr)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot parse uuid %v", err.Error())
	}

	err = domain.PagerdutyExtension.Delete(pagerdutyUUID)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot delete pagerduty for uuid %v, err: %v", pagerdutyUUIDStr, err.Error())
	}

	return api.Send(w, http.StatusOK, types.JSON{
		"msg": "pagerduty sucessfully deleted",
	}, nil)
}
