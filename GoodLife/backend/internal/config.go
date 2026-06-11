package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
}

func Load() (*Config, error) {
	// viper automatically reads environment variables
	viper.AutomaticEnv()

	// default values
	viper.SetDefault("SERVER_PORT", "8080")

	cfg := &Config{
		ServerPort: viper.GetString("SERVER_PORT"),
	}
	return cfg, nil
}
