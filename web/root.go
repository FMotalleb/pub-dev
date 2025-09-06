package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func init() {
	RegisterEndpoint(
		func(web *echo.Echo) {
			web.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "Ok!")
			})
		},
	)
}
