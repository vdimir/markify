package view

import "html/template"

// PageType HTML template identifier
type PageType int

const (
	// EditorTpl template "editor.html"
	EditorTpl PageType = iota
	// PageTpl template "page.html"
	PageTpl
	// StatusTpl template "status.html"
	StatusTpl
	// URLPromptTpl template "url_prompt.html"
	URLPromptTpl
	// Number of templates
	_TplCount
)

var templateNames = map[PageType]string{
	EditorTpl:    "editor.html",
	PageTpl:      "page.html",
	StatusTpl:    "status.html",
	URLPromptTpl: "url_prompt.html",
}

// TemplateContext provides context for template
type TemplateContext interface {
	// Name of template file to render
	Name() string
}

// EditorContext context for editor.html
type EditorContext struct {
	Title       string
	Msg         string
	InitialText string
}

// Name of the page
func (c *EditorContext) Name() string {
	return templateNames[EditorTpl]
}

// PageContext context for page.html
type PageContext struct {
	Title string
	Body  template.HTML
}

// Name of the page
func (c *PageContext) Name() string {
	return templateNames[PageTpl]
}

// StatusContext context for status.html
type StatusContext struct {
	Title     string
	HeaderMsg string
	Msg       string
}

// Name of the page
func (c *StatusContext) Name() string {
	return templateNames[StatusTpl]
}

// URLPromptContext context for utl_prompt.html
type URLPromptContext struct {
	Title       string
	Msg         string
	InitialText string
}

// Name of the page
func (c *URLPromptContext) Name() string {
	return templateNames[URLPromptTpl]
}
