package server

import (
	"fmt"
	"net/http"
)

func StartHTTPServer(port int, store *Store) error {
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

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, mux)
}
