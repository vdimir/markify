package view

import (
	"bytes"
	"html/template"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLTemplateRender(t *testing.T) {

	checlAllRender := func(r HTMLPageRender, pageCtxs []TemplateContext) {
		assert.Equal(t, len(templateNames), len(pageCtxs))

		for _, ctx := range pageCtxs {
			wr := &bytes.Buffer{}
			err := r.RenderPage(wr, ctx)
			assert.NoError(t, err)
			assert.NotEmpty(t, wr.String())
		}
	}

	r, err := NewRender(http.Dir("../../assets"))
	assert.NoError(t, err)

	checlAllRender(r, []TemplateContext{
		&EditorContext{},
		&PageContext{},
		&StatusContext{},
		&URLPromptContext{},
	})

	checlAllRender(r, []TemplateContext{
		&EditorContext{
			Title:       "Title",
			Msg:         "Msg",
			InitialText: "InitialText",
		},
		&PageContext{
			Title: "Title",
			Body:  template.HTML("<h1>Hello</h1>"),
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
}
