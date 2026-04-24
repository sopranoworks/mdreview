package mcp

import (
	"os/exec"
	"strings"
)

// GetTailscaleIP returns the Tailscale IPv4 address of the current machine.
// If Tailscale is not running or the command fails, it returns "127.0.0.1".
func GetTailscaleIP() string {
	cmd := exec.Command("tailscale", "ip", "-4")
	out, err := cmd.Output()
	if err != nil {
		return "127.0.0.1"
	}
	ip := strings.TrimSpace(string(out))
	if ip == "" {
		return "127.0.0.1"
	}
	return ip
}
