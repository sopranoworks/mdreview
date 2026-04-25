package server

import (
	"os"
	"testing"
)

func TestState(t *testing.T) {
	port := 12345
	token := "test-token"
	err := WriteState(port, token)
	if err != nil {
		t.Fatalf("WriteState failed: %v", err)
	}
	defer RemoveState()

	pid, rPort, rToken, err := ReadState()
	if err != nil {
		t.Fatalf("ReadState failed: %v", err)
	}

	if pid != os.Getpid() {
		t.Errorf("expected pid %d, got %d", os.Getpid(), pid)
	}
	if rPort != port {
		t.Errorf("expected port %d, got %d", port, rPort)
	}
	if rToken != token {
		t.Errorf("expected token %s, got %s", token, rToken)
	}
}
