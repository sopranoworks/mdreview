package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	baseDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	baseDir = filepath.Join(baseDir, "testdata")

	tests := []struct {
		name       string
		targetPath string
		wantErr    bool
	}{
		{
			name:       "valid path",
			targetPath: "file.md",
			wantErr:    false,
		},
		{
			name:       "valid path in subdirectory",
			targetPath: "subdir/file.md",
			wantErr:    false,
		},
		{
			name:       "path traversal attempt",
			targetPath: "../../etc/passwd",
			wantErr:    true,
		},
		{
			name:       "path traversal attempt with absolute path",
			targetPath: "/etc/passwd",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidatePath(baseDir, tt.targetPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
