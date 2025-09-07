package web

import (
	"context"
	"errors"
	"net"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/fmotalleb/pub-dev/config"
)

func Start(ctx context.Context) error {
	cfg, err := config.Get(ctx)
	if err != nil {
		return err
	}
	if cfg.HTTPListenAddr == "" {
		return errors.New("`http_listen` is not set")
	}

	server := echo.New()
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
		authMiddleware(cfg.Auth),
	)

	SetupRoutes(ctx, server)

	if err := server.Start(cfg.HTTPListenAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func authMiddleware(rules []config.AuthRule) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			if !checkAuth(c, rules, path) {
				c.Response().Header().Add("WWW-Authenticate", `Bearer realm="pub", message="Obtain a token from administrator"`)
				return c.String(http.StatusUnauthorized, "unauthorized")
			}
			return next(c)
		}
	}
}

func checkAuth(c echo.Context, rules []config.AuthRule, path string) bool {
	for _, r := range rules {
		for _, bp := range r.BasePath {
			if strings.HasPrefix(path, bp) {
				return slices.Contains(r.Tokens, getBearer(c))
			}
		}
	}
	return true
}

func getBearer(c echo.Context) string {
	headerParts := 2
	header := c.Request().Header.Get("Authorization")
	head := strings.SplitN(header, " ", headerParts)
	if len(head) != headerParts {
		return ""
	}
	if strings.ToLower(head[0]) != "bearer" {
		return ""
	}
	return head[1]
}
