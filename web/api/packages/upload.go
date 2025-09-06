package packages

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fmotalleb/go-tools/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/fmotalleb/pub-dev/config"
	"github.com/fmotalleb/pub-dev/pub"
)

const directoryMakePermission = 0o755

var tempDirRoot = os.TempDir()

func init() {
	RegisterEndpoint(func(g *echo.Group) {
		// TODO: add authorization middleware
		g.GET("versions/new", handleNewUpload)
		g.POST("versions/newUpload", handleUpload)
		g.GET("versions/newUploadFinish", handleFinalize)
	})
}

type NewUploadResponse struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

func handleNewUpload(ctx echo.Context) error {
	cfg := config.GetForce(ctx.Request().Context())
	return ctx.JSON(http.StatusOK, NewUploadResponse{
		URL:    cfg.BaseURL + "api/packages/versions/newUpload",
		Fields: map[string]string{},
	})
}

func handleUpload(c echo.Context) error {
	l := log.Of(c.Request().Context()).Named("upload")

	randID := generateTempID()
	tempDir := filepath.Join(tempDirRoot, randID)

	if err := os.MkdirAll(tempDir, directoryMakePermission); err != nil {
		l.Error("failed to create temp directory", zap.Error(err), zap.String("path", tempDir))
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "missing file field")
	}

	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to open uploaded file")
	}
	defer src.Close()

	tempPath := filepath.Join(tempDir, file.Filename)
	if err := saveUploadedFile(src, tempPath); err != nil {
		l.Error("failed to save file", zap.Error(err))
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	redirectURL := buildFinalizeURL(randID, file.Filename)
	cfg := config.GetForce(c.Request().Context())
	return c.Redirect(http.StatusFound, cfg.BaseURL+redirectURL)
}

func handleFinalize(c echo.Context) error {
	l := log.Of(c.Request().Context()).Named("finalize")
	randID := c.QueryParam("temp")
	name := c.QueryParam("name")
	tempPath := filepath.Join(tempDirRoot, randID, name)

	if _, err := os.Stat(tempPath); err != nil {
		return c.String(http.StatusNotFound, "upload not found")
	}
	l.Info("finalizing upload", zap.String("path", tempPath))

	spec, err := pub.ReadSpecFromTar(tempPath)
	if err != nil {
		l.Error("failed to read pubspec", zap.Error(err))
		return c.String(http.StatusBadRequest, "invalid package data")
	}

	cfg := config.GetForce(c.Request().Context())
	finalRoot := filepath.Join(cfg.StoragePath, spec.Name, spec.Version)

	if err := os.MkdirAll(finalRoot, directoryMakePermission); err != nil {
		l.Error("failed to create final directory", zap.Error(err))
		return c.String(http.StatusInternalServerError, "server error")
	}

	finalPath := filepath.Join(finalRoot, "package.tar.gz")
	if err := moveFile(tempPath, finalPath); err != nil {
		l.Error("failed to move package", zap.Error(err))
		return c.String(http.StatusInternalServerError, "server error")
	}
	err = writeSpecData(spec, l, c, finalRoot, finalPath)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": map[string]string{
			"message": "package published successfully",
		},
	})
}

// --- Helpers ---

func generateTempID() string {
	return rand.Text()[0:10]
}

func writeSpecData(spec *pub.Spec, l *zap.Logger, c echo.Context, finalRoot string, finalPath string) error {
	cfg := config.GetForce(c.Request().Context())
	raw := make(map[string]any)
	raw["pubspec"] = spec.Raw
	raw["version"] = spec.Version
	raw["archive_url"] = cfg.BaseURL + path.Join("storage", "packages", spec.Name+"-"+spec.Version+".tar.gz")
	data, err := os.ReadFile(finalPath)
	if err != nil {
		return err
	}
	hash := sha256.Sum256(data)
	hashString := hex.EncodeToString(hash[:])
	raw["archive_sha256"] = hashString
	raw["published"] = time.Now().Format(time.RFC3339)
	specData, err := json.Marshal(raw)
	if err != nil {
		l.Error("failed to marshal spec data", zap.Error(err))
		return c.String(http.StatusInternalServerError, "server error")
	}

	f, err := os.Create(path.Join(finalRoot, "package.json"))
	if err != nil {
		l.Error("failed to open spec data file", zap.Error(err))
		return c.String(http.StatusInternalServerError, "server error")
	}
	_, err = f.Write(specData)
	if err != nil {
		l.Error("failed to write spec data file", zap.Error(err))
		return c.String(http.StatusInternalServerError, "server error")
	}
	return nil
}

func buildFinalizeURL(tempID, filename string) string {
	u := &url.URL{Path: "api/packages/versions/newUploadFinish"}
	q := u.Query()
	q.Set("temp", tempID)
	q.Set("name", filename)
	u.RawQuery = q.Encode()
	return u.String()
}

func saveUploadedFile(src io.Reader, dstPath string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func moveFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		src.Close()
		os.Remove(srcPath)
	}()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return os.Remove(srcPath)
}
