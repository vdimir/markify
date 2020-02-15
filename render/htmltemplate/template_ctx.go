package htmltemplate

import "html/template"

const (
	// IndexTemplatePage identifier for index page
	IndexTemplatePage = "index.html"
	// StatusTemplatePage identifier for status page
	StatusTemplatePage = "status.html"
	// ContentTemplatePage identifier for page with content such article
	ContentTemplatePage = "page.html"
	// TextEditorTemplate identifier for page with text-editing area
	TextEditorTemplate = "editor.html"
)

// InfoPageContext cotains info data for info_page.html and index.html
type InfoPageContext struct {
	Title     string
	HeaderMsg string
	MainMsg   string
}

// MdPageContext cotains info data for page.html
type MdPageContext struct {
	Title string
	Body  template.HTML
}
