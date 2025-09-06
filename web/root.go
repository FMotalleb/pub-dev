package web

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func init() {
	RegisterEndpoint(
		func(_ context.Context, web *echo.Echo) {
			web.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "Ok!")
			})
		},
	)
}
