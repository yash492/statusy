package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yash492/statusy/cmd/server/handlers"
	"github.com/yash492/statusy/cmd/server/middlewares"
	"github.com/yash492/statusy/pkg/api"
)

func registerRoutes(r chi.Router) {
	r.Route("/api", routes)
}

func routes(r chi.Router) {
	r.Route("/services", servicesRoutes)
	r.Route("/subscriptions", subscriptionRoutes)

}

func subscriptionRoutes(r chi.Router) {
	r.Method(http.MethodPost, "/", api.Handler(handlers.AddSubscription))
	r.Method(http.MethodGet, "/services", api.Handler(handlers.ServicesForSubsciptions))
	r.With(middlewares.Subscription).Route("/{subscriptionID}", subSubscriptionRoutes)
}

func subSubscriptionRoutes(r chi.Router) {
	r.Method(http.MethodGet, "/", api.Handler(handlers.SubscriptionByID))
	r.Method(http.MethodPut, "/", api.Handler(handlers.EditSubscription))
}

func servicesRoutes(r chi.Router) {
	r.With(middlewares.Service).Route("/{serviceID}", subServiceRoutes)
}

func subServiceRoutes(r chi.Router) {
	r.Method(http.MethodGet, "/components", api.Handler(handlers.ComponentsByService))
}
