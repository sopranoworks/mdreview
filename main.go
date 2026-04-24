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

	// Start HTTP Preview Server in background
	go func() {
		fmt.Fprintf(os.Stderr, "Starting HTTP preview server on port %d...\n", *port)
		if err := server.StartHTTPServer(*port, store); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Start MCP Server (stdio)
	fmt.Fprintf(os.Stderr, "Starting MCP server (workspace: %s)...\n", *workspace)
	mcpServer := mcp.NewMCPServer(*workspace, *port, store)
	mcpServer.Run()
}
