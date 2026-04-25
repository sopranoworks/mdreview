package mcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestProxyPreviewRequest(t *testing.T) {
	token := "test-token"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ipc/preview" {
			t.Errorf("expected path /ipc/preview, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if r.Header.Get("X-IPC-Token") != token {
			t.Errorf("expected token %s, got %s", token, r.Header.Get("X-IPC-Token"))
		}
		var req struct {
			Path string `json:"path"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.Path != "test.md" {
			t.Errorf("expected path test.md, got %s", req.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": "test-id"})
	}))
	defer ts.Close()

	addr := ts.Listener.Addr().String()
	parts := strings.Split(addr, ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	url, err := ProxyPreviewRequest(port, token, "test.md")
	if err != nil {
		t.Fatalf("ProxyPreviewRequest failed: %v", err)
	}

	expectedSuffix := "/rev/test-id"
	if !strings.HasSuffix(url, expectedSuffix) {
		t.Errorf("expected URL to end with %s, got %s", expectedSuffix, url)
	}
}
