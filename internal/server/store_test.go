package server

import (
	"fmt"
	"testing"
)

func TestStore(t *testing.T) {
	store := NewStore()
	id := "test-id"
	content := "<h1>Hello World</h1>"

	store.Set(id, content)

	got, ok := store.Get(id)
	if !ok {
		t.Errorf("expected to find content for id %s", id)
	}
	if got != content {
		t.Errorf("expected %q, got %q", content, got)
	}

	_, ok = store.Get("non-existent")
	if ok {
		t.Error("expected not to find content for non-existent id")
	}
}

func TestStoreThreadSafety(t *testing.T) {
	store := NewStore()
	const iterations = 1000

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < iterations; j++ {
				store.Set(fmt.Sprintf("id-%d-%d", id, j), "content")
				store.Get(fmt.Sprintf("id-%d-%d", id, j))
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
