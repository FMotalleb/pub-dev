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
		p := path.Join(storage, p.Name())
		WritePackageMeta(ctx, l, p)
	}
	return nil
}

func WritePackageMeta(ctx context.Context, l *zap.Logger, p string) {
	pl := l.With(zap.String("package", p))
	meta, err := buildMetaData(ctx, p)
	if err != nil {
		pl.Error("failed to generate metadata for package, skipping", zap.Error(err))
		return
	}
	targetPath := path.Join(p, "listing.json")
	target, err := os.Create(targetPath)
	if err != nil {
		pl.Error("failed to create listing.json for package", zap.Error(err))
		return
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		pl.Error("failed to convert hashmap to json for package", zap.Error(err))
		return
	}
	_, err = target.Write(metaJSON)
	if err != nil {
		pl.Error("failed to write json to list file for package", zap.Error(err))
		return
	}
}

func buildMetaData(ctx context.Context, pDir string) (map[string]any, error) {
	l := log.Of(ctx).
		Named("buildMetadata").
		With(zap.String("package", pDir))
	var directories []os.DirEntry
	var err error
	if directories, err = os.ReadDir(pDir); err != nil {
		l.Error("failed to read package directory", zap.Error(err))
		return nil, err
	}
	raw := make(map[string]any, 0)
	raw["name"] = path.Base(pDir)
	versions := make([]map[string]any, 0, len(directories))
	var latest map[string]any
	for _, d := range directories {
		if !d.IsDir() {
			continue
		}
		data := make(map[string]any, 0)
		verData := path.Join(pDir, d.Name(), "package.json")
		if err := utils.ReadJSON(verData, &data); err != nil {
			l.Error("failed to read package version data", zap.Error(err), zap.String("version", d.Name()))
			return nil, err
		}
		if latest == nil {
			latest = data
		}
		latest = newer(latest, data)
		versions = append(versions, data)
	}
	raw["versions"] = versions
	raw["latest"] = latest
	return raw, nil
}

func newer(i, j map[string]any) map[string]any {
	iTime := getPublishTime(i)
	jTime := getPublishTime(j)
	if iTime.After(jTime) {
		return i
	}
	return j
}

func getPublishTime(meta map[string]any) time.Time {
	defTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	published, ok := meta["published"]
	if !ok {
		return defTime
	}
	pubStr, ok := published.(string)
	if !ok {
		return defTime
	}
	result, err := time.Parse(time.RFC3339, pubStr)
	if err != nil {
		return defTime
	}
	return result
}
