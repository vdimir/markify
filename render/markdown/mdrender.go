package markdown

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"

	highlighting "github.com/yuin/goldmark-highlighting"

	"github.com/alecthomas/chroma/formatters/html"
)

// Render renders markdown to html
type Converter struct {
	markdown goldmark.Markdown
}

type Document struct {
	Title   string
	Preview string
	Body    string
}

// NewRender create new renderer
func NewConverter() *Converter {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(true),
				),
			),
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

	return &Converter{md}
}

// Converter markdown to html
func (r *Converter) Convert(data []byte) (*Document, error) {
	var ctx = parser.NewContext()

	var htmlBuf bytes.Buffer
	if err := r.markdown.Convert(data, &htmlBuf, parser.WithContext(ctx)); err != nil {
		return nil, err
	}

	doc := &Document{
		Body: htmlBuf.String(),
	}
	previewText, ok := ctx.Get(titleParserCtxKey).(*PagePreviewText)
	if ok && previewText != nil {
		doc.Title = previewText.Title
		doc.Preview = previewText.Preview
	}
	return doc, nil
}
