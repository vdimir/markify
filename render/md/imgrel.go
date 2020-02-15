package md

import (
	"net/url"
	"path"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var (
	// LinkBaseURL context key for markdown source url to replace relative links
	LinkBaseURL = parser.NewContextKey()
)

// ImgReplaceRelativeLink replaces relative image links to absolute
type ImgReplaceRelativeLink struct{}

// Extend goldmark converter
func (t ImgReplaceRelativeLink) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&imgRelativeLink{}, 100),
		),
	)
}

type imgRelativeLink struct{}

func (t *imgRelativeLink) Transform(n *ast.Document, reader text.Reader, pc parser.Context) {
	srcURL, ok := pc.Get(LinkBaseURL).(*url.URL)
	if !ok || srcURL == nil {
		return
	}
	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		if entering && (n.Kind() == ast.KindImage) {
			origDest := string(n.(*ast.Image).Destination)
			origURL, err := url.Parse(origDest)
			if err != nil {
				return s, err
			}
			if origURL.Host == "" {
				origURL.Host = srcURL.Host
				origURL.Scheme = srcURL.Scheme
				origURL.Path = path.Join(path.Dir(srcURL.Path), origURL.Path)
			}
			n.(*ast.Image).Destination = []byte(origURL.String())
		}
		return s, nil
	})
}
