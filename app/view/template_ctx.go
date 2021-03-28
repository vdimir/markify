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

// OpenGraphInfo contain Opengraph metadata
type OpenGraphInfo struct {
	Title       string
	Type        string
	URL         string
	Image       string
	Description string
}

// PageContext context for page.html
type PageContext struct {
	Title  string
	Body   template.HTML
	OgInfo *OpenGraphInfo
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

// URLPromptContext context for url_prompt.html
type URLPromptContext struct {
	Title       string
	Msg         string
	InitialText string
}

// Name of the page
func (c *URLPromptContext) Name() string {
	return templateNames[URLPromptTpl]
}
