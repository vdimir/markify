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

// NewDebugRender creates debug HTMLPageRender
func NewDebugRender(fs fs.FS) (HTMLPageRender, error) {
	_, err := NewRender(fs)
	return &DebugRender{fs}, err
}

// RenderPage render page and writes data to wr. Realoads all templates every time
func (htmlRend *DebugRender) RenderPage(wr io.Writer, tplContext TemplateContext) error {
	inner, err := NewRender(htmlRend)
	if err != nil {
		return err
	}
	err = inner.RenderPage(wr, tplContext)
	return err
}
