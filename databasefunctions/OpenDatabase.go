package databasefunctions

import (
	"github.com/dgraph-io/badger"
)

// OpenDatabase desciprtor
func OpenDatabase(dir string) (*badger.DB, error) {
	//Define Storage pathway for data.
	opts := badger.DefaultOptions(dir)

	// Open the Badger database located in the project-0/data directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(opts)

	return db, err
}
