package api

import (
	"github.com/labstack/echo/v4"
)

var setup []func(*echo.Group)

func RegisterEndpoint(register func(*echo.Group)) {
	setup = append(setup, register)
}

func Setup(web *echo.Echo) {
	api := web.Group("api/")
	for _, e := range setup {
		e(api)
	}
}
