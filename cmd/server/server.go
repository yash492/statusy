package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	registerRoutes(r)

	http.ListenAndServe(":8080", r)
}
