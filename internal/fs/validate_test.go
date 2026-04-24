package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mdreview-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	baseDir := filepath.Join(tmpDir, "workspace")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		t.Fatalf("failed to create workspace dir: %v", err)
	}

	// Create a file outside the workspace
	secretFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(secretFile, []byte("secret"), 0644); err != nil {
		t.Fatalf("failed to create secret file: %v", err)
	}

	// Create a symlink inside the workspace pointing outside
	symlinkPath := filepath.Join(baseDir, "escape-link")
	if err := os.Symlink(secretFile, symlinkPath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

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
		{
			name:       "symlink escape attempt",
			targetPath: "escape-link",
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
