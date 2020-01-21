package databasefunctions

import (
	ehfuncs "github.com/CodeZipline/project-0/errorhandlerfunctions"
	"github.com/dgraph-io/badger"
)

// DbDelete Searches through the LSM tree and places a delete marker during commit stage
func DbDelete(db *badger.DB, k string) string {
	// Start a writable transaction, delete
	txn := db.NewTransaction(true)
	//Implicityly called when Commit() is called or used to discard read
	// only transaction, either way safe to defer this function to the end
	defer txn.Discard()

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(k))
		return err
	})
	ehfuncs.Ehandler(err)
	return k
}
