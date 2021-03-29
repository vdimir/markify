package markdown

import (
	"bytes"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/stretchr/testify/require"
	"github.com/vdimir/markify/testutil"
)

const testDataPath = "../../testdata"

func TestRenderSimple(t *testing.T) {
	mdData := testutil.MustReadData(t, path.Join(testDataPath, "page.md"))
	data := mustRenderMd(t, mdData)
	assert.NotEmpty(t, data)
}

func TestTOCShortcode(t *testing.T) {
	mdData := []byte("" +
		"## Header *text **text*** 0.1\n" +
		"# Header 1\n" +
		"Table of contents:\n" +
		"{{ toc }}\n" +
		"# Header2\n" +
		"Text\n" +
		"## Header 2.1\n" +
		"Other Text\n" +
		"## Header 2.2\n" +
		"### *Header* 2.2.1\n" +
		"### Header 2.2.2\n" +
		"## __Header__ 2.3\n" +
		"# Header 3\n" +
		"# Header 4\n" +
		"### Header 4.1.1\n" +
		"### Header 4.1.2\n" +
		"## Header 4.2\n" +
		"")
	data := mustRenderMd(t, mdData)

	checkContaining(t, data, map[string]bool{"{{ toc }}": false, "<nav class=\"toc-block\">": true})
	// TODO generated html correctness
}

func TestTitleExtractor(t *testing.T) {
	{
		mdData := []byte("# header 11")

		rndr := NewConverter()

		doc, err := rndr.Convert(mdData)
		require.NoError(t, err)
		require.Equal(t, "header 11", doc.Title)
		require.Equal(t, "", doc.Preview)
	}
	{
		mdData := []byte("" +
			"# header 11\n" +
			"text 123\n" +
			"456" + "\n\n" +
			"text 789" + "\n" +
			"# header 22\n" +
			"\n")
		rndr := NewConverter()

		doc, err := rndr.Convert(mdData)
		require.NoError(t, err)
		require.Equal(t, "header 11", doc.Title)
		require.Equal(t, "text 123 456", doc.Preview)
	}
	{
		mdData := []byte("" +
			"text 11\n\n" +
			"text 22\n\n" +
			"text 33\n\n" +
			"\n")

		rndr := NewConverter()

		doc, err := rndr.Convert(mdData)
		require.NoError(t, err)
		require.Equal(t, "", doc.Title)
		require.Equal(t, "text 11", doc.Preview)
	}
	{
		mdData := []byte("" +
			"# text 11\n\n" +
			"## text 22\n\n" +
			"# text 33\n\n" +
			"\n")

		rndr := NewConverter()

		doc, err := rndr.Convert(mdData)
		require.NoError(t, err)
		require.Equal(t, "text 11", doc.Title)
		require.Equal(t, "", doc.Preview)
	}
}

// --- helpers ---

func mustRenderMd(t *testing.T, mdData []byte) string {
	rndr := NewConverter()
	doc, err := rndr.Convert(mdData)
	require.NoError(t, err)
	return doc.Body
}

func traverseHTMLNodes(root *html.Node, process func(n *html.Node)) {
	var walker func(*html.Node)
	walker = func(n *html.Node) {
		process(n)

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(root)
}

// checkContaining checks that data contains (un)expected  substirings
// expectedSubslices values should be true if string required
// and false if string should be absent
func checkContaining(t *testing.T, data string, expectedSubslices map[string]bool) {
	for k, shouldContain := range expectedSubslices {
		contains := strings.Contains(data, k)
		msg := "does contains"
		if shouldContain {
			msg = "does not contains"
		}
		assert.Truef(t, contains == shouldContain, "%q %s %q", string(data), msg, k)
	}
}

type kvPair struct {
	key string
	val string
}

// checkHTMLAttrs checks that specified tags are presented with given attributes.
// If attrinues is nil, only tag existence checked
func checkHTMLAttrs(t *testing.T, data []byte, expNodes map[string]*kvPair) {
	root, err := html.Parse(bytes.NewReader(data))
	require.NoError(t, err)

	traverseHTMLNodes(root, func(n *html.Node) {
		if n.Type != html.ElementNode {
			return
		}
		expectedAttrs, ok := expNodes[n.Data]
		if !ok {
			return
		}
		if expectedAttrs == nil {
			delete(expNodes, n.Data)
			return
		}
		for _, a := range n.Attr {
			attrFound := a.Key == expectedAttrs.key && strings.HasPrefix(a.Val, expectedAttrs.val)
			if attrFound {
				delete(expNodes, n.Data)
				break
			}
		}
	})

	for k, v := range expNodes {
		if v == nil {
			assert.Failf(t, "Test failed", "<%s> tag not found in %s", k, data)
		} else {
			assert.Failf(t, "Test failed", "<%s %s=%q ...> tag not found in %s", k, v.key, v.val, data)
		}
	}
}
