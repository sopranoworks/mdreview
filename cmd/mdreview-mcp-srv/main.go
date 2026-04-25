package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sopranoworks/mdreview/internal/server"
	"github.com/sopranoworks/mdreview/internal/version"
)

func main() {
	fmt.Printf("mdreview-mcp-srv version %s\n", version.Version)

	store := server.NewStore()
	store.StartCleanup(1*time.Minute, 10*time.Minute)

	token, err := server.GenerateToken()
	if err != nil {
		log.Fatalf("failed to generate token: %v", err)
	}

	port, idleChan, err := server.StartHTTPServer(0, store, token, 10*time.Minute)
	if err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}

	if err := server.WriteState(port, token); err != nil {
		log.Fatalf("failed to write state: %v", err)
	}
	defer server.RemoveState()

	fmt.Printf("mdreview-mcp-srv listening on port %d\n", port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		fmt.Println("Shutting down (signal)...")
	case <-idleChan:
		fmt.Println("Shutting down (idle timeout)...")
	}
}
