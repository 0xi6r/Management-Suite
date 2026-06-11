package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	DatabaseURL string
	JWTSecret string
}

func Load() (*Config, error) {
	// viper automatically reads environment variables
	viper.AutomaticEnv()

	// default values
	viper.SetDefault("SERVER_PORT", "8181")
	viper.SetDefault("DATABASE_URL", "postgres://goodlife:goodlife@localhost:5432/goodlife?sslmode=disable")

	cfg := &Config{
		ServerPort: viper.GetString("SERVER_PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
	}
	return cfg, nil
}
