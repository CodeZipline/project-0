package databasefunctions

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
)

// DbGCContinous will run for a [dur] amount of minutes with garbage collection running in intervals of [inter] minutes
func DbGCContinous(db *badger.DB, inter int, dur int, dRS float64) {

	ticker := time.NewTicker(time.Duration(inter) * time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(time.Duration(dur) * time.Minute)
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
		again:
			fmt.Printf("GC with value %F discardRatioSpace, ", dRS)
			fmt.Println("at time: ", t)
			err := db.RunValueLogGC(dRS)
			if err == nil {
				goto again
			} else if err == badger.ErrNoRewrite {
				fmt.Println("No files found to be rewritten in the value logs.")
			}
		}
	}
}
