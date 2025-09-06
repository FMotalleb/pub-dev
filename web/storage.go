package web

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/fmotalleb/pub-dev/config"
)

func init() {
	RegisterEndpoint(
		func(ctx context.Context, web *echo.Echo) {
			cfg := config.GetForce(ctx)
			web.Static("storage/packages", cfg.StoragePath)
		},
	)
}
