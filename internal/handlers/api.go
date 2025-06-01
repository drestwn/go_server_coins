package handlers

import (
    "go_server/internal/middleware"
    "net/http"
    "github.com/go-chi/chi"
    chimiddle "github.com/go-chi/chi/middleware"
)

type CoinBalanceParams struct {
    Username string `schema:"username"`
}

type CoinBalanceResponse struct {
    Balance int64 `json:"Balance"`
    Code    int   `json:"Code"`
}

func InternalErrorHandler(w http.ResponseWriter) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("Internal server error"))
}

func Handler(r *chi.Mux) {
    // Global middleware
    r.Use(chimiddle.StripSlashes)
	r.Route("/account", func(router chi.Router) {
        // Public endpoint (no authorization)
        router.Post("/coins", CreateCoinBalance)

        // Protected endpoint (requires authorization)
        router.With(middleware.Authorization).Get("/coins", GetCoinBalance)
    })
}