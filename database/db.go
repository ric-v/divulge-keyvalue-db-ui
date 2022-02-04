package database

import "errors"

const (
	BOLT_DB = "boltdb"
	BUNT_DB = "buntdb"
)

type KeyValuePair struct {
	Key   string
	Value string
}

// DB interface for underlying database packages
type DB interface {
	Add(string, string, ...interface{}) error
	CloseDB()
	Get(key string) (string, error) // TODO: add list all keys
	Delete(string) error
	List(args ...interface{}) ([]KeyValuePair, error)
}

// BoltBuckets interface for underlying bolt DB buckets
type Buckets interface {
	Add()
	Get()
	List()
	Delete()
}

// DBType for identifying underlying database packages
type DBType string

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
