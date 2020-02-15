package htmltemplate

import (
	"html/template"
	"io"
	"net/http"

	"github.com/vdimir/markify/util"
)

const templatePrefixDir = "/template"

// HTMLPageRender renders HTML pages from templates
type HTMLPageRender interface {
	RenderPage(wr io.Writer, path string, data interface{}) error
}

// Render contains set of templates and prerenred pages
type Render struct {
	tpl *template.Template
}

// NewRender creates new pages repository
func NewRender(statikFS http.FileSystem) (HTMLPageRender, error) {
	htmlRend := &Render{
		tpl: template.New("root"),
	}

	err := util.WalkFiles(statikFS, templatePrefixDir, func(data []byte, filePath string) error {
		_, err := htmlRend.tpl.New(filePath).Parse(string(data))
		return err
	})

	if err != nil {
		return nil, err
	}

	return htmlRend, nil
}

// RenderPage render page and writes data to wr
func (htmlRend *Render) RenderPage(wr io.Writer, name string, data interface{}) error {
	err := htmlRend.tpl.ExecuteTemplate(wr, name, data)
	return err
}

// DebugRender like Render bute reload templates at every request
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
func (htmlRend *DebugRender) RenderPage(wr io.Writer, path string, data interface{}) error {
	inner, err := NewRender(htmlRend.statikFS)
	if err != nil {
		return err
	}
	err = inner.RenderPage(wr, path, data)
	return err
}
