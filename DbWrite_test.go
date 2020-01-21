package test

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"testing"

	dbfuncs "github.com/CodeZipline/project-0/databasefunctions"
	ehfuncs "github.com/CodeZipline/project-0/errorhandlerfunctions"
	s "github.com/CodeZipline/project-0/structs"
)

// TestDbWrite checks that when a write execution returns the key and value should be the same as the write function parameter.
func TestDbWrite(t *testing.T) {
	// TTL is a flag that is set to true, with a default of 24 hours for a Time To Live(TTL).
	var TTL int

	// GCDURATION is a flag that is set to determine how long the database Garbage Collection function should run for.
	var GCDURATION int

	// GCINTERVAL is a flag that is set to determine frequent a single call to GC function is made within the duration of GCDURATION
	var GCINTERVAL int

	// GCDRS is a flag that is set to filter file space usage with its lifetime value log write amplification in mind.
	var GCDRS float64

	// CONFIGFILE will be the location of the initial database set up parameters.
	const CONFIGFILE string = "config.json"

	// DbArchitecture ...
	var DbArchitecture s.DatabaseArchitecture

	flag.IntVar(&TTL, "ttl", 24, "Flag to set for Time To Live writes.")
	flag.IntVar(&GCDURATION, "gcd", 1, "Flag to set for Duration of Garbage Collection function calls. (Minutes)")
	flag.IntVar(&GCINTERVAL, "gci", 10, "Flag to set for Intervals of Garbage Collection function calls. (Seconds)")
	flag.Float64Var(&GCDRS, "gcdrs", 0.5, "Flag to set for the discard ratio space of a file. (1-0)")

	flag.Parse()

	//Uses the ReadFile, open the file and read all content and save results and returns as bytes
	configData, err := ioutil.ReadFile(CONFIGFILE)
	ehfuncs.Ehandler(err)
	//Reads the byte array and converts it into struct
	err = json.Unmarshal(configData, &DbArchitecture)
	ehfuncs.Ehandler(err)

	db, err := dbfuncs.OpenDatabase(DbArchitecture.DD)
	ehfuncs.Ehandler(err)
	//It's largely used to clean up connections, which is what you're doing with this code and http response bodies.
	// db.Close() closes the underlying db connection and you're telling the code to execute it when the function exits."
	defer db.Close()

	var tk, tv string

	tk, tv = dbfuncs.DbWrite(db, "testingKey", "testingValue")
	// Error if the write did not set the value and key to be the inputted value.
	if tk != "testingKey" && tv != "testingValue" {
		t.Errorf("Not return inputed keys and values")
	}
}
