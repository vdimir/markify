package md

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	gutil "github.com/yuin/goldmark/util"
)

// templateEmbedHTMLRenderer render embebed representation formatting specified template
type templateEmbedHTMLRenderer struct {
	name   string
	center bool
	tpl    *template.Template
}

// RegisterFuncs for templateEmbedHTMLRenderer
func (r *templateEmbedHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ShortCodeNodeKinds[r.name], r.renderTemplateEmbed)
}

func (r *templateEmbedHTMLRenderer) renderTemplateEmbed(w gutil.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if !entering {
		if r.center {
			_, _ = w.WriteString("</center>")
		}
		return gast.WalkContinue, nil
	}

	if r.center {
		_, _ = w.WriteString("<center>")
	}
	params := []string{}
	for _, p := range n.(*ShortCodeNode).params {
		params = append(params, string(p))
	}
	buf := &bytes.Buffer{}
	if err := r.tpl.Execute(buf, params); err != nil {
		log.Printf("[ERROR] %s", err)
		_, _ = w.WriteString(fmt.Sprintf("Unable to display %s", r.name))
	} else {
		_, _ = w.Write(buf.Bytes())
	}

	return gast.WalkContinue, nil
}

// --- extender ---

// tplEmbedExtender is an extension that allows to use template embeding
type tplEmbedExtender struct {
	name string
	tpl  *template.Template
}

// Extend with tplEmbedExtender
func (e *tplEmbedExtender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		gutil.Prioritized(&shortCodeParser{
			keyword: e.name,
		}, 500),
	))
	rend := &templateEmbedHTMLRenderer{
		name:   e.name,
		center: false,
		tpl:    e.tpl,
	}
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		gutil.Prioritized(rend, 500),
	))
}
