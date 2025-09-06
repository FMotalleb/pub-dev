package config

import (
	"time"

	// Autoload .env file.
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	HTTPListenAddr string `mapstructure:"http_listen" env:"HTTP_LISTEN"`

	DatabaseConnection string        `mapstructure:"database" env:"DATABASE" default:"sqlite://data/packages.db" validate:"required,uri"`
	StoragePath        string        `mapstructure:"storage" default:"./packages"`
	RotateAfter        time.Duration `mapstructure:"max_age" env:"MAX_AGE" default:"1h"`
}
