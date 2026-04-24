package mcp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GetTailscaleIP returns the Tailscale IPv4 address of the current machine.
// If Tailscale is not running or the command fails, it returns "127.0.0.1".
func GetTailscaleIP() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "tailscale", "ip", "-4")
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to get Tailscale IP: %v\n", err)
		return "127.0.0.1"
	}
	ip := strings.TrimSpace(string(out))
	if ip == "" {
		return "127.0.0.1"
	}
	return ip
}
