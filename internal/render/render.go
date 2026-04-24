package render

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

// RenderMarkdown renders markdown string to HTML string.
func RenderMarkdown(source string) (string, error) {
	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	p := bluemonday.UGCPolicy()
	return p.Sanitize(buf.String()), nil
}
