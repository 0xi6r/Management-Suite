package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/0xi6r/Management-Suite/GoodLife/internal/config"
)

func main() {
	// Load .env file if present (local dev)
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	addr := ":" + cfg.ServerPort
	fmt.Printf("GoodLife API starting on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
