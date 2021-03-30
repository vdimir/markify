package view

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLTemplateRender(t *testing.T) {

	checkRender := func(r HTMLPageRender, pageCtxs []TemplateContext) {
		for _, ctx := range pageCtxs {
			wr := &bytes.Buffer{}
			err := r.RenderPage(wr, ctx)
			assert.NoError(t, err)
			assert.NotEmpty(t, wr.String())
		}
	}

	checkAllRender := func(r HTMLPageRender, pageCtxs []TemplateContext) {
		checkRender(r, pageCtxs)
	}

	r, err := NewRender(os.DirFS("../app/assets"))
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
