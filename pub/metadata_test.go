package pub

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fmotalleb/pub-dev/utils"
)

func TestReadPackage(t *testing.T) {
	tempDir := t.TempDir()
	pkgPath := filepath.Join(tempDir, "test-package")
	require.NoError(t, os.Mkdir(pkgPath, 0o755))

	// Create dummy versions
	v1Path := filepath.Join(pkgPath, "1.0.0")
	require.NoError(t, os.Mkdir(v1Path, 0o755))
	v1 := &PackageVersion{
		Version:   "1.0.0",
		Published: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	require.NoError(t, utils.WriteJSON(filepath.Join(v1Path, "package.json"), v1))

	v2Path := filepath.Join(pkgPath, "2.0.0")
	require.NoError(t, os.Mkdir(v2Path, 0o755))
	v2 := &PackageVersion{
		Version:   "2.0.0",
		Published: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	require.NoError(t, utils.WriteJSON(filepath.Join(v2Path, "package.json"), v2))

	pkg, err := ReadPackage(t.Context(), pkgPath)
	require.NoError(t, err)

	assert.Equal(t, "test-package", pkg.Name)
	assert.Len(t, pkg.Versions, 2)
	assert.Equal(t, v2, pkg.Latest)
}

func TestPackage_WriteMeta(t *testing.T) {
	tempDir := t.TempDir()
	pkgPath := filepath.Join(tempDir, "test-package")
	require.NoError(t, os.Mkdir(pkgPath, 0o755))

	pkg := &Package{
		Name: "test-package",
		Latest: &PackageVersion{
			Version:   "1.0.0",
			Published: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Versions: []*PackageVersion{
			{
				Version:   "1.0.0",
				Published: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	err := pkg.WriteMeta(pkgPath)
	require.NoError(t, err)

	metaPath := filepath.Join(pkgPath, "listing.json")
	_, err = os.Stat(metaPath)
	assert.NoError(t, err)

	readPkg := &Package{}
	require.NoError(t, utils.ReadJSON(metaPath, readPkg))
	assert.Equal(t, pkg, readPkg)
}
