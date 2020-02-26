package docstore

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

type MdMeta struct {
	SrcURL           []byte `bson:",omitempty"`
	Title            []byte `bson:",omitempty"`
	OwnerID          []byte `bson:",omitempty"`
	CreationTime     int    `bson:",omitempty"`
	UpdateTime       int    `bson:",omitempty"`
	MdDocumentParams `bson:",omitempty,inline"`
}

type MdDocument struct {
	MdMeta
	Text         []byte `bson:",omitempty"`
	RenderedHTML []byte `bson:",omitempty"`
}
