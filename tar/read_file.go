package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func ReadFile(tarGzPath, targetFile string) ([]byte, error) {
	// Open the tar.gz file
	f, err := os.Open(tarGzPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tar.gz: %w", err)
	}
	defer f.Close()

	// Wrap gzip reader
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gz.Close()

	// Create a tar reader
	tr := tar.NewReader(gz)

	// Iterate through tar entries
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // end of archive
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar entry: %w", err)
		}

		if header.Typeflag == tar.TypeReg && header.Name == targetFile {
			// Found the file, read contents
			data, err := io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("failed to read file contents: %w", err)
			}
			return data, nil
		}
	}

	return nil, fmt.Errorf("file %s not found in archive", targetFile)
}
