package web

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/fmotalleb/pub-dev/config"
	"github.com/fmotalleb/pub-dev/web/handlers"
)

func SetupRoutes(ctx context.Context, e *echo.Echo, cfg *config.Config) {

	// Root and static routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok!")
	})
	e.Static("storage/packages", cfg.PubStorage)

	// API routes
	api := e.Group("/api")
	p := api.Group("/packages")

	// Package routes
	p.GET("/:package", handlers.GetPackageInfo)
	p.GET("/versions/new", handlers.HandleNewUpload)
	p.POST("/versions/newUpload", handlers.HandleUpload)
	p.GET("/versions/newUploadFinish", handlers.HandleFinalize)
}
