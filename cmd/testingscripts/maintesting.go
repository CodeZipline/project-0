package main

import (
	"fmt"
	"github.com/CodeZipline/project-0/configurations"
	"github.com/CodeZipline/project-0/databasefunctions"
	"github.com/CodeZipline/project-0/errorhandlerfunctions"
	"github.com/dgraph-io/badger"
	"sync"
	"time"
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

type readOp struct {
    key  string
    resp chan string
}
type writeOp struct {
    key  string
    val  string
    resp chan bool
}

func main() {
	db, err := databasefunctions.OpenDatabase(configurations.DbArchitecture.DD)
	errorhandlerfunctions.Ehandler(err)
	defer db.Close()

	// Test functions to check for proper execution of basic functions
	testDbWrite(db)
	testDbWrite(db)
	testDbRead(db)
	testDbFullReadOnKeys(db)
	testDbDelete(db)


    var readOpCounter uint64
    var writeOpCounter uint64

    reads := make(chan readOp)
    writes := make(chan writeOp)

    go func() {
        for {
            select {
            case read := <-reads:
                read.resp <- databasefunctions.DbRead(db, read.key)
            case write := <-writes:
                databasefunctions.DbWrite(db, write.key, write.val)
                write.resp <- true
            }
        }
    }()

    for r := 0; r < 100; r++ {
        go func() {
            for {
                read := readOp{
                    key:  rand.Intn(5),
                    resp: make(chan int)}
                reads <- read
                <-read.resp
                atomic.AddUint64(&readOps, 1)
                time.Sleep(time.Millisecond)
            }
        }()
    }

    for w := 0; w < 10; w++ {
        go func() {
            for {
                write := writeOp{
                    key:  rand.Intn(5),
                    val:  rand.Intn(100),
                    resp: make(chan bool)}
                writes <- write
                <-write.resp
                atomic.AddUint64(&writeOps, 1)
                time.Sleep(time.Millisecond)
            }
        }()
    }

    time.Sleep(time.Second)

    readOpsFinal := atomic.LoadUint64(&readOps)
    fmt.Println("readOps:", readOpsFinal)
    writeOpsFinal := atomic.LoadUint64(&writeOps)
    fmt.Println("writeOps:", writeOpsFinal)
}

}
