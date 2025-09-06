package config

import (
	// Autoload .env file.
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	HTTPListenAddr string `mapstructure:"http_listen" env:"HTTP_LISTEN"`

	DatabaseConnection string `mapstructure:"database" env:"DATABASE" default:"sqlite://storage/packages.db" validate:"required,uri"`
	StoragePath        string `mapstructure:"storage" default:"storage/packages"`
}
