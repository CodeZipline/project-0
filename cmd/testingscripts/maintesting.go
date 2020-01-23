package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/CodeZipline/project-0/configurations"
	"github.com/CodeZipline/project-0/databasefunctions"
	"github.com/CodeZipline/project-0/errorhandlerfunctions"
	"github.com/dgraph-io/badger"
)

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// testDbWrite will determine if the input data will be written and returned once DbWrite has executed
func testDbWrite(db *badger.DB) {
	var testingk, testingv string

	testingk, testingv = databasefunctions.DbWrite(db, "tK", "tV")

	if testingk != "tK" {
		fmt.Printf("input key: %s != output key: %s. \n", "tK", testingk)
	} else if testingv != "tV" {
		fmt.Printf("input value: %s != output value: %s. \n", "tV", testingv)
	} else {
		fmt.Println("Value written correctly.")
	}
}

// testDbRead will determine if the input data from the last function is read properly
func testDbRead(db *badger.DB) {
	var testingv string

	testingv = databasefunctions.DbRead(db, "tK")

	if testingv != "tV" {
		fmt.Printf("input value: %s != output value: %s. \n", "tV", testingv)
	} else {
		fmt.Println("Value read correctly.")
	}
}

// testDbFullReadOnKeys will check that the previous read key is found in the slice of all keys in the database
func testDbFullReadOnKeys(db *badger.DB) {
	var retSlice []string
	var testingk string
	testingk = "tK"
	retSlice = databasefunctions.DbFullReadOnKeys(db)
	found := func() bool {
		for _, item := range retSlice {
			if item == testingk {
				return true
			}
		}
		return false
	}()
	if !found {
		fmt.Printf("The input key %s, is not found in the database.\n", testingk)
	} else {
		fmt.Println("Found key correctly.")
	}
}

// testDbDelete will read a known key, delete it and check again if it exist in the database
func testDbDelete(db *badger.DB) {
	// Check for existing key
	var retSlice []string
	var testingk string
	testingk = "tK"
	retSlice = databasefunctions.DbFullReadOnKeys(db)
	testDbFullReadOnKeys(db)

	// Perform the delete operation
	fmt.Println("Deleting the key.")
	_ = databasefunctions.DbDelete(db, testingk)
	newretSlice := databasefunctions.DbFullReadOnKeys(db)
	same := Equal(retSlice, newretSlice)
	if same {
		fmt.Printf("The input key %s has not been deleted from the database.\n", testingk)
	} else {
		fmt.Println("Deleted key correctly.")
	}
	// Cleaning up
	fmt.Println("Garbage Collection")
	databasefunctions.DbGCContinous(db, 10, 1, 0.5)

}

type kvPairs struct {
	key string
	val string
}

func main() {

	Logfile, err := os.OpenFile("logs/log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file %v \n", err)
		panic(err)
	}
	log.SetOutput(Logfile)

	Answers := []string{
		"Coding Answer : Found on stackoverflow.com",
		"Life Answer : 42",
		"Answer To Any Problem : Our Lord and Savior Jesus Christ",
		"9 + 10 : 21",
		"What Ended In 1986? : 1985",
		"Is This The Krusty Krab? : NO, This is Patrick!",
		"End All Be All Coding Language : Go",
	}

	rand.Seed(time.Now().UTC().UnixNano())

	db, err := databasefunctions.OpenDatabase(configurations.DbArchitecture.DD)
	errorhandlerfunctions.Ehandler(err)
	defer db.Close()

	kvMessagesChan := make(chan kvPairs)

	go func() {
		for {
			select {
			case resp := <-kvMessagesChan:
				log.Printf("Reader finished. The key: %s, The Value: %s. \n", resp.key, databasefunctions.DbRead(db, resp.key))
			}
		}
	}()

	for w := 0; w < 3; w++ {
		go func() {
			for {
				write := kvPairs{
					key: Answers[rand.Intn(len(Answers))],
					val: "success"}
				dbKey, dbValue := databasefunctions.DbWrite(db, write.key, write.val)
				log.Printf("Writer finished. Key: %s, Value: %s.\n", dbKey, dbValue)
				kvMessagesChan <- write

				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)

	// Test functions to check for proper execution of basic functions
	testDbWrite(db)
	testDbRead(db)
	testDbFullReadOnKeys(db)
	fmt.Println("Checking delete function:")
	testDbDelete(db)

	os.Exit(0)

}
