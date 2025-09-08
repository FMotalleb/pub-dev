package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/fmotalleb/pub-dev/config"
)

func Middleware(rules []config.AuthRule) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
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
				return validate(r.Tokens, c)
			}
		}
	}
	return true
}

func validate(tokens []string, c echo.Context) bool {
	head := getBearer(c)
	for _, t := range tokens {
		if getMatcher(t)(head) {
			return true
		}
	}
	return false
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

const authHeaderParts = 2

func getMatcher(matcher string) func(string) bool {
	parts := strings.SplitN(matcher, ":", authHeaderParts)
	if len(parts) != authHeaderParts {
		// No type prefix, fallback to simple equality
		return func(in string) bool {
			return matcher == in
		}
	}

	hashType, hashValue := parts[0], parts[1]

	switch hashType {
	case "sha256":
		return func(in string) bool {
			sum := sha256.Sum256([]byte(in))
			return hex.EncodeToString(sum[:]) == hashValue
		}
	case "bcrypt":
		return func(in string) bool {
			err := bcrypt.CompareHashAndPassword([]byte(hashValue), []byte(in))
			return err == nil
		}
	default:
		// Unknown type, fallback to string equality
		return func(in string) bool {
			return matcher == in
		}
	}
}
