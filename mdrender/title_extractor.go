package md

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// MdTitleKey stores title extracted from first header
var MdTitleKey = parser.NewContextKey()

const maxTitleLen = 120

type titleExtractorTransformer struct{}

// PagePreviewText contains title of document and beginning of content
type PagePreviewText struct {
	Title string
	Body  string
}

func (t *titleExtractorTransformer) Transform(n *gast.Document, reader text.Reader, pc parser.Context) {
	if _, ok := pc.Get(MdTitleKey).(*PagePreviewText); ok {
		return
	}

	pc.Set(MdTitleKey, &PagePreviewText{})

	_ = ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		buf := &bytes.Buffer{}
		dst := pc.Get(MdTitleKey).(*PagePreviewText)

		if n.Kind() == ast.KindHeading && dst.Title == "" {
			extractTextFromNode(n, reader, buf)
			dst.Title = buf.String()
		}

		if n.Kind() == ast.KindParagraph && dst.Body == "" {
			extractTextFromNode(n, reader, buf)
			dst.Body = buf.String()
		}

		if dst.Title != "" && dst.Body != "" {
			return ast.WalkStop, stopWalkError{}

		}

		return ast.WalkContinue, nil
	})
}
