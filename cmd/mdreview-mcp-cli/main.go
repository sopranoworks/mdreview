package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sopranoworks/mdreview/internal/mcp"
	"github.com/sopranoworks/mdreview/internal/version"
)

func main() {
	portFlag := flag.Int("port", 0, "Preferred port for the HTTP sidecar server")
	workspaceFlag := flag.String("workspace", "", "Path to the workspace root")
	flag.Parse()

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("mdreview-mcp-cli version %s\n", version.Version)
		return
	}

	workspacePath := *workspaceFlag
	if workspacePath == "" {
		var err error
		workspacePath, err = os.Getwd()
		if err != nil {
			log.Fatalf("failed to get current directory: %v", err)
		}
	}

	port, token, err := mcp.DiscoverOrStartServer(*portFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to discover or start sidecar server: %v\n", err)
		os.Exit(1)
	}

	server := mcp.NewMCPServer(workspacePath).WithSidecar(port, token)
	server.Run()
}
