package md

import (
	"bytes"
	"net/url"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// Options config for markdown renderer
type Options struct {
	BaseURL           *url.URL
	DisableShortcodes bool
}

func (opt *Options) fillParserCtx(ctx parser.Context) {
	if opt.BaseURL != nil {
		ctx.Set(LinkBaseURL, opt.BaseURL)
	}
	ctx.Set(EnableShortcodes, !opt.DisableShortcodes)
}

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
			ImgReplaceRelativeLink{},
			EmbedTweet,
			EmbedInstagram,
			EmbedGist,
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
func (r *Render) Render(data []byte, opt *Options) (*bytes.Buffer, parser.Context, error) {
	var ctx = parser.NewContext()
	if opt != nil {
		opt.fillParserCtx(ctx)
	}

	var htmlBuf bytes.Buffer
	if err := r.md.Convert(data, &htmlBuf, parser.WithContext(ctx)); err != nil {
		return nil, nil, err
	}

	return &htmlBuf, ctx, nil
}
