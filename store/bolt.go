package store

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

const dataBktName = "__data__"
const metaBktName = "__metadata__"

// Bolt store data in BoldDB
type Bolt struct {
	fileName string
	db       *bolt.DB
}

// NewBoltStorage create Bolt Store.
func NewBoltStorage(fileName string) (*Bolt, error) {
	bkts := [][]byte{[]byte(dataBktName), []byte(metaBktName)}
	db, err := newBoltWithBuckets(fileName, bkts, bbolt.Options{})
	return &Bolt{
		db:       db,
		fileName: fileName,
	}, err
}

// Save data in storage
func (b *Bolt) SetBlob(key string, reader io.Reader, meta map[string]string, _ time.Duration) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "can't read data from reader")
	}
	metadata, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket([]byte(metaBktName)).Put([]byte(key), metadata); err != nil {
			return err
		}
		return tx.Bucket([]byte(dataBktName)).Put([]byte(key), data)
	})
}

func (b *Bolt) GetBlob(key string) (io.Reader, error) {
	var data []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(dataBktName)).Get([]byte(key))
		return nil
	})
	return bytes.NewReader(data), err
}

func (b *Bolt) GetMeta(key string) (map[string]string, error) {
	var data []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(metaBktName)).Get([]byte(key))
		return nil
	})
	if err != nil {
		return nil, err
	}
	var meta map[string]string
	err = json.Unmarshal(data, &meta)
	return meta, err
}

func (b *Bolt) DeleteBlob(key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket([]byte(metaBktName)).Delete([]byte(key)); err != nil {
			return err
		}
		return tx.Bucket([]byte(dataBktName)).Delete([]byte(key))
	})
}

// Close storage
func (b *Bolt) Close() error {
	return b.db.Close()
}

// newBoltWithBuckets create bolt.DB with specified buckets.
func newBoltWithBuckets(fileName string, bkts [][]byte, options bbolt.Options) (*bbolt.DB, error) {
	db, err := bbolt.Open(fileName, 0600, &options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to make boltdb for %s", fileName)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		for _, bkt := range bkts {
			if _, e := tx.CreateBucketIfNotExists(bkt); e != nil {
				return errors.Wrapf(e, "failed to create top level bucket %s", bkt)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize boltdb db %q buckets", fileName)
	}
	return db, nil
}
