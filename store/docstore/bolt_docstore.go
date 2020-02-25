package docstore

import (
	"path"

	"github.com/vdimir/markify/store"
	"github.com/vdimir/markify/store/kvstore"
	"go.etcd.io/bbolt"
)

var metaBucketName = []byte("__meta__")
var textBucketName = []byte("__text__")

type BoltDocStore struct {
	docStorage    *bbolt.DB
	renderStorage kvstore.Store
}

func NewBoltDocStore(dbPath string) DocStore {
	renderStorage, err := kvstore.NewBoltStorage(path.Join(dbPath, "render.db"), bbolt.Options{})
	if err != nil {
		panic(err)
	}
	bkts := [][]byte{textBucketName, metaBucketName}
	docStorage, err := store.NewBoltDB(path.Join(dbPath, "doc.db"), bkts, bbolt.Options{})
	if err != nil {
		panic(err)
	}
	return &BoltDocStore{
		docStorage:    docStorage,
		renderStorage: renderStorage,
	}
}

// SaveDocument save new document and return key
func (s *BoltDocStore) SaveDocument(doc *MdDocument) (DBKey, error) {
	panic("")
}

func (s *BoltDocStore) LoadDocument(key DBKey, parts DocProjection, doc *MdDocument) error {
	panic("")
}

func (s *BoltDocStore) UpdateDocument(key DBKey, parts DocProjection, doc *MdDocument) error {
	panic("")
}

func (s *BoltDocStore) Close() error {
	var err error
	err = s.docStorage.Close()
	if err != nil {
		return err
	}
	err = s.renderStorage.Close()
	if err != nil {
		return err
	}
	return nil
}
