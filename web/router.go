package web

import (
	"context"
	"net/http"

	"github.com/fmotalleb/pub-dev/config"
	"github.com/fmotalleb/pub-dev/web/handlers"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(ctx context.Context, e *echo.Echo) {
	cfg := config.GetForce(ctx)

	// Root and static routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok!")
	})
	e.Static("storage/packages", cfg.StoragePath)

	// API routes
	api := e.Group("/api")
	p := api.Group("/packages")

	// Package routes
	p.GET("/:package", handlers.GetPackageInfo)
	p.GET("/versions/new", handlers.HandleNewUpload)
	p.POST("/versions/newUpload", handlers.HandleUpload)
	p.GET("/versions/newUploadFinish", handlers.HandleFinalize)
}
