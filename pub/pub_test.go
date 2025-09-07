package pub

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSpec(t *testing.T) {
	tests := []struct {
		name       string
		content    []byte
		want       *Spec
		wantErr    bool
		wantErrStr string
	}{
		{
			name:    "valid spec",
			content: []byte("name: test\nversion: 1.0.0"),
			want: &Spec{
				Name:    "test",
				Version: "1.0.0",
				Raw:     map[string]any{"name": "test", "version": "1.0.0"},
			},
		},
		{
			name:       "invalid spec - no name",
			content:    []byte("version: 1.0.0"),
			wantErr:    true,
			wantErrStr: "invalid pubspec data: name",
		},
		{
			name:       "invalid spec - no version",
			content:    []byte("name: test"),
			wantErr:    true,
			wantErrStr: "invalid pubspec data: version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSpec(tt.content)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrStr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestReadSpecFromTar(t *testing.T) {
	tempDir := t.TempDir()
	tarPath := filepath.Join(tempDir, "test.tar.gz")

	// Create a dummy tar.gz file
	_, err := os.Create(tarPath)
	require.NoError(t, err)

	// TODO: Add a real tar.gz file with a pubspec.yaml file
}
