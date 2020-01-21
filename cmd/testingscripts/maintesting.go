package main

import (
	"fmt"
	"github.com/CodeZipline/project-0/configurations"
	"github.com/CodeZipline/project-0/databasefunctions"
	"github.com/CodeZipline/project-0/errorhandlerfunctions"
	"github.com/dgraph-io/badger"
)

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

func main() {
	db, err := databasefunctions.OpenDatabase(configurations.DbArchitecture.DD)
	errorhandlerfunctions.Ehandler(err)
	defer db.Close()
	testDbWrite(db)
	testDbRead(db)
	testDbFullReadOnKeys(db)
}
