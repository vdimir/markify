package view

import "html/template"

// TemplateContext provides context for template
type TemplateContext interface {
	// FileName with template file to render
	FileName() string
}

// EditorContext context for editor.html
type EditorContext struct {
	Title       string
	Msg         string
	InitialText string
}

// Name of the page
func (c *EditorContext) FileName() string {
	return "editor.html"
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
	CreateTime string
}

// Name of the page
func (c *PageContext) FileName() string {
	return "page.html"
}

// StatusContext context for status.html
type StatusContext struct {
	Title     string
	HeaderMsg string
	Msg       string
}

// Name of the page
func (c *StatusContext) FileName() string {
	return "status.html"
}
