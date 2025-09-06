package web

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/fmotalleb/pub-dev/config"
)

var (
	routes = make([]func(context.Context, *echo.Echo), 0)
	server = echo.New()
)

func RegisterEndpoint(register func(context.Context, *echo.Echo)) {
	routes = append(routes, register)
}

func init() {
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
}

func Start(ctx context.Context) error {
	cfg, err := config.Get(ctx)
	if err != nil {
		return err
	}
	if cfg.HTTPListenAddr == "" {
		return errors.New("`http_listen` is not set")
	}
	server.Server = &http.Server{
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		IdleTimeout:       time.Minute,
		WriteTimeout:      time.Minute,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	for _, r := range routes {
		r(ctx, server)
	}
	if err := server.Start(cfg.HTTPListenAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
