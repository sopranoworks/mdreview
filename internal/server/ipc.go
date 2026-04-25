package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const stateFileName = "mdreview.state"

func getStateFilePath() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("mdreview.%d.state", os.Getuid()))
}

// GenerateToken generates a random 16-byte hex token.
func GenerateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// WriteState writes the current process ID, port, and token to a state file.
func WriteState(port int, token string) error {
	data := fmt.Sprintf("%d:%d:%s", os.Getpid(), port, token)
	return os.WriteFile(getStateFilePath(), []byte(data), 0600)
}

// ReadState reads the process ID, port, and token from the state file.
func ReadState() (int, int, string, error) {
	data, err := os.ReadFile(getStateFilePath())
	if err != nil {
		return 0, 0, "", err
	}
	parts := strings.Split(string(data), ":")
	if len(parts) != 3 {
		return 0, 0, "", fmt.Errorf("invalid state file format")
	}
	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, "", err
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, "", err
	}
	return pid, port, parts[2], nil
}

// RemoveState removes the state file.
func RemoveState() error {
	return os.Remove(getStateFilePath())
}
