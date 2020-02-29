package engine

import (
	"html/template"

	"github.com/vdimir/markify/store/docstore"
)

type documentWrapper struct {
	key   []byte
	dbDoc *docstore.MdDocument
}

// Document base interface
type Document interface {
	Title() string
}

// DocumentRender provides HTML formatted document
type DocumentRender interface {
	Document
	HTMLBody() template.HTML
}

// DocumentText provides raw document text
type DocumentText interface {
	Document
	MdText() []byte
}

// DocumentFull contains both raw text and HTML render
type DocumentFull interface {
	DocumentRender
	DocumentText
}

// DocumentSaved contains key in database
type DocumentSaved interface {
	Key() []byte
}

// DocumentFullSaved contains all document data saved in db
type DocumentFullSaved interface {
	DocumentFull
	DocumentSaved
}

// Title of document (may be emtry)
func (doc *documentWrapper) Title() string {
	return string(doc.dbDoc.Title)
}

// HTMLBody rendered document
func (doc *documentWrapper) HTMLBody() template.HTML {
	if doc.dbDoc.RenderedHTML == nil {
		panic("RenderedHTML is nil")
	}
	return template.HTML(doc.dbDoc.RenderedHTML)
}

// MdText raw document text
func (doc *documentWrapper) MdText() []byte {
	if doc.dbDoc.Text == nil {
		panic("Text is nil")
	}
	return doc.dbDoc.Text
}

// Key in database
func (doc *documentWrapper) Key() []byte {
	if doc.key == nil {
		panic("key is nil")
	}
	return doc.key
}

// UserDocumentData contains untreated user input used to create new document
type UserDocumentData struct {
	Data             []byte
	IsURL            bool
	EnableShortcodes bool
	EnableRelImgLink bool
}

// NewUserDocumentData creates default UserDocumentData
func NewUserDocumentData(data []byte) *UserDocumentData {
	return &UserDocumentData{
		Data:             data,
		IsURL:            false,
		EnableShortcodes: true,
		EnableRelImgLink: false,
	}
}
