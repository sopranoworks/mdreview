package render

import (
	"bytes"

	"github.com/yuin/goldmark"
)

// RenderMarkdown renders markdown string to HTML string.
func RenderMarkdown(source string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
