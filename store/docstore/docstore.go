package docstore

type DocStore interface {
	SaveDocument(doc *MdDocument) (DBKey, error)
	UpdateDocument(key DBKey, parts DocProjection, doc *MdDocument) error
	LoadDocument(key DBKey, parts DocProjection, doc *MdDocument) error
}
