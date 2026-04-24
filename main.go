package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"mdreview/internal/mcp"
	"mdreview/internal/server"
)

func main() {
	port := flag.Int("port", 8080, "Port for the HTTP preview server")
	workspace := flag.String("workspace", ".", "Path to the workspace root")
	flag.Parse()

	// Allow environment variables to override flags
	if envPort := os.Getenv("MDREVIEW_PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", port)
	}
	if envWorkspace := os.Getenv("MDREVIEW_WORKSPACE"); envWorkspace != "" {
		*workspace = envWorkspace
	}

	store := server.NewStore()

	// Start HTTP Preview Server in background
	go func() {
		if err := server.StartHTTPServer(*port, store); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Start MCP Server (stdio)
	mcpServer := mcp.NewMCPServer(*workspace, *port, store)
	mcpServer.Run()
}
