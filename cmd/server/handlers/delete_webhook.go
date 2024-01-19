package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

func DeleteWebhookExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	webhookUUIDStr := chi.URLParam(r, "uuid")

	webhookUUID, err := uuid.Parse(webhookUUIDStr)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot parse uuid %v", err.Error())
	}

	err = domain.WebhookExtension.Delete(webhookUUID)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot delete webhook for uuid %v, err: %v", webhookUUIDStr, err.Error())
	}

	return api.Send(w, http.StatusOK, types.JSON{
		"msg": "webhook sucessfully deleted",
	}, nil)
}
