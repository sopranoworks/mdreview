package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHTTPServer_Auth(t *testing.T) {
	store := NewStore()
	token := "test-token"
	
	port, _, err := StartHTTPServer(0, store, token, 1*time.Minute)
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	client := &http.Client{}
	url := fmt.Sprintf("http://127.0.0.1:%d/ipc/ping", port)

	// Test without token
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}

	// Test with wrong token
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("X-IPC-Token", "wrong-token")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}

	// Test with correct token
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("X-IPC-Token", token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHTTPServer_Preview(t *testing.T) {
	store := NewStore()
	token := "test-token"
	
	port, _, err := StartHTTPServer(0, store, token, 1*time.Minute)
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	client := &http.Client{}
	url := fmt.Sprintf("http://127.0.0.1:%d/ipc/preview", port)

	// Create a temporary markdown file
	tmpFile, err := os.CreateTemp("", "test*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("# Hello")
	tmpFile.Close()

	// Test valid request
	body := fmt.Sprintf(`{"path": "%s"}`, tmpFile.Name())
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("X-IPC-Token", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result struct {
		ID string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if result.ID == "" {
		t.Error("expected non-empty ID")
	}

	// Test non-existent file
	body = `{"path": "/non/existent.md"}`
	req, _ = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("X-IPC-Token", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent file, got %d", resp.StatusCode)
	}

	// Test non-md file
	tmpFile2, _ := os.CreateTemp("", "test*.txt")
	defer os.Remove(tmpFile2.Name())
	tmpFile2.Close()
	body = fmt.Sprintf(`{"path": "%s"}`, tmpFile2.Name())
	req, _ = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("X-IPC-Token", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for non-md file, got %d", resp.StatusCode)
	}
}

func TestHTTPServer_IdleTimeout(t *testing.T) {
	store := NewStore()
	token := "test-token"
	timeout := 100 * time.Millisecond
	
	_, idleChan, err := StartHTTPServer(0, store, token, timeout)
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	select {
	case <-idleChan:
		// Success
	case <-time.After(500 * time.Millisecond):
		t.Error("server did not time out as expected")
	}
}
