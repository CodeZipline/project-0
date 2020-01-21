package databasefunctions

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	c "github.com/CodeZipline/project-0/configurations"
	"github.com/dgraph-io/badger"
)

// rwMutexStruct holds mutex for synchronizing read and writes to the database
type rwMutexStruct struct {
	//counters map[string]*int64
	mutex sync.RWMutex
}

// DbMethods is a collection of print statements for user to know what commands can be execute with this program
func DbMethods() {
	fmt.Println("Press r to read from the database.")
	fmt.Println("Press r_keys to read all key values in the database")
	fmt.Println("Press w to write to the database.")
	fmt.Println("Press w_TTL to write to the database, with TTL flag value.")
	fmt.Println("Press gc to perform a garbage collection method.")
	fmt.Println("Press d to delete a key from the database.")
	fmt.Println("Press q to quit this program.")
}

// DbMenu is a switch statement for calling other database related functions
func DbMenu(db *badger.DB) {
	s := stdInCommand("commandMode")

	if reflect.TypeOf(s).Kind() != reflect.String {
		log.Fatalln("invalid input")
		return
	}

	var rwM rwMutexStruct

	switch s {
	case "r":
		rwM.mutex.RLock()
		defer rwM.mutex.RUnlock()
		k, _ := stdInRead("readMode")
		fmt.Println(DbRead(db, k))
	case "r_keys":
		rwM.mutex.RLock()
		defer rwM.mutex.RUnlock()
		fmt.Println(DbFullReadOnKeys(db))
	case "w":
		rwM.mutex.Lock()
		defer rwM.mutex.Unlock()
		k, v := stdInRead("writeMode")
		fmt.Println(DbWrite(db, k, v))
	case "w_TTL":
		rwM.mutex.Lock()
		defer rwM.mutex.Unlock()
		k, v := stdInRead("writeMode")
		fmt.Println(DbWriteTTL(db, k, v, c.TTL))
	case "gc":
		dbGCContinous(db, c.GCINTERVAL, c.GCDURATION, c.GCDRS)
	case "d":
		k, _ := stdInRead("deleteMode")
		fmt.Println(DbDelete(db, k))
	case "q":
		DbQuit(db)
	default:
		log.Fatal("Exiting with default case")
	}
}
