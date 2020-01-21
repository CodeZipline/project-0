package databasefunctions

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
)

// DbQuit will exit properly
func DbQuit(db *badger.DB) {
	fmt.Println("Quiting Program")
	db.Close()
	os.Exit(0)
}
