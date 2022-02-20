package database

import (
	"errors"

	"github.com/boltdb/bolt"
)

type buckets interface {
	AddBucket(string)
	SetBucket(string)
	ListBuckets() ([]string, error)
	DeleteBucket(string) error
}

// BoltDB - BoltDB struct
type BoltDB struct {
	*bolt.DB
	bucketName []byte
	buckets
}

// Conn godoc - Returns the underlying database connection
func (db *BoltDB) Conn() interface{} {
	return db
}

// openBolt godoc - Creates a new BoltDB instance
func openBolt(fileName string) (db *BoltDB, err error) {

	db = &BoltDB{}
	// open new bolt file
	db.DB, err = bolt.Open(fileName, 0600, nil)
	return
}

// CloseDB godoc - Closes the database
func (db *BoltDB) CloseDB() {
	db.Close()
}

// Add godoc - Adds a new key/value pair to the database
func (db *BoltDB) Add(key string, value string, args ...interface{}) error {

	// add the key-value pair to the boltdb file
	err := db.Update(func(tx *bolt.Tx) error {

		// create a bucket if it doesn't exist
		b, err := tx.CreateBucketIfNotExists(db.bucketName)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
	return err
}

// Get godoc - Gets a value from the database
func (db *BoltDB) Get(key string) (string, error) {

	var value string
	err := db.View(func(tx *bolt.Tx) error {

		// open the bucket
		b := tx.Bucket(db.bucketName)
		if b == nil {
			return errors.New("bucket not found")
		}

		// get value from db
		value = string(b.Get([]byte(key)))
		return nil
	})
	return value, err
}

// Delete godoc - Deletes a key/value pair from the database
func (db *BoltDB) Delete(key string) error {
	return db.Update(func(tx *bolt.Tx) error {

		// open the bucket
		b := tx.Bucket(db.bucketName)
		if b == nil {
			return errors.New("bucket not found")
		}

		// delete key-value pair from db
		return b.Delete([]byte(key))
	})
}

// List godoc - Lists all key/value pairs in the database
func (db *BoltDB) List(args ...interface{}) (data []KeyValuePair, err error) {

	err = db.View(func(tx *bolt.Tx) error {

		// open the bucket
		b := tx.Bucket(db.bucketName)
		if b == nil {
			return errors.New("bucket not found")
		}

		// get all keys
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			data = append(data, KeyValuePair{string(k), string(v)})
		}
		return nil
	})
	return
}

// AddBucket godoc - Adds a new bucket to the database
func (db *BoltDB) AddBucket(bucketName string) error {

	// add the bucket to the boltdb file
	return db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
}

// SetBucket godoc - Sets the current bucket
func (db *BoltDB) SetBucket(bucketName string) {
	db.bucketName = []byte(bucketName)
}

// SetBucket godoc - Sets the current bucket
func (db *BoltDB) GetDefBucket() string {
	return string(db.bucketName)
}

// ListBuckets godoc - Lists all buckets in the database
func (db *BoltDB) ListBuckets() (data []string, err error) {

	// open the db in view mode
	err = db.View(func(tx *bolt.Tx) error {

		// iterate on each buckets from the db
		err = tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			data = append(data, string(name))
			return err
		})
		return err
	})
	return
}

// DeleteBucket godoc - Deletes a bucket from the database
func (db *BoltDB) DeleteBucket(bucketName string) error {

	// delete the bucket from the boltdb file
	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
}
