package view

import (
	"embed"
	"github.com/vdimir/markify/util"
	"html/template"
	"io"
	"io/fs"
	"os"
)

//go:embed template/*
var embeddedTemplateFS embed.FS

// HTMLPageView renders HTML pages from templates
type HTMLPageView interface {
	RenderPage(wr io.Writer, tplContext TemplateContext) error
}

// Render contains set of templates and pre-rendered pages
type Render struct {
	tpl *template.Template
}

func NewView(filePath string) (HTMLPageView, error) {
	if filePath != "" {
		return &DebugRender{os.DirFS(filePath)}, nil
	}
	subFs, err := fs.Sub(embeddedTemplateFS, "template")
	if err != nil {
		return nil, err
	}
	return newView(subFs)
}

// RenderPage render page and writes data to wr
func (htmlRend *Render) RenderPage(wr io.Writer, tplContext TemplateContext) error {
	err := htmlRend.tpl.ExecuteTemplate(wr, tplContext.FileName(), tplContext)
	return err
}

func newView(fs fs.FS) (HTMLPageView, error) {
	htmlRend := &Render{
		tpl: template.New("root"),
	}

	err := util.WalkFiles(fs, ".", func(data []byte, filePath string) error {
		_, err := htmlRend.tpl.New(filePath).Parse(string(data))
		return err
	})

	if err != nil {
		return nil, err
	}

	return htmlRend, nil
}

