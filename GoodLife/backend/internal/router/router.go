package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	mw "github.com/0xi6r/Management-Suite/GoodLife/backend/internal/middleware"
)

func New(logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	// Chi's built-in middlewares
	//r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	//r.Use(middleware.Logger) // logs every request (standard lib logger)
	//r.Use(middleware.Recoverer)

	// custom zap-based req logger (replaces middleware.Logger)
	r.Use(mw.RequestLogger(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return r
}
