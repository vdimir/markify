package view

import (
	"bytes"
	"html/template"
	"net/http"
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
		assert.Equal(t, len(templateNames), len(pageCtxs))
		checkRender(r, pageCtxs)
	}

	r, err := NewRender(http.Dir("../../assets"))
	assert.NoError(t, err)

	checkAllRender(r, []TemplateContext{
		&EditorContext{},
		&PageContext{},
		&StatusContext{},
		&URLPromptContext{},
	})

	checkAllRender(r, []TemplateContext{
		&EditorContext{
			Title:       "Title",
			Msg:         "Msg",
			InitialText: "InitialText",
		},
		&PageContext{
			Title: "Title",
			Body:  template.HTML("<h1>Hello</h1>"),
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
		&URLPromptContext{
			Title:       "Title",
			Msg:         "Msg",
			InitialText: "InitialText",
		},
	})

	checkRender(r, []TemplateContext{
		&PageContext{
			Title: "Title",
			Body:  template.HTML("<h1>Hello</h1>"),
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
