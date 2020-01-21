package errorhandlerfunctions

import (
	"fmt"
	"log"
	"os"
)

// init ...
func init() {
	Logfile, err := os.OpenFile("logs/log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file %v", err)
		panic(err)
	}
	log.SetOutput(Logfile)
}

// Ehandler This function is for error handling
func Ehandler(err error) {
	if err != nil {
		log.Fatalf("Error message: %v", err)
	}
}
