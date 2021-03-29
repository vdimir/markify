package view

import (
	"io"
	"net/http"
)

// DebugRender like Render but reloads templates at every request
// For testing purposes only
type DebugRender struct {
	statikFS http.FileSystem
}

// NewDebugRender creates debug HTMLPageRender
func NewDebugRender(statikFS http.FileSystem) (HTMLPageRender, error) {
	_, err := NewRender(statikFS)
	return &DebugRender{
		statikFS: statikFS,
	}, err
}

// RenderPage render page and writes data to wr. Realoads all templates every time
func (htmlRend *DebugRender) RenderPage(wr io.Writer, tplContext TemplateContext) error {
	inner, err := NewRender(htmlRend.statikFS)
	if err != nil {
		return err
	}
	err = inner.RenderPage(wr, tplContext)
	return err
}
