package kvstore

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/pkg/errors"
	"github.com/vdimir/markify/store"
	"github.com/vdimir/markify/util"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

const dataBktName = "__data__"
const tsBktName = "__timestamp__"

// Bolt store data in BoldDB
type Bolt struct {
	fileName string
	db       *bolt.DB
}

// NewBoltStorage create Bolt Store.
func NewBoltStorage(fileName string, options bolt.Options) (Store, error) {
	bkts := [][]byte{[]byte(dataBktName), []byte(tsBktName)}
	db, err := store.NewBoltDB(fileName, bkts, options)
	return &Bolt{
		db:       db,
		fileName: fileName,
	}, err
}

// Save data in storage
func (b *Bolt) Save(key []byte, data []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket([]byte(tsBktName)).Get(key) != nil {
			return errors.Errorf("key %v already exists", key)
		}
		return b.update(tx, key, data, []byte(dataBktName))
	})
}

// Update data in storage
func (b *Bolt) Update(key []byte, data []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return b.update(tx, key, data, []byte(dataBktName))
	})
}

// Load data from key, return nil if  key  not found
func (b *Bolt) Load(key []byte) ([]byte, error) {
	var data []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(dataBktName)).Get(key)
		return nil
	})
	return data, err
}

// Timestamp returns key update time in unix nano or zero in key not found
func (b *Bolt) Timestamp(key []byte) (int64, error) {
	var ts int64
	err := b.db.View(func(tx *bolt.Tx) error {
		ts = b.getTimestamp(tx, key)
		return nil
	})
	return ts, err
}

// NewEntry insert data with new key
func (b *Bolt) NewEntry(data []byte) ([]byte, error) {
	key := util.GetUID()
	err := b.Save(key, data)
	return key, err
}

// Close storage
func (b *Bolt) Close() error {
	return b.db.Close()
}

func (b *Bolt) getTimestamp(tx *bbolt.Tx, key []byte) int64 {
	data := tx.Bucket([]byte(tsBktName)).Get(key)
	if data == nil {
		return 0
	}
	var ts int64
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &ts)
	if err != nil {
		panic(err)
	}
	return ts
}

func (b *Bolt) update(tx *bbolt.Tx, key []byte, data []byte, bkt []byte) error {
	var err error
	dataBkt := tx.Bucket(bkt)
	if err = dataBkt.Put(key, data); err != nil {
		return errors.Wrapf(err, "can't put to bucket with %v", key)
	}
	tsBuf := &bytes.Buffer{}
	if err = binary.Write(tsBuf, binary.LittleEndian, time.Now().UnixNano()); err != nil {
		return errors.Wrapf(err, "can't serialize timestamp")
	}
	tsBkt := tx.Bucket([]byte(tsBktName))
	if err = tsBkt.Put(key, tsBuf.Bytes()); err != nil {
		return errors.Wrapf(err, "can't put to bucket with %v", key)
	}
	return err
}
