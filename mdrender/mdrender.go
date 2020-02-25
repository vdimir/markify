package md

import (
	"bytes"
	"net/url"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
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
		goldmark.WithRendererOptions(),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	return &Render{md: md}, nil
}

// Render markdown to html
func (r *Render) Render(data []byte, opt *Options) (*bytes.Buffer, error) {
	var ctx = parser.NewContext()
	if opt != nil {
		opt.fillParserCtx(ctx)
	}

	var htmlBuf bytes.Buffer
	if err := r.md.Convert(data, &htmlBuf, parser.WithContext(ctx)); err != nil {
		return nil, err
	}
	// usedShortcodes := ctx.Get(UsedShortcodesKey).(map[string]bool)

	return &htmlBuf, nil
}
