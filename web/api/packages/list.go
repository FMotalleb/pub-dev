package packages

import (
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/fmotalleb/go-tools/log"

	"github.com/fmotalleb/pub-dev/config"
	"github.com/fmotalleb/pub-dev/utils"
)

func init() {
	RegisterEndpoint(func(g *echo.Group) {
		// TODO: add authorization middleware
		g.GET(":package", getPackageInfo)
	})
}

func getPackageInfo(ctx echo.Context) error {
	cfg := config.GetForce(ctx.Request().Context())
	l := log.Of(ctx.Request().Context())
	p := ctx.Param("package")
	p = strings.Trim(p, "/ :-")
	if p == "" {
		return ctx.String(http.StatusNotFound, "package name is empty")
	}
	listing := path.Join(cfg.StoragePath, p, "listing.json")

	var raw ListingResponse
	err := utils.ReadJSONTemplate(listing, &raw, *cfg)
	if err != nil {
		l.Error("failed to parse json for package", zap.String("path", listing), zap.Error(err))
		return ctx.String(http.StatusInternalServerError, "internal server error")
	}
	return ctx.JSON(http.StatusOK, raw)
}

type (
	ListingResponse struct {
		Name     string        `json:"name"`
		Latest   ListingItem   `json:"latest"`
		Versions []ListingItem `json:"versions"`
	}
	ListingItem struct {
		ArchiveSHA256 string         `json:"archive_sha256"`
		ArchiveURL    string         `json:"archive_url"`
		PublishDate   string         `json:"published"`
		Version       string         `json:"version"`
		PubSpec       map[string]any `json:"pubspec"`
	}
)
