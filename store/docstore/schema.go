package docstore

import (
	"net/url"
	"time"
)

type DBKey []byte

type DocProjection int32

const (
	ProjMeta DocProjection = 2 >> iota
	ProjText
	ProjRender
	ProjAll = ProjMeta | ProjText | ProjRender
)

type MdDocumentParams struct {
	EnableShortcodes bool
}

type MdDocument struct {
	Key          DBKey            // *
	SrcURL       *url.URL         // *
	Title        []byte           // *
	OwnerID      []byte           // * Meta
	CreationTime time.Time        // *
	UpdateTime   time.Time        // *
	Params       MdDocumentParams // *
	Text         []byte           // Text
	RenderedHTML []byte           // Render
}
