package view

import (
	"html/template"
	"io"
	"net/http"

	"github.com/vdimir/markify/util"
)

const templatePrefixDir = "/template"

// HTMLPageRender renders HTML pages from templates
type HTMLPageRender interface {
	RenderPage(wr io.Writer, tplContext TemplateContext) error
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
func (htmlRend *Render) RenderPage(wr io.Writer, tplContext TemplateContext) error {
	err := htmlRend.tpl.ExecuteTemplate(wr, tplContext.Name(), tplContext)
	return err
}
