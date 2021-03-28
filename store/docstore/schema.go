package docstore

// DBKey  is key in database
type DBKey []byte

// MdMeta contains document meta info
type MdMeta struct {
	Title            []byte `bson:",omitempty"`
	Description      []byte `bson:",omitempty"`
	OwnerID          []byte `bson:",omitempty"`
	CreationTime     int64  `bson:",omitempty"`
	UpdateTime       int64  `bson:",omitempty"`
}

// MdDocument represent markdown document
type MdDocument struct {
	MdMeta
	Text         []byte `bson:",omitempty"`
	RenderedHTML []byte `bson:",omitempty"`
}
