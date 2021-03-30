package view

import (
	"github.com/vdimir/markify/util"
	"html/template"
	"io"
	"io/fs"
)

const templatePrefixDir = "template"

// HTMLPageRender renders HTML pages from templates
type HTMLPageRender interface {
	RenderPage(wr io.Writer, tplContext TemplateContext) error
}

// Render contains set of templates and pre-rendered pages
type Render struct {
	tpl *template.Template
}

// NewRender creates new pages repository
func NewRender(fs fs.FS) (HTMLPageRender, error) {
	htmlRend := &Render{
		tpl: template.New("root"),
	}

	err := util.WalkFiles(fs, templatePrefixDir, func(data []byte, filePath string) error {
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
	err := htmlRend.tpl.ExecuteTemplate(wr, tplContext.FileName(), tplContext)
	return err
}
