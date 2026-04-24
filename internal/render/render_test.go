package render

import (
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	input := "# Hello"
	expected := "<h1>Hello</h1>\n"
	output, err := RenderMarkdown(input)
	if err != nil {
		t.Fatalf("RenderMarkdown failed: %v", err)
	}
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}
