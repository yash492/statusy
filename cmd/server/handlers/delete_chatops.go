package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

func DeleteChatopsExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	chatopsUUIDStr := chi.URLParam(r, "uuid")

	chatopsType := r.URL.Query().Get("type")

	chatopsUUID, err := uuid.Parse(chatopsUUIDStr)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot parse uuid %v", err.Error())
	}

	err = domain.ChatopsExtension.Delete(chatopsUUID)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot delete %v for uuid %v, err: %v", chatopsType, chatopsUUIDStr, err.Error())
	}

	return api.Send(w, http.StatusOK, types.JSON{
		"msg": "chatops extension is sucessfully deleted",
	}, nil)
}
