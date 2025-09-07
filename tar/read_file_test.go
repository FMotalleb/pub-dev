package tar

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestTarGz(t *testing.T, dir string, files map[string]string) string {
	t.Helper()
	tarGzPath := filepath.Join(dir, "test.tar.gz")
	f, err := os.Create(tarGzPath)
	assert.NoError(t, err)
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	tw := tar.NewWriter(gz)
	defer tw.Close()

	for name, content := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0600,
			Size: int64(len(content)),
		}
		err := tw.WriteHeader(hdr)
		assert.NoError(t, err)
		_, err = tw.Write([]byte(content))
		assert.NoError(t, err)
	}

	return tarGzPath
}

func TestReadFile(t *testing.T) {
	tempDir := t.TempDir()
	pubspecContent := "name: test_package\nversion: 1.0.0"
	files := map[string]string{
		"pubspec.yaml": pubspecContent,
		"README.md":    "This is a test package.",
	}
	tarGzPath := createTestTarGz(t, tempDir, files)

	// Test reading an existing file
	data, err := ReadFile(tarGzPath, "pubspec.yaml")
	assert.NoError(t, err)
	assert.Equal(t, pubspecContent, string(data))
}

func TestReadFile_FileNotInArchive(t *testing.T) {
	tempDir := t.TempDir()
	files := map[string]string{
		"README.md": "This is a test package.",
	}
	tarGzPath := createTestTarGz(t, tempDir, files)

	_, err := ReadFile(tarGzPath, "pubspec.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file pubspec.yaml not found in archive")
}

func TestReadFile_ArchiveNotExist(t *testing.T) {
	tempDir := t.TempDir()
	tarGzPath := filepath.Join(tempDir, "not_exist.tar.gz")

	_, err := ReadFile(tarGzPath, "pubspec.yaml")
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}
