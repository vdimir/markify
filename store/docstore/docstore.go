package docstore

// DocStore provides interface for stroring MdDocument
// Document can be loaded and updated by parts, to reduce amount of data
type DocStore interface {
	// SaveDocument save new document and return key
	SaveDocument(doc *MdDocument) (DBKey, error)
	// UpdateDocument updates specified document parts in database
	UpdateDocument(key DBKey, parts DocProjection, doc *MdDocument) error
	// LoadDocument load specified document parts from database
	LoadDocument(key DBKey, parts DocProjection, doc *MdDocument) error
}
