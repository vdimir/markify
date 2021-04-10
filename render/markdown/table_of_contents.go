package markdown

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	gutil "github.com/yuin/goldmark/util"
)

const tocKeyword = "toc"

var tocResultKey = parser.NewContextKey()

type tocTree struct {
	HeadingID []byte
	Title     string
	Children  []*tocTree
	parent    *tocTree
	level     int
}

func newTocTree(level int) *tocTree {
	return &tocTree{
		HeadingID: nil,
		Title:     "",
		Children:  []*tocTree{},
		level:     level,
	}
}

func (t *tocTree) withTitleAndID(title string, headingID []byte) *tocTree {
	t.Title = title
	t.HeadingID = headingID
	return t
}

func (t *tocTree) addChild(child *tocTree) *tocTree {
	underlyingLevel := t.level + 1
	if child.level < underlyingLevel {
		if t.parent == nil {
			// actually it is error, but we don't fail, produce some result
			return t
		}
		return t.parent.addChild(child)
	}
	if child.level == underlyingLevel {
		child.parent = t
		t.Children = append(t.Children, child)
		return child
	}
	dummyChild := newTocTree(underlyingLevel)
	dummyChild.parent = t
	t.Children = append(t.Children, dummyChild)
	return dummyChild.addChild(child)
}

func (t *tocTree) writeAHref(w io.Writer, lo int, hi int) {
	if t.Title == "" || !t.inRange(lo, hi) {
		return
	}
	if t.HeadingID != nil {
		w.Write([]byte(fmt.Sprintf("<a href=\"#%s\">", t.HeadingID)))
	} else {
		w.Write([]byte("<a>"))
	}
	w.Write([]byte(t.Title))
	w.Write([]byte("</a>"))
}

func (t *tocTree) writeBlockBeg(w io.Writer, lo int, hi int) {
	if t.level == 0 {
		w.Write([]byte("<nav class=\"toc-block\">"))
	} else if t.inRange(lo, hi) {
		w.Write([]byte("<li>"))
	}
}

func (t *tocTree) writeBlockEnd(w io.Writer, lo int, hi int) {
	if t.level == 0 {
		w.Write([]byte("</nav>"))
	} else if t.inRange(lo, hi) {
		w.Write([]byte("</li>"))
	}
}

func (t *tocTree) inRange(lo int, hi int) bool {
	return lo <= t.level && t.level <= hi
}

func (t *tocTree) toHTML(w io.Writer, lo int, hi int) {
	t.writeBlockBeg(w, lo, hi)
	t.writeAHref(w, lo, hi)
	if len(t.Children) > 0 {
		if t.inRange(lo, hi) {
			w.Write([]byte("<ul>"))
		}
		for _, ch := range t.Children {
			ch.toHTML(w, lo, hi)
		}
		if t.inRange(lo, hi) {
			w.Write([]byte("</ul>"))
		}
	}
	t.writeBlockEnd(w, lo, hi)
}

// --- transformer ---

// tocTransformer extract toc from document
type tocTransformer struct{}

func (t *tocTransformer) Transform(n *gast.Document, reader text.Reader, pc parser.Context) {
	if enabled, ok := pc.Get(EnableShortcodes).(bool); ok && !enabled {
		// do not collect toc if shortcodes not enabled
		return
	}

	toctree := newTocTree(0)
	currentRoot := toctree
	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n.Kind() != ast.KindHeading || !entering {
			return ast.WalkContinue, nil
		}
		headingNode := n.(*ast.Heading)

		var headerID []byte
		if id, idFound := headingNode.AttributeString("id"); idFound {
			headerID = id.([]byte)
		}

		headingText := &bytes.Buffer{}
		err := extractTextFromNode(n, reader, headingText)
		if err != nil {
			return ast.WalkContinue, err
		}
		newNode := newTocTree(headingNode.Level).withTitleAndID(headingText.String(), headerID)

		currentRoot = currentRoot.addChild(newNode)
		return ast.WalkSkipChildren, nil
	})

	if ns, ok := pc.Get(tocResultKey).([]gast.Node); ok && n != nil {
		for _, n := range ns {
			n.(*ShortCodeNode).context = toctree
		}
	}
}

// --- render ---

// tocHTMLRenderer renders toc
type tocHTMLRenderer struct{}

// RegisterFuncs for tocHTMLRenderer
func (r *tocHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	if kind, ok := ShortCodeNodeKinds[tocKeyword]; ok {
		reg.Register(kind, r.renderTableOfContets)
	} else {
		panic(fmt.Errorf("missing NodeKind for %s", tocKeyword))
	}
}

func (r *tocHTMLRenderer) renderTableOfContets(w gutil.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if n, ok := n.(*ShortCodeNode); entering && ok {
		tocCtx := n.context.(*tocTree)
		if tocCtx == nil {
			return gast.WalkContinue, nil
		}
		lo, hi := 1, 6
		if len(n.params) >= 2 {
			userHi, err1 := strconv.ParseInt(string(n.params[0]), 10, 32)
			userLo, err2 := strconv.ParseInt(string(n.params[1]), 10, 32)
			if err1 == nil && err2 == nil {
				lo, hi = int(userHi), int(userLo)
			}
		}
		tocCtx.toHTML(w, lo, hi)
	}

	return gast.WalkContinue, nil
}

// --- extender ---

// tocExtender is an extension for toc shortcode
type tocExtender struct{}

// TableOfContentsShortcode allows to insert Table Of Contents with {{ toc }} shortcode
var TableOfContentsShortcode = &tocExtender{}

// Extend with tocExtender
func (e *tocExtender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		gutil.Prioritized(&shortCodeParser{
			keyword:          tocKeyword,
			parserContextKey: &tocResultKey,
		}, 200),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		gutil.Prioritized(&tocHTMLRenderer{}, 200),
	))
	m.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(&tocTransformer{}, 10),
	))
}
