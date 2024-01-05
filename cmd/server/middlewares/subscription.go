package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/types"
)

func Subscription(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subscriptionID := chi.URLParam(r, "subscriptionID")

		uuid, err := uuid.Parse(subscriptionID)
		if err != nil {
			api.MiddlewareErrorf(w, http.StatusInternalServerError, "cannot parse uuid")
		}

		ctx := ContextWrapAll(r.Context(), map[string]any{
			types.SubscriptionIDCtx: uuid,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
