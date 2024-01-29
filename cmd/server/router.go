package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yash492/statusy/cmd/server/handlers"
	"github.com/yash492/statusy/cmd/server/middlewares"
	"github.com/yash492/statusy/pkg/api"
)

func registerRoutes(r chi.Router) {
	r.Get("/*", handlers.ServeFrontendHandler)
	r.Route("/api", routes)
}

func routes(r chi.Router) {
	r.Method(http.MethodGet, "/dashboard", api.Handler(handlers.DashboardList))

	r.Route("/services", servicesRoutes)
	r.Route("/subscriptions", subscriptionRoutes)
	r.Route("/integrations", integrationRoutes)

}

func subscriptionRoutes(r chi.Router) {
	r.Method(http.MethodPost, "/", api.Handler(handlers.AddSubscription))
	r.Method(http.MethodGet, "/services", api.Handler(handlers.ServicesForSubsciptions))
	r.With(middlewares.Subscription).Route("/{subscriptionID}", func(r chi.Router) {
		r.Method(http.MethodDelete, "/", api.Handler(handlers.DeleteSubscription))
		r.Method(http.MethodGet, "/", api.Handler(handlers.SubscriptionByID))
		r.Method(http.MethodPut, "/", api.Handler(handlers.EditSubscription))
		r.Method(http.MethodGet, "/incidents", api.Handler(handlers.SubscriptionIncidents))

	})
}

func servicesRoutes(r chi.Router) {
	r.With(middlewares.Service).Route("/{serviceID}", func(r chi.Router) {
		r.Method(http.MethodGet, "/components", api.Handler(handlers.ComponentsByService))
	})
}

func integrationRoutes(r chi.Router) {
	r.Method(http.MethodPut, "/chatops", api.Handler(handlers.SaveChatOpsExtension))
	r.Method(http.MethodPut, "/webhook", api.Handler(handlers.SaveWebhookExtension))

	r.Method(http.MethodGet, "/chatops", api.Handler(handlers.GetChatopsExtension))
	r.Method(http.MethodGet, "/webhook", api.Handler(handlers.GetWebhookExtension))

	r.Method(http.MethodDelete, "/chatops/{uuid}", api.Handler(handlers.DeleteChatopsExtension))
	r.Method(http.MethodDelete, "/webhook/{uuid}", api.Handler(handlers.DeleteWebhookExtension))

	r.Route("/incident-management", func(r chi.Router) {
		r.Method(http.MethodGet, "/", api.Handler(handlers.GetIncidentManagementExtension))

		r.Method(http.MethodPut, "/squadcast", api.Handler(handlers.SaveSquadcastExtension))
		r.Method(http.MethodPut, "/pagerduty", api.Handler(handlers.SavePagerdutyExtension))

		r.Method(http.MethodDelete, "/squadcast/{uuid}", api.Handler(handlers.DeleteSquadcastExtension))
		r.Method(http.MethodDelete, "/pagerduty/{uuid}", api.Handler(handlers.DeletePagerdutyExtension))
	})
}
