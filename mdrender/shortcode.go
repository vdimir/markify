package md

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var (
	// EnableShortcodes context key for indicates if shortcodes enabled
	EnableShortcodes = parser.NewContextKey()
	// UsedShortcodesKey store list of used shortcodes on page
	UsedShortcodesKey = parser.NewContextKey()
)

// ShortCodeNodeKinds mapping `keyword -> NodeKind`
var ShortCodeNodeKinds = map[string]gast.NodeKind{
	"tweet":     gast.NewNodeKind("EmbedTweet"),
	"gist":      gast.NewNodeKind("EmbedGist"),
	"instagram": gast.NewNodeKind("EmbedInstagram"),
	tocKeyword:  gast.NewNodeKind("TOCShortcode"),
}

// A ShortCodeNode struct represents a tweet ast node
type ShortCodeNode struct {
	gast.BaseInline
	kind    ast.NodeKind
	params  [][]byte
	context interface{}
}

// Kind implements Node.Kind.
func (n *ShortCodeNode) Kind() gast.NodeKind {
	return n.kind
}

// Dump for ShortCodeNode
func (n *ShortCodeNode) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// --- parser ---

type shortCodeParser struct {
	keyword          string
	minArgc          int
	parserContextKey *parser.ContextKey
}

func (s *shortCodeParser) Trigger() []byte {
	return []byte{'{'}
}

func (s *shortCodeParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	if enabled, ok := pc.Get(EnableShortcodes).(bool); ok && !enabled {
		return nil
	}

	ln, _ := block.PeekLine()
	if !bytes.HasPrefix(ln, []byte("{{")) {
		return nil
	}
	endPos := bytes.Index(ln, []byte("}}"))
	if endPos < 0 {
		return nil
	}
	body := bytes.TrimSpace(ln[2:endPos])
	kwMatch := false
	params := [][]byte{}
	for _, w := range bytes.Split(body, []byte{' '}) {
		if len(w) == 0 {
			continue
		}
		if kwMatch {
			params = append(params, w)
			continue
		}
		if s.keyword == string(w) {
			kwMatch = true
		} else {
			return nil
		}
	}

	usedShortcodes, ok := pc.Get(UsedShortcodesKey).(map[string]bool)
	if !ok {
		usedShortcodes = map[string]bool{}
		pc.Set(UsedShortcodesKey, usedShortcodes)
	}
	usedShortcodes[s.keyword] = true

	if s.minArgc > 0 && len(params) < s.minArgc {
		return nil
	}
	block.Advance(endPos + 2)
	kind, ok := ShortCodeNodeKinds[s.keyword]
	if !ok {
		panic(fmt.Errorf("missing NodeKind for %s", s.keyword))
	}
	node := &ShortCodeNode{
		kind:   kind,
		params: params,
	}
	if s.parserContextKey != nil {
		var elems []ast.Node
		if exitedElems, ok := pc.Get(*s.parserContextKey).([]gast.Node); ok && exitedElems != nil {
			elems = exitedElems
		}
		pc.Set(*s.parserContextKey, append(elems, node))
	}
	return node
}
