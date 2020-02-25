package engine

import (
	"html/template"

	"github.com/vdimir/markify/store/docstore"
)

type documentWrapper struct {
	key   []byte
	dbDoc *docstore.MdDocument
}

type Document interface {
	Title() string
}

type DocumentRender interface {
	Document
	HTMLBody() template.HTML
}

type DocumentText interface {
	Document
	MdText() []byte
}

type DocumentFull interface {
	Document
	HTMLBody() template.HTML
	MdText() []byte
}

type DocumentSaved interface {
	Key() []byte
}

type DocumentFullSaved interface {
	DocumentFull
	DocumentSaved
}

func (doc *documentWrapper) Title() string {
	return string(doc.dbDoc.Title)
}

func (doc *documentWrapper) HTMLBody() template.HTML {
	return template.HTML(doc.dbDoc.RenderedHTML)
}

func (doc *documentWrapper) MdText() []byte {
	return doc.dbDoc.Text
}

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

func NewUserDocumentData(data []byte) *UserDocumentData {
	return &UserDocumentData{
		Data:             data,
		IsURL:            false,
		EnableShortcodes: true,
		EnableRelImgLink: false,
	}
}
