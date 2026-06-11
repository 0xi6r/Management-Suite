package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func RegisterRoutes(r chi.Router, pool *pgxpool.Pool, jwtSecret string, logger *zap.Logger) {
	h := NewHandler(pool, jwtSecret, logger)
	r.Post("/login", h.Login)
}
