package config

import (
	// Autoload .env file.
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	HTTPListenAddr string `mapstructure:"http_listen" env:"HTTP_LISTEN"`

	BaseURL string `mapstructure:"base_url" env:"BASE_URL"`

	DatabaseConnection string `mapstructure:"database" env:"DATABASE" default:"sqlite://storage/packages.db" validate:"required,uri"`
	StoragePath        string `mapstructure:"storage" default:"./storage/packages"`
}
