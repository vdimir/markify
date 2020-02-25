package docstore

import (
	"go.mongodb.org/mongo-driver/bson"
)

type dbDocMeta struct {
	SrcURL           []byte
	Title            []byte
	OwnerID          []byte
	CreationTime     int
	UpdateTime       int
	EnableShortcodes bool
}

type dbDocText struct {
	inner []byte
}

func (doc *MdDocument) MetaSerializer() bson.Marshaler {
	return nil
}

// func (d *dbDocMeta) DestBucket() []byte {
// 	return []byte(store.DataBktName)
// }

// func (d dbText) DestBucket() []byte {
// 	return textBucketName
// }

func (doc *MdDocument) MarshalBSON() ([]byte, error) {
	docHeader := &dbDocMeta{
		Title:        doc.Title,
		OwnerID:      doc.OwnerID,
		CreationTime: doc.CreationTime.Second(),
		UpdateTime:   doc.UpdateTime.Second(),
	}
	if doc.SrcURL != nil {
		docHeader.SrcURL = []byte(doc.SrcURL.String())
	}
	return bson.Marshal(docHeader)
}
