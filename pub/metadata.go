package pub

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/fmotalleb/go-tools/log"
	"go.uber.org/zap"

	"github.com/fmotalleb/pub-dev/utils"
)

const filePermission = 0o600

type Package struct {
	Name     string            `json:"name"`
	Latest   *PackageVersion   `json:"latest"`
	Versions []*PackageVersion `json:"versions"`
}

type PackageVersion struct {
	Version       string         `json:"version"`
	Pubspec       map[string]any `json:"pubspec"`
	ArchiveURL    string         `json:"archive_url"`
	ArchiveSHA256 string         `json:"archive_sha256"`
	Published     time.Time      `json:"published"`
}

func RecalculateMetadata(ctx context.Context, storage string) error {
	l := log.Of(ctx).
		Named("RecalculateMetadata").
		With(zap.String("storage", storage))
	ctx = log.WithLogger(ctx, l)
	packages, err := os.ReadDir(storage)
	if err != nil {
		l.Error("failed to read storage directory", zap.Error(err))
		return err
	}
	for _, p := range packages {
		if !p.IsDir() {
			continue
		}
		pPath := path.Join(storage, p.Name())
		pkg, err := ReadPackage(ctx, pPath)
		if err != nil {
			l.Error("failed to read package", zap.Error(err), zap.String("package", p.Name()))
			continue
		}
		if err := pkg.WriteMeta(pPath); err != nil {
			l.Error("failed to write package metadata", zap.Error(err), zap.String("package", p.Name()))
		}
	}
	return nil
}

func ReadPackage(ctx context.Context, pDir string) (*Package, error) {
	l := log.Of(ctx).
		Named("ReadPackage").
		With(zap.String("package", pDir))
	versions, err := os.ReadDir(pDir)
	if err != nil {
		l.Error("failed to read package directory", zap.Error(err))
		return nil, err
	}

	pkg := &Package{
		Name:     path.Base(pDir),
		Versions: make([]*PackageVersion, 0, len(versions)),
	}

	for _, v := range versions {
		if !v.IsDir() {
			continue
		}
		verPath := path.Join(pDir, v.Name(), "package.json")
		var ver PackageVersion
		if err := utils.ReadJSON(verPath, &ver); err != nil {
			l.Error("failed to read package version data", zap.Error(err), zap.String("version", v.Name()))
			continue
		}
		pkg.AddVersion(&ver)
	}

	return pkg, nil
}

func (p *Package) AddVersion(v *PackageVersion) {
	p.Versions = append(p.Versions, v)
	if p.Latest == nil || v.Published.After(p.Latest.Published) {
		p.Latest = v
	}
}

func (p *Package) WriteMeta(pDir string) error {
	metaJSON, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	targetPath := path.Join(pDir, "listing.json")
	return os.WriteFile(targetPath, metaJSON, filePermission)
}
