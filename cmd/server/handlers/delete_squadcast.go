package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

func DeleteSquadcastExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	squadcastUUIDStr := chi.URLParam(r, "uuid")

	squadcastUUID, err := uuid.Parse(squadcastUUIDStr)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot parse uuid %v", err.Error())
	}

	err = domain.SquadcastExtension.Delete(squadcastUUID)
	if err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot delete squadcast for uuid %v, err: %v", squadcastUUIDStr, err.Error())
	}

	return api.Send(w, http.StatusOK, types.JSON{
		"msg": "squadcast sucessfully deleted",
	}, nil)
}
