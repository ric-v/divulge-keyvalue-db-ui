package database

import "errors"

const (
	// BOLT_DB - BoltDB
	BOLT_DB = "boltdb"
	// BUNT_DB - BuntDB
	BUNT_DB = "buntdb"
)

type KeyValuePair struct {
	Key   string
	Value string
}

// DB interface for underlying database packages
type DB interface {
	Conn() interface{}
	Add(string, string, ...interface{}) error
	CloseDB()
	Get(key string) (string, error) // TODO: add list all keys
	Delete(string) error
	List(args ...interface{}) ([]KeyValuePair, error)
}

// DBType for identifying underlying database packages
type DBType string

// NewDB godoc - creates a new DB instance abstracting the underlying database package
func NewDB(fileName string, dbtype string) (DB, error) {

	// TODO: Add support for other database packages
	switch dbtype {
	case BOLT_DB:
		return openBolt(fileName)
	case BUNT_DB:
		return openBunt(fileName)
	default:
		return nil, errors.New("invalid database type")
	}
}
