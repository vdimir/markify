package view

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLTemplateRender(t *testing.T) {

	checkRender := func(r HTMLPageView, pageCtxs []TemplateContext) {
		for _, ctx := range pageCtxs {
			wr := &bytes.Buffer{}
			err := r.RenderPage(wr, ctx)
			assert.NoError(t, err)
			assert.NotEmpty(t, wr.String())
		}
	}

	checkAllRender := func(r HTMLPageView, pageCtxs []TemplateContext) {
		checkRender(r, pageCtxs)
	}

	r, err := NewView("template")
	require.NoError(t, err)

	checkAllRender(r, []TemplateContext{
		&EditorContext{},
		&PageContext{},
		&StatusContext{},
	})

	checkAllRender(r, []TemplateContext{
		&EditorContext{
			Title:       "Title",
			Msg:         "Msg",
			InitialText: "InitialText",
		},
		&PageContext{
			Title: "Title",
			Body:  "<h1>Hello</h1>",
			OgInfo: &OpenGraphInfo{
				Title:       "ogtitle",
				URL:         "http://markify.dev/foobar",
				Image:       "",
				Type:        "article",
				Description: "article description",
			},
		},
		&StatusContext{
			Title:     "Title",
			HeaderMsg: "HeaderMsg",
			Msg:       "Msg",
		},
	})

	checkRender(r, []TemplateContext{
		&PageContext{
			Title: "Title",
			Body:  "<h1>Hello</h1>",
			OgInfo: &OpenGraphInfo{
				Title:       "ogtitle",
				URL:         "http://markify.dev/foobar",
				Image:       "",
				Type:        "article",
				Description: "article description",
			},
		},
	})
}
