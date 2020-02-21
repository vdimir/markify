package htmltemplate

import "html/template"

const (
	// URLInputPage identifier for index page
	URLInputPage = "url_prompt.html"
	// StatusTemplatePage identifier for status page
	StatusTemplatePage = "status.html"
	// ContentTemplatePage identifier for page with content such article
	ContentTemplatePage = "page.html"
	// TextEditorTemplate identifier for page with text-editing area
	TextEditorTemplate = "editor.html"
)

// InfoPageContext cotains info data pages with Title and Header
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
