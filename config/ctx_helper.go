package config

import (
	"context"
	"errors"
)

type CtxKey string

const (
	configKey = CtxKey("config")
)

var (
	ErrorNoConfigFound   = errors.New("no config attached to the current context")
	ErrorConfigMissMatch = errors.New("config was removed or replaced, key is defined but type miss match")
)

func Attach(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

func Get(ctx context.Context) (*Config, error) {
	var v any
	if v = ctx.Value(configKey); v == nil {
		return nil, ErrorNoConfigFound
	}
	if cfg, ok := v.(*Config); ok {
		return cfg, nil
	}
	return nil, ErrorConfigMissMatch
}
