package server

import (
	"fmt"
	"net"
	"net/http"
)

// StartHTTPServer starts the HTTP server on the specified port.
// If port is 0, it will pick an available port.
// It returns the actual port used and an error if any.
func StartHTTPServer(port int, store *Store) (int, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rev/", func(w http.ResponseWriter, r *http.Request) {
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

	// Use net.Listen to allow for port 0 (auto-selection)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, err
	}

	actualPort := ln.Addr().(*net.TCPAddr).Port

	go func() {
		if err := http.Serve(ln, mux); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	return actualPort, nil
}
