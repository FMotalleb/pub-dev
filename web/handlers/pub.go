package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fmotalleb/go-tools/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/fmotalleb/pub-dev/config"
	"github.com/fmotalleb/pub-dev/pub"
	"github.com/fmotalleb/pub-dev/utils"
)

const (
	directoryMakePermission = 0o755
	tempRandSize            = 10
)

var tempDirRoot = os.TempDir()

type NewUploadResponse struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

func GetPackageInfo(ctx echo.Context) error {
	cfg := config.GetForce(ctx.Request().Context())
	l := log.Of(ctx.Request().Context())
	p := ctx.Param("package")
	p = strings.Trim(p, "/ :-")
	if p == "" {
		return ctx.String(http.StatusNotFound, "package name is empty")
	}
	listing := path.Join(cfg.PubStorage, p, "listing.json")

	var raw pub.Package
	err := utils.ReadJSONTemplate(listing, &raw, *cfg)
	if err != nil {
		if os.IsNotExist(err) {
			return ctx.NoContent(http.StatusOK)
		}
		l.Error("failed to parse json for package", zap.String("path", listing), zap.Error(err))
		return ctx.String(http.StatusInternalServerError, "internal server error")
	}
	return ctx.JSON(http.StatusOK, raw)
}

func HandleNewUpload(ctx echo.Context) error {
	cfg := config.GetForce(ctx.Request().Context())
	return ctx.JSON(http.StatusOK, NewUploadResponse{
		URL:    cfg.BaseURL + "api/packages/versions/newUpload",
		Fields: map[string]string{},
	})
}

func HandleUpload(c echo.Context) error {
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

func HandleFinalize(c echo.Context) error {
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
	packageDir := filepath.Join(cfg.PubStorage, spec.Name)
	finalDir := filepath.Join(packageDir, spec.Version)

	if err = os.MkdirAll(finalDir, directoryMakePermission); err != nil {
		l.Error("failed to create final directory", zap.Error(err))
		return c.String(http.StatusInternalServerError, "server error")
	}

	finalPath := filepath.Join(finalDir, "package.tar.gz")
	if err = moveFile(tempPath, finalPath); err != nil {
		l.Error("failed to move package", zap.Error(err))
		return c.String(http.StatusInternalServerError, "server error")
	}

	version, err := createPackageVersion(spec, finalPath)
	if err != nil {
		l.Error("failed to create listing item", zap.Error(err))
		return c.String(http.StatusInternalServerError, "server error")
	}

	if err = utils.WriteJSON(filepath.Join(finalDir, "package.json"), version); err != nil {
		l.Error("failed to write spec data", zap.Error(err))
		return c.JSON(http.StatusOK, map[string]any{
			"error": map[string]any{
				"code":    1,
				"message": "failed to generate spec data",
			},
		})
	}

	pkg, err := pub.ReadPackage(c.Request().Context(), packageDir)
	if err != nil {
		l.Error("failed to read package", zap.Error(err))
		return c.String(http.StatusInternalServerError, "server error")
	}

	if err := pkg.WriteMeta(packageDir); err != nil {
		l.Error("failed to write package metadata", zap.Error(err))
		return c.String(http.StatusInternalServerError, "server error")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": map[string]string{
			"message": "package published successfully",
		},
	})
}

// --- Helpers ---

func generateTempID() string {
	b := make([]byte, tempRandSize)
	if _, err := rand.Read(b); err != nil {
		return "random"
	}
	return hex.EncodeToString(b)
}

func createPackageVersion(spec *pub.Spec, packagePath string) (*pub.PackageVersion, error) {
	hash, err := calculateSHA256(packagePath)
	if err != nil {
		return nil, err
	}

	return &pub.PackageVersion{
		Pubspec:       spec.Raw,
		Version:       spec.Version,
		ArchiveURL:    "{{ .BaseURL }}" + path.Join("storage", "packages", spec.Name, spec.Version, "package.tar.gz"),
		ArchiveSHA256: hash,
		Published:     time.Now(),
	}, nil
}

func calculateSHA256(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
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

// copy and remove old file because the move is not working on cross device.
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
	return os.RemoveAll(filepath.Dir(srcPath))
}
