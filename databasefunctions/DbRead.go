package databasefunctions

import (
	ehfuncs "github.com/CodeZipline/project-0/errorhandlerfunctions"
	"github.com/dgraph-io/badger"
)

// DbRead is a key value search
func DbRead(db *badger.DB, k string) string {
	// Start a readable transaction
	txn := db.NewTransaction(false)
	//Implicityly called when Commit() is called or used to discard read
	// only transaction, either way safe to defer this function to the end
	defer txn.Discard()
	//k, _ = stdInRead("readMode")

	var v string

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		ehfuncs.Ehandler(err)

		err = item.Value(func(val []byte) error {
			// This func with val would only be called if item.Value encounters no error.
			//fmt.Printf("The key is %s, the value is: %s\n", k, val)

			v = string(val)
			return nil
		})
		ehfuncs.Ehandler(err)

		return nil
	})
	ehfuncs.Ehandler(err)
	return v
}

// DbFullReadOnKeys will create an iterator that will not fetch values and append keys to a slice that will be returned.
func DbFullReadOnKeys(db *badger.DB) []string {
	var retString []string
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		// PrefetchValue is set to false to only search on key values, fetching no values.
		opts.PrefetchValues = false
		// Iterator over tree for key
		it := txn.NewIterator(opts)
		defer it.Close()
		// Rewind moves to zeroth postion, valid checks for when iteration is over, next advacnes the iterator by one
		for it.Rewind(); it.Valid(); it.Next() {
			// Pointer to the key-value pair, expires when next is called
			item := it.Item()
			k := item.Key()
			// Must not print values since they are being fetched and can cause errors, since prefetch
			//  is set to false values will nto be copied over and lock still remains
			//fmt.Printf("key=%s\n", k)
			retString = append(retString, string(k))
		}
		return nil
	})
	ehfuncs.Ehandler(err)

	return retString
}
