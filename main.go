package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"mdreview/internal/mcp"
	"mdreview/internal/server"
)

func main() {
	port := flag.Int("port", 8080, "Port for the HTTP preview server")
	workspace := flag.String("workspace", ".", "Path to the workspace root")
	flag.Parse()

	// Allow environment variables to override flags
	if envPort := os.Getenv("MDREVIEW_PORT"); envPort != "" {
		p, err := strconv.Atoi(envPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid MDREVIEW_PORT: %v. Using default %d\n", err, *port)
		} else {
			*port = p
		}
	}
	if envWorkspace := os.Getenv("MDREVIEW_WORKSPACE"); envWorkspace != "" {
		*workspace = envWorkspace
	}

	store := server.NewStore()

	// Start HTTP Preview Server
	actualPort, err := server.StartHTTPServer(*port, store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start HTTP server on port %d: %v. Retrying with automatic port selection...\n", *port, err)
		actualPort, err = server.StartHTTPServer(0, store)
		if err != nil {
			log.Fatalf("Critical: Failed to start HTTP server even with automatic port: %v", err)
		}
	}

	fmt.Fprintf(os.Stderr, "HTTP preview server active on port %d\n", actualPort)

	// Start MCP Server (stdio)
	fmt.Fprintf(os.Stderr, "Starting MCP server (workspace: %s)...\n", *workspace)
	mcpServer := mcp.NewMCPServer(*workspace, actualPort, store)
	mcpServer.Run()
}
