package db

import (
	bolt "go.etcd.io/bbolt"
	"go.etcd.io/bbolt/errors"
)

var defaultBucket = []byte("default")

// wrapper around the bolt.DB
type database struct {
	db *bolt.DB
}

func NewDatabase(dbPath string) (*database, func() error, error) {
	boltDB, err := bolt.Open(dbPath, 0600, nil)

	if err != nil {
		return nil, nil, err
	}

	dDB := &database{boltDB}

	if err := dDB.createDefaultBucket(); err != nil {
		boltDB.Close()
		return nil, nil, err
	}

	return dDB, boltDB.Close, nil
}

func (d *database) createDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		return err
	})
}

// set key in the default bucket
func (d *database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), value)
	})
}

// get key from the default bucket
func (d *database) GetKey(key string) ([]byte, error) {
	var value []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(defaultBucket))
		if bucket == nil {
			return errors.ErrBucketNotFound
		}

		value = bucket.Get([]byte(key))
		return nil
	})
	return value, err
}
