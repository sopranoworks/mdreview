package render

import (
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic markdown",
			input:    "# Hello",
			expected: "<h1>Hello</h1>\n",
		},
		{
			name:     "XSS sanitization - script tag",
			input:    "# Hello <script>alert('xss')</script>",
			expected: "<h1>Hello </h1>\n",
		},
		{
			name:     "XSS sanitization - onload attribute",
			input:    `Hello <img src="x" onerror="alert('xss')">`,
			expected: "<p>Hello <img src=\"x\"></p>\n",
		},
		{
			name:     "XSS sanitization - javascript link",
			input:    `[click me](javascript:alert('xss'))`,
			expected: "<p>click me</p>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := RenderMarkdown(tt.input)
			if err != nil {
				t.Fatalf("RenderMarkdown failed: %v", err)
			}
			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}
