package md

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// Render renders markdown to html
type Render struct {
	md goldmark.Markdown
}

// NewRender create new renderer
func NewRender() (*Render, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
			TableOfContentsShortcode,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(&titleExtractorTransformer{}, 500),
			),
		),
		goldmark.WithRendererOptions(),
	)

	return &Render{md: md}, nil
}

// Render markdown to html
func (r *Render) Render(data []byte) (*bytes.Buffer, parser.Context, error) {
	var ctx = parser.NewContext()

	var htmlBuf bytes.Buffer
	if err := r.md.Convert(data, &htmlBuf, parser.WithContext(ctx)); err != nil {
		return nil, nil, err
	}

	return &htmlBuf, ctx, nil
}
