package main

import (
	"net/http"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/0xi6r/Management-Suite/GoodLife/backend/internal"
	"github.com/0xi6r/Management-Suite/GoodLife/backend/internal/router"
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

	// building the Chi router (internal/router)
	handler := router.New(logger)

	addr := ":" + cfg.ServerPort
	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
