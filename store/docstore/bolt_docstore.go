package docstore

import (
	"path"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/pkg/errors"
	"github.com/vdimir/markify/util"

	"github.com/vdimir/markify/store"
	"github.com/vdimir/markify/store/kvstore"
	"go.etcd.io/bbolt"
)

var metaBucketName = []byte("__meta__")
var textBucketName = []byte("__text__")

type boltDocStore struct {
	docStorage    *bbolt.DB
	renderStorage kvstore.Store
}

// NewBoltDocStore create boltDocStore. Database files  store in dbPath folder
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
	return &boltDocStore{
		docStorage:    docStorage,
		renderStorage: renderStorage,
	}
}

// SaveDocument save new document and return key
func (s *boltDocStore) SaveDocument(doc *MdDocument) (DBKey, error) {
	if doc.CreationTime == 0 || doc.UpdateTime == 0 {
		return nil, errors.Errorf("fields CreationTime and UpdateTime required")
	}
	if doc.Text == nil || len(doc.Text) == 0 {
		return nil, errors.Errorf("empty text")
	}

	key := util.GetUID()
	blob, err := bson.Marshal(doc)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal error in save")
	}
	err = s.docStorage.Update(func(tx *bbolt.Tx) error {
		err := tx.Bucket(metaBucketName).Put(key, blob)
		if err != nil {
			return errors.Wrapf(err, "cannot put data to bolt")
		}
		return nil
	})
	return key, err
}

// LoadDocument load specified document parts from database
func (s *boltDocStore) LoadDocument(key DBKey, doc *MdDocument) error {
	err := s.docStorage.View(func(tx *bbolt.Tx) error {
		blob := tx.Bucket(metaBucketName).Get(key)
		if blob == nil {
			return errors.Errorf("key %v not found", key)
		}
		err := bson.Unmarshal(blob, doc)
		if err != nil {
			return errors.Wrapf(err, "unmarshal error in load")
		}
		return nil
	})
	return err
}

// UpdateDocument updates specified document parts in database
func (s *boltDocStore) UpdateDocument(key DBKey, doc *MdDocument) error {
	blob, err := bson.Marshal(doc)
	if err != nil {
		return err
	}
	err = s.docStorage.Update(func(tx *bbolt.Tx) error {
		err := tx.Bucket(metaBucketName).Put(key, blob)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *boltDocStore) Close() error {
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
