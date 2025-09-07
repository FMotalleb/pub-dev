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
	"github.com/fmotalleb/pub-dev/web/middleware/auth"
)

type Server struct {
	*echo.Echo
	cfg *config.Config
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	if cfg.HTTPListenAddr == "" {
		return nil, errors.New("`http_listen` is not set")
	}
	server := &Server{
		Echo: echo.New(),
		cfg:  cfg,
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

	server.Use(
		middleware.Logger(),
		middleware.Recover(),
		auth.Middleware(cfg.Auth),
	)

	SetupRoutes(ctx, server.Echo, cfg)
	return server, nil
}

func (s *Server) Start() error {
	if err := s.Echo.Start(s.cfg.HTTPListenAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
