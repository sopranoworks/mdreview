package fs

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidatePath ensures that targetPath is within baseDir.
// It returns the absolute path or an error if it escapes the workspace.
func ValidatePath(baseDir, targetPath string) (string, error) {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute base path: %w", err)
	}

	var absTarget string
	if filepath.IsAbs(targetPath) {
		absTarget = filepath.Clean(targetPath)
	} else {
		absTarget = filepath.Join(absBase, targetPath)
	}

	rel, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}

	if strings.HasPrefix(rel, "..") || strings.HasPrefix(rel, "/") {
		return "", fmt.Errorf("path escapes workspace: %s", targetPath)
	}

	return absTarget, nil
}
