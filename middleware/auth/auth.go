package auth

import (
	"net/http"
	"slices"
	"strings"

	"github.com/fmotalleb/pub-dev/config"
	"github.com/labstack/echo/v4"
)

func Middleware(rules []config.AuthRule) echo.MiddlewareFunc {
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
