package server

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sopranoworks/mdreview/internal/render"
)

// StartHTTPServer starts the HTTP server on the specified port.
func StartHTTPServer(port int, store *Store, token string, idleTimeout time.Duration) (int, <-chan struct{}, error) {
	idleChan := make(chan struct{})
	timer := time.NewTimer(idleTimeout)

	resetTimer := func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(idleTimeout)
	}

	go func() {
		<-timer.C
		close(idleChan)
	}()

	authMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			resetTimer()

			// Check if remote address is local
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil || (host != "127.0.0.1" && host != "::1" && host != "localhost") {
				http.Error(w, "Forbidden: local access only", http.StatusForbidden)
				return
			}

			if r.Header.Get("X-IPC-Token") != token {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next(w, r)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/rev/", func(w http.ResponseWriter, r *http.Request) {
		resetTimer()
		id := r.URL.Path[len("/rev/"):]
		if id == "" {
			http.NotFound(w, r)
			return
		}
		content, ok := store.Get(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, content)
	})

	mux.HandleFunc("/ipc/preview", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Path string `json:"path"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Path validation
		if !strings.HasSuffix(req.Path, ".md") {
			http.Error(w, "Only .md files are allowed", http.StatusBadRequest)
			return
		}
		if _, err := os.Stat(req.Path); os.IsNotExist(err) {
			http.Error(w, "File does not exist", http.StatusNotFound)
			return
		}

		content, err := os.ReadFile(req.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rendered, err := render.RenderMarkdown(string(content))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := fmt.Sprintf("%x", sha256.Sum256([]byte(token+req.Path)))
		store.Set(id, rendered)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}))

	mux.HandleFunc("/ipc/ping", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong")
	}))

	// Bind to all interfaces to allow Tailscale access
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, nil, err
	}

	actualPort := ln.Addr().(*net.TCPAddr).Port

	go func() {
		if err := http.Serve(ln, mux); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	return actualPort, idleChan, nil
}
