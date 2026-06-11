package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	mw "github.com/0xi6r/Management-Suite/GoodLife/backend/internal/middleware"
)

func New(logger *zap.Logger, pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	// Chi's built-in middlewares
	//r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	//r.Use(middleware.Logger) // logs every request (standard lib logger)
	//r.Use(middleware.Recoverer)
	r.Use(mw.RequestLogger(logger))

	// custom zap-based req logger (replaces middleware.Logger)
	//r.Use(mw.RequestLogger(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("db unavailable"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return r
}
