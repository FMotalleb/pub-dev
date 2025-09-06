package packages

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/fmotalleb/pub-dev/config"
)

func init() {
	RegisterEndpoint(func(g *echo.Group) {
		// TODO: add authorization middleware
		g.GET(":package", getPackageInfo)
	})
}

func getPackageInfo(ctx echo.Context) error {
	cfg := config.GetForce(ctx.Request().Context())
	var err error
	var directories []os.DirEntry
	p := ctx.Param("package")
	p = strings.Trim(p, "/ :-")
	if p == "" {
		return ctx.String(http.StatusNotFound, "package name is empty")
	}
	pDir := path.Join(cfg.StoragePath, p)
	if directories, err = os.ReadDir(pDir); err != nil {
		return ctx.String(http.StatusNotFound, "package not found")
	}
	raw := make(map[string]any, 3)
	raw["name"] = p
	versions := make([]map[string]any, len(directories))
	for i, d := range directories {
		verData := path.Join(pDir, d.Name(), "package.json")
		if err := readToJson(verData, &versions[i]); err != nil {
			return ctx.String(http.StatusInternalServerError, "package version issue found")
		}
	}
	raw["versions"] = versions
	raw["latest"] = versions[len(versions)-1]
	// raw["latest"]
	return ctx.JSON(http.StatusOK, raw)
}

func readToJson(path string, dst *map[string]any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, dst)
	return err
}
