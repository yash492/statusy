package middlewares

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/types"
)

func Service(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		serviceIDStr := chi.URLParam(r, "serviceID")
		serviceID, err := strconv.Atoi(serviceIDStr)
		if err != nil {
			api.MiddlewareErrorf(w, http.StatusInternalServerError, "cannot decode %v service ID", serviceIDStr)
			return
		}

		ctx := ContextWrapAll(r.Context(), map[string]any{
			types.ServiceIDCtx: uint(serviceID),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
