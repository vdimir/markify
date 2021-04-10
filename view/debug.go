package view

import (
	"io"
	"io/fs"
)

// DebugRender like Render but reloads templates at every request
// For testing purposes only
type DebugRender struct {
	fs.FS
}

// RenderPage render page and writes data to wr. Reloads all templates every time
func (htmlRend *DebugRender) RenderPage(wr io.Writer, tplContext TemplateContext) error {
	inner, err := newView(htmlRend.FS)
	if err != nil {
		return err
	}
	err = inner.RenderPage(wr, tplContext)
	return err
}
