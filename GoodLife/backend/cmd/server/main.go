package main

import (
	"net/http"
	"context"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/0xi6r/Management-Suite/GoodLife/backend/internal"
	"github.com/0xi6r/Management-Suite/GoodLife/backend/internal/router"
	"github.com/0xi6r/Management-Suite/GoodLife/backend/internal/db"
)

func main() {
	// Load .env file if present (local dev)
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// a structured logger (dev mode)
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger) // makes zap.L() available everywhere

	logger.Info("Goodlife API starting", zap.String("port", cfg.ServerPort),)

	// connect to database
	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("database conn failed", zap.Error(err))
	}
	defer pool.Close()
	logger.Info("db conn established")

	// building the Chi router (internal/router)
	handler := router.New(logger, pool)

	addr := ":" + cfg.ServerPort
	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
