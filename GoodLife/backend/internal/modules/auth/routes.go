package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/0xi6r/Management-Suite/GoodLife/backend/internal/middleware"
)

func RegisterRoutes(r chi.Router, pool *pgxpool.Pool, jwtSecret string, logger *zap.Logger) {
	h := NewHandler(pool, jwtSecret, logger)
	r.Post("/login", h.Login)

	// Protected group: requires valid JWT
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(jwtSecret, logger))
		r.Get("/me", h.Me)
	})
}
