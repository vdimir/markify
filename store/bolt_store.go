package store

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/vdimir/markify/util"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

const dataBktName = "data"
const tsBktName = "timestamp"

// Bolt store data in BoldDB
type Bolt struct {
	fileName string
	db       *bolt.DB
}

// NewBoltStorage create Bolt Store.
func NewBoltStorage(fileName string, options bolt.Options) (*Bolt, error) {
	db, err := bolt.Open(fileName, 0600, &options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to make boltdb for %s", fileName)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		if _, e := tx.CreateBucketIfNotExists([]byte(dataBktName)); e != nil {
			return errors.Wrapf(e, "failed to create top level bucket %s", dataBktName)
		}
		if _, e := tx.CreateBucketIfNotExists([]byte(tsBktName)); e != nil {
			return errors.Wrapf(e, "failed to create top level bucket %s", tsBktName)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize boltdb db %q buckets", fileName)
	}
	return &Bolt{
		db:       db,
		fileName: fileName,
	}, nil
}

// Save data in storage
func (b *Bolt) Save(key []byte, data []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		var err error
		if err = tx.Bucket([]byte(dataBktName)).Put(key, data); err != nil {
			return errors.Wrapf(err, "can't put to bucket with %v", key)
		}
		tsBuf := &bytes.Buffer{}
		if err = binary.Write(tsBuf, binary.LittleEndian, time.Now().UnixNano()); err != nil {
			return errors.Wrapf(err, "can't serialize timestamp")
		}
		if err = tx.Bucket([]byte(tsBktName)).Put(key, tsBuf.Bytes()); err != nil {
			return errors.Wrapf(err, "can't put to bucket with %v", key)
		}
		return err
	})

	return err
}

// Load data from key, return nil if  key  not found
func (b *Bolt) Load(key []byte) ([]byte, error) {
	var data []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(dataBktName)).Get(key)
		if data == nil {
			return nil
		}
		return nil
	})
	return data, err
}

// Timestamp returns key update time in unix nano or zero in key not found
func (b *Bolt) Timestamp(key []byte) (int64, error) {
	var data []byte
	var ts int64

	err := b.db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(tsBktName)).Get(key)
		if data == nil {
			return nil
		}
		err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &ts)
		if err != nil {
			return err
		}
		return nil
	})
	return ts, err
}

// NewKey insert data with new key
func (b *Bolt) NewKey(data []byte) ([]byte, error) {
	key := util.GetUID()
	b.Save(key, data)
	return key, nil
}

// Close storage
func (b *Bolt) Close() error {
	return b.db.Close()
}
