package handlers

import (
	"go_server/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(r *chi.Mux) {
	// global middleware
	r.Use(chimiddle.StripSlashes) // Ignores trailing slashes, e.g., "google.com/a/b/" ignores slash after b

	r.Route("/account", func(router chi.Router) {
		// middleware for /account route
		router.Use(middleware.Authorization)
		router.Get("/coins", GetCoinBalance)
	})
}
