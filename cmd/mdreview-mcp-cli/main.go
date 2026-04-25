package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sopranoworks/mdreview/internal/mcp"
	"github.com/sopranoworks/mdreview/internal/version"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("mdreview-mcp-cli version %s\n", version.Version)
		return
	}

	workspacePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	port, token, err := mcp.DiscoverOrStartServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to discover or start sidecar server: %v\n", err)
		os.Exit(1)
	}

	server := mcp.NewMCPServer(workspacePath).WithSidecar(port, token)
	server.Run()
}
