package store

import (
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

// NewBoltDB create bolt.DB with specified buckets.
func NewBoltDB(fileName string, bkts [][]byte, options bbolt.Options) (*bbolt.DB, error) {
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
