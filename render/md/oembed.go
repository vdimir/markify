package md

import (
	"fmt"
	"log"
	"regexp"

	mutil "github.com/vdimir/markify/util"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"

	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	gutil "github.com/yuin/goldmark/util"
)

type oembedResp struct {
	Data string `json:"html"`
}

var resourceIDPattern = regexp.MustCompile("^[A-Za-z0-9_\\-]+$")

// oembedHTMLRenderer render embebed representation of URL fom third-party sites
type oembedHTMLRenderer struct {
	name        string
	center      bool
	urlTemplate string
}

// RegisterFuncs for oembedHTMLRenderer
func (r *oembedHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	if kind, ok := ShortCodeNodeKinds[r.name]; ok {
		reg.Register(kind, r.renderOembed)
	} else {
		panic(fmt.Errorf("missing NodeKind for %s", r.name))
	}
}

func getOembedData(urlPath string) (string, error) {
	resp := &oembedResp{}
	err := mutil.GetJSON(urlPath, resp)
	if err != nil {
		return "", errors.Wrap(err, "cannot get OEmbed code via API")
	}
	return resp.Data, nil
}

func (r *oembedHTMLRenderer) renderOembed(w gutil.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if !entering {
		if r.center {
			_, _ = w.WriteString("</center>")
		}
		return gast.WalkContinue, nil
	}

	if r.center {
		_, _ = w.WriteString("<center>")
	}

	params := n.(*ShortCodeNode).params
	if len(params) == 0 || !resourceIDPattern.Match(params[0]) {
		_, _ = w.WriteString(fmt.Sprintf("Wrong arguments for %s!", r.name))
		return gast.WalkContinue, nil
	}

	resourceID := params[0]
	urlPath := fmt.Sprintf(r.urlTemplate, resourceID)
	data, err := getOembedData(urlPath)

	if err != nil {
		log.Printf("[ERROR] %s", err)
		data = fmt.Sprintf("Unable to load %s %s!", r.name, resourceID)
	}
	_, _ = w.WriteString(data)

	return gast.WalkContinue, nil
}

// --- extender ---

// oembedExtender is an extension that allows to use embeding
type oembedExtender struct {
	name        string
	urlTemplate string
}

// Extend with oembedExtender
func (e *oembedExtender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		gutil.Prioritized(&shortCodeParser{
			keyword: e.name,
		}, 500),
	))
	rend := &oembedHTMLRenderer{
		name:        e.name,
		center:      true,
		urlTemplate: e.urlTemplate,
	}
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		gutil.Prioritized(rend, 500),
	))
}
