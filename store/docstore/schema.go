package docstore

// DBKey  is key in database
type DBKey []byte

// DocProjection specifies part of document to work with
type DocProjection int32

const (
	// ProjMeta contains document meta informatrion
	ProjMeta DocProjection = 2 >> iota
	// ProjText contains raw document text
	ProjText
	// ProjRender contains rendered HTML
	ProjRender
	// ProjAll contains full document data
	ProjAll = ProjMeta | ProjText | ProjRender
)

// MdDocumentParams contains document user parametes
type MdDocumentParams struct {
	EnableShortcodes bool
}

// MdMeta contains document meta info
type MdMeta struct {
	SrcURL           []byte `bson:",omitempty"`
	Title            []byte `bson:",omitempty"`
	OwnerID          []byte `bson:",omitempty"`
	CreationTime     int    `bson:",omitempty"`
	UpdateTime       int    `bson:",omitempty"`
	MdDocumentParams `bson:",omitempty,inline"`
}

// MdDocument represent markdown document
type MdDocument struct {
	MdMeta
	Text         []byte `bson:",omitempty"`
	RenderedHTML []byte `bson:",omitempty"`
}
