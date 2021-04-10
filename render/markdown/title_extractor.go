package markdown

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// titleParserCtxKey stores title extracted from first header
var titleParserCtxKey = parser.NewContextKey()

const maxTitleLen = 120

type titleExtractorTransformer struct{}

// PagePreviewText contains title of document and beginning of content
type PagePreviewText struct {
	Title   string
	Preview string
}

func (t *titleExtractorTransformer) Transform(n *gast.Document, reader text.Reader, pc parser.Context) {
	if _, ok := pc.Get(titleParserCtxKey).(*PagePreviewText); ok {
		return
	}

	pc.Set(titleParserCtxKey, &PagePreviewText{})

	_ = ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		buf := &bytes.Buffer{}
		dst := pc.Get(titleParserCtxKey).(*PagePreviewText)

		if n.Kind() == ast.KindHeading && dst.Title == "" {
			extractTextFromNode(n, reader, buf)
			dst.Title = buf.String()
		}

		if n.Kind() == ast.KindParagraph && dst.Preview == "" {
			extractTextFromNode(n, reader, buf)
			dst.Preview = buf.String()
		}

		if dst.Title != "" && dst.Preview != "" {
			return ast.WalkStop, stopWalkError
		}

		return ast.WalkContinue, nil
	})
}
