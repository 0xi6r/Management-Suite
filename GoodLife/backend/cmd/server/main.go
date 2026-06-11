package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/0xi6r/Management-Suite/GoodLife/backend/internal"
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

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	addr := ":" + cfg.ServerPort
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
