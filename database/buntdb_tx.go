package database

import (
	"github.com/tidwall/buntdb"
)

// BuntDB - BuntDB struct
type BuntDB struct {
	*buntdb.DB
}

// Conn godoc - Returns the underlying database connection
func (db *BuntDB) Conn() interface{} {
	return db
}

// openBunt godoc - Creates a new BuntDB instance
func openBunt(fileName string) (db *BuntDB, err error) {

	db = &BuntDB{}
	// open new bolt file
	db.DB, err = buntdb.Open(fileName)
	return
}

// CloseDB godoc - Closes the database
func (db *BuntDB) CloseDB() {
	db.Close()
}

// Add godoc - Adds a new key/value pair to the database
func (db *BuntDB) Add(key string, value string, args ...interface{}) error {

	// add the key-value pair to the buntDB file
	err := db.Update(func(tx *buntdb.Tx) error {

		// insert new key-value pair to db
		_, _, err := tx.Set(key, value, nil)
		return err
	})
	return err
}

// Get godoc - Gets a value from the database
func (db *BuntDB) Get(key string) (string, error) {

	var value string
	err := db.View(func(tx *buntdb.Tx) error {

		// get value from db
		var err error
		value, err = tx.Get(key)
		return err
	})
	return value, err
}

// Delete godoc - Deletes a key/value pair from the database
func (db *BuntDB) Delete(key string) error {
	return db.Update(func(tx *buntdb.Tx) error {

		// delete key-value pair from db
		_, err := tx.Delete(key)
		return err
	})
}

// List godoc - Lists all keys in the database
func (db *BuntDB) List(args ...interface{}) (data []KeyValuePair, err error) {

	err = db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(key, value string) bool {
			data = append(data, KeyValuePair{key, value})
			return true
		})
		return nil
	})
	return
}
