package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmotalleb/pub-dev/config"
	"github.com/fmotalleb/pub-dev/middleware/auth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		wantStatus int
		path       string
		header     string
		rules      []config.AuthRule
	}{
		{
			name:       "No auth needed",
			wantStatus: http.StatusOK,
			path:       "/api/packages/mypackage",
		},
		{
			name:       "Auth needed, but no token provided",
			wantStatus: http.StatusUnauthorized,
			path:       "/api/packages/versions/newUpload",
			rules: []config.AuthRule{
				{
					BasePath: []string{"/api/packages/versions/newUpload"},
					Tokens:   []string{"my-secret-token"},
				},
			},
		},
		{
			name:       "Auth needed, with invalid token",
			wantStatus: http.StatusUnauthorized,
			path:       "/api/packages/versions/newUpload",
			header:     "Bearer invalid-token",
			rules: []config.AuthRule{
				{
					BasePath: []string{"/api/packages/versions/newUpload"},
					Tokens:   []string{"my-secret-token"},
				},
			},
		},
		{
			name:       "Auth needed, with valid token",
			wantStatus: http.StatusOK,
			path:       "/api/packages/versions/newUpload",
			header:     "Bearer my-secret-token",
			rules: []config.AuthRule{
				{
					BasePath: []string{"/api/packages/versions/newUpload"},
					Tokens:   []string{"my-secret-token"},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Authorization", tt.header)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := auth.Middleware(tt.rules)(func(c echo.Context) error {
				return c.String(http.StatusOK, "next")
			})
			if assert.NoError(t, h(c)) {
				assert.Equal(t, tt.wantStatus, rec.Code)
			}
		})
	}
}
