package config

import (
	// Autoload .env file.
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	HTTPListenAddr string `mapstructure:"http_listen" env:"HTTP_LISTEN"`

	BaseURL string `mapstructure:"base_url" env:"BASE_URL"`

	PubStorage string `mapstructure:"pub_storage" default:"./storage/pub"`

	Auth []AuthRule `mapstructure:"auth"`
}

type AuthRule struct {
	BasePath []string `mapstructure:"path"`
	Tokens   []string `mapstructure:"token"`
}
