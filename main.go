package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sopranoworks/mdreview/internal/mcp"
	"github.com/sopranoworks/mdreview/internal/server"
)

func main() {
	portFlag := flag.String("port", "", "Port for the HTTP preview server")
	workspaceFlag := flag.String("workspace", "", "Path to the workspace root")
	flag.Parse()

	port := 8080
	if *portFlag != "" {
		p, err := strconv.Atoi(*portFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid port %q: %v. Using default 8080\n", *portFlag, err)
		} else {
			port = p
		}
	}

	workspace := "."
	if *workspaceFlag != "" {
		workspace = *workspaceFlag
	}

	// Allow environment variables to override flags
	if envPort := os.Getenv("MDREVIEW_PORT"); envPort != "" {
		p, err := strconv.Atoi(envPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid MDREVIEW_PORT: %v. Using current %d\n", err, port)
		} else {
			port = p
		}
	}
	if envWorkspace := os.Getenv("MDREVIEW_WORKSPACE"); envWorkspace != "" {
		workspace = envWorkspace
	}

	store := server.NewStore()

	// Start HTTP Preview Server
	actualPort, err := server.StartHTTPServer(port, store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start HTTP server on port %d: %v. Retrying with automatic port selection...\n", port, err)
		actualPort, err = server.StartHTTPServer(0, store)
		if err != nil {
			log.Fatalf("Critical: Failed to start HTTP server even with automatic port: %v", err)
		}
	}

	fmt.Fprintf(os.Stderr, "HTTP preview server active on port %d\n", actualPort)

	// Start MCP Server (stdio)
	fmt.Fprintf(os.Stderr, "Starting MCP server (workspace: %s)...\n", workspace)
	mcpServer := mcp.NewMCPServer(workspace, actualPort, store)
	mcpServer.Run()
}
