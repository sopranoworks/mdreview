package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/sopranoworks/mdreview/internal/server"
)

// DiscoverOrStartServer attempts to find a running sidecar server or starts a new one.
func DiscoverOrStartServer(preferredPort int) (int, string, error) {
	pid, port, token, err := server.ReadState()
	if err == nil {
		// Check if process is running
		process, err := os.FindProcess(pid)
		if err == nil {
			// On Unix, FindProcess always succeeds. We need to send signal 0 to check if it's alive.
			err = process.Signal(syscall.Signal(0))
			if err == nil {
				// Process is running, try to ping
				if pingServer(port, token) {
					return port, token, nil
				}
			}
		}
	}

	// Server not found or unresponsive, start it
	return startServer(preferredPort)
}

func pingServer(port int, token string) bool {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/ipc/ping", port), nil)
	if err != nil {
		return false
	}
	req.Header.Set("X-IPC-Token", token)
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func startServer(preferredPort int) (int, string, error) {
	srvPath := "mdreview-mcp-srv"

	// Try to find it in the same directory as the current executable
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		localSrv := filepath.Join(dir, "mdreview-mcp-srv")
		if _, err := os.Stat(localSrv); err == nil {
			srvPath = localSrv
		}
	}

	args := []string{}
	if preferredPort > 0 {
		args = append(args, "-port", fmt.Sprintf("%d", preferredPort))
	}
	cmd := exec.Command(srvPath, args...)
	// Ensure it runs in the background and doesn't get killed when the CLI exits

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return 0, "", fmt.Errorf("failed to start %s: %w", srvPath, err)
	}

	// Wait for state file to appear and server to be responsive
	for i := 0; i < 20; i++ {
		time.Sleep(500 * time.Millisecond)
		_, port, token, err := server.ReadState()
		if err == nil {
			if pingServer(port, token) {
				return port, token, nil
			}
		}
	}

	return 0, "", fmt.Errorf("timed out waiting for sidecar server to start")
}

// ProxyPreviewRequest sends a preview request to the sidecar server and returns the preview URL.
func ProxyPreviewRequest(port int, token string, path string) (string, error) {
	reqBody, err := json.Marshal(map[string]string{"path": path})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/ipc/preview", port), bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-IPC-Token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%d/rev/%s", GetTailscaleIP(), port, result.ID), nil
}
