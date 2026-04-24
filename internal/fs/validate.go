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

	// Resolve symlinks in base directory
	absBase, err = filepath.EvalSymlinks(absBase)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlinks in base path: %w", err)
	}

	var absTarget string
	if filepath.IsAbs(targetPath) {
		absTarget = filepath.Clean(targetPath)
	} else {
		absTarget = filepath.Join(absBase, targetPath)
	}

	// Resolve symlinks in target path
	// Note: EvalSymlinks requires the path to exist.
	// If it doesn't exist, we still want to validate the path.
	// However, for security, we should resolve what we can.
	resolvedTarget, err := filepath.EvalSymlinks(absTarget)
	if err == nil {
		absTarget = resolvedTarget
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
