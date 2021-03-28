package docstore

// DocStore provides interface for storing MdDocument
type DocStore interface {
	// SaveDocument save new document and return key
	SaveDocument(doc *MdDocument) (DBKey, error)
	// UpdateDocument updates specified document parts in database
	UpdateDocument(key DBKey, doc *MdDocument) error
	// LoadDocument load specified document parts from database
	LoadDocument(key DBKey, doc *MdDocument) error
}
