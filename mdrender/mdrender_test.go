package md

import (
	"bytes"
	"net/url"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/stretchr/testify/require"
	"github.com/vdimir/markify/testutil"
)

const testDataPath = "../testdata"

func TestRenderSimple(t *testing.T) {
	mdData := testutil.MustReadData(t, path.Join(testDataPath, "page.md"))
	data := mustRenderMd(t, mdData)
	assert.NotEmpty(t, data)
}

func TestRelativeImgLink(t *testing.T) {
	rndr, err := NewRender()
	require.NoError(t, err)

	mdData := testutil.MustReadData(t, path.Join(testDataPath, "page.md"))

	pageURL, err := url.ParseRequestURI("http://aaa.com/foo/bar.md")
	require.NoError(t, err)

	buf, _, err := rndr.Render(mdData, &Options{BaseURL: pageURL})
	assert.NoError(t, err)

	htmlPrs, err := html.Parse(buf)
	require.NoError(t, err)

	imgExpectedTags := map[string]string{
		"img":  "http://aaa.com/foo/static/foo.png",
		"img2": "http://aaa.com/foo/bar.png",
	}

	traverseHTMLNodes(htmlPrs, func(n *html.Node) {
		if n.Type != html.ElementNode || n.Data != "img" {
			return
		}
		var src, alt string
		for _, a := range n.Attr {
			if a.Key == "src" {
				src = a.Val
			}
			if a.Key == "alt" {
				alt = a.Val
			}
		}
		if expectedSrc, ok := imgExpectedTags[alt]; ok {
			assert.Equal(t, expectedSrc, src)
			delete(imgExpectedTags, alt)
		}
	})
	assert.Empty(t, imgExpectedTags)
}

func TestTweetShortcode(t *testing.T) {
	mdData := []byte("" +
		"Foo\nTweet {{ tweet     1224047348109795330  }} ok?\n" +
		"\n")

	expNodes := map[string]*kvPair{
		"blockquote": nil,
		"a":          &kvPair{"href", "https://twitter.com/AlgebraFact/status/1224047348109795330"},
	}
	checkHTMLAttrs(t, mustRenderMd(t, mdData), expNodes)
}

func TestTweetShortcodeInCode(t *testing.T) {
	mdData := []byte("" +
		"Foo\nTweet `{{ tweet 1224047348109795330 }}`\n" +
		"Tweet { tweet 1224047348109795311 }}\n" +
		"")

	expectedSubslices := map[string]bool{
		"blockquote":                      false,
		"{{ tweet 1224047348109795330 }}": true,
		"{ tweet 1224047348109795311 }}":  true,
	}

	checkContaining(t, mustRenderMd(t, mdData), expectedSubslices)
}

func TestTweetShortcodeError(t *testing.T) {
	mdData := []byte("" +
		"Foo\nTweet {{ tweet 000 }}\n" +
		"Tweet {{ tweet aaa }}\n" +
		"")

	expectedSubslices := map[string]bool{
		"blockquote":               false,
		"{{ tweet 000 }}":          false,
		"{{ tweet aaa }}":          false,
		"Unable to load tweet 000": true,
		"Unable to load tweet aaa": true,
	}
	checkContaining(t, mustRenderMd(t, mdData), expectedSubslices)
}

func TestInstagramShortcode(t *testing.T) {
	t.Skip()
	
	mdData := []byte("Foo\nInsagram {{ instagram B7gs_jFKWA0  }} ok?\n")

	expNodes := map[string]*kvPair{
		"blockquote": nil,
		"a":          &kvPair{"href", "https://www.instagram.com/p/B7gs_jFKWA0"},
	}
	checkHTMLAttrs(t, mustRenderMd(t, mdData), expNodes)
}

func TestInstagramShortcodeError(t *testing.T) {
	expectedSubslices := func(s string) map[string]bool {
		return map[string]bool{s: true}
	}
	var data []byte

	data = mustRenderMd(t, []byte("Foo\nInsagram {{ instagram x;//sdfs }} ok?\n"))
	checkContaining(t, data, expectedSubslices("Wrong arguments for instagram"))

	data = mustRenderMd(t, []byte("Foo\nInsagram {{ instagram }} ok?\n\n"))
	checkContaining(t, data, expectedSubslices("Wrong arguments for instagram"))
}

func TestGistShortcode(t *testing.T) {
	mdData := []byte("Foo\nGist {{ gist spf13 7896402 }} ok?\n\n")

	expNodes := map[string]*kvPair{
		"script": &kvPair{"src", "https://gist.github.com/spf13/7896402.js"},
	}
	checkHTMLAttrs(t, mustRenderMd(t, mdData), expNodes)
}

func TestGistShortcodeError(t *testing.T) {
	mdData := []byte("Foo\nGist {{ gist spf13  }} ok?\n")

	expectedSubslices := map[string]bool{
		"Unable to display gist": true,
	}
	checkContaining(t, mustRenderMd(t, mdData), expectedSubslices)
}

func TestDisableShortcode(t *testing.T) {
	mdData := []byte("" +
		"Tweet {{ tweet 1224047348109795330 }}\n" +
		"")

	expectedSubslices := map[string]bool{
		"blockquote":                      false,
		"a":                               false,
		"{{ tweet 1224047348109795330 }}": true,
	}

	rndr, err := NewRender()
	require.NoError(t, err)

	buf, _, err := rndr.Render(mdData, &Options{DisableShortcodes: true})
	require.NoError(t, err)

	checkContaining(t, buf.Bytes(), expectedSubslices)
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

		rndr, err := NewRender()
		require.NoError(t, err)

		_, ctx, err := rndr.Render(mdData, nil)
		require.NoError(t, err)
		require.Equal(t, &PagePreviewText{"header 11", ""}, ctx.Get(MdTitleKey))
	}
	{
		mdData := []byte("" +
			"# header 11\n" +
			"text 123\n" +
			"456" + "\n\n" +
			"text 789" + "\n" +
			"# header 22\n" +
			"\n")

		rndr, err := NewRender()
		require.NoError(t, err)

		_, ctx, err := rndr.Render(mdData, nil)
		require.NoError(t, err)
		require.Equal(t, &PagePreviewText{"header 11", "text 123 456"}, ctx.Get(MdTitleKey))
	}
	{
		mdData := []byte("" +
			"text 11\n\n" +
			"text 22\n\n" +
			"text 33\n\n" +
			"\n")

		rndr, err := NewRender()
		require.NoError(t, err)

		_, ctx, err := rndr.Render(mdData, nil)
		require.NoError(t, err)
		require.Equal(t, &PagePreviewText{"", "text 11"}, ctx.Get(MdTitleKey))
	}
	{
		mdData := []byte("" +
			"# text 11\n\n" +
			"## text 22\n\n" +
			"# text 33\n\n" +
			"\n")

		rndr, err := NewRender()
		require.NoError(t, err)

		_, ctx, err := rndr.Render(mdData, nil)
		require.NoError(t, err)
		require.Equal(t, &PagePreviewText{"text 11", ""}, ctx.Get(MdTitleKey))
	}
}

// --- helpers ---

func mustRenderMd(t *testing.T, mdData []byte) []byte {
	rndr, err := NewRender()
	require.NoError(t, err)

	buf, _, err := rndr.Render(mdData, nil)
	require.NoError(t, err)
	return buf.Bytes()
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
func checkContaining(t *testing.T, data []byte, expectedSubslices map[string]bool) {
	for k, shouldContain := range expectedSubslices {
		contains := bytes.Contains(data, []byte(k))
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
