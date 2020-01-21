package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	c "github.com/CodeZipline/project-0/configurations"
	dbfuncs "github.com/CodeZipline/project-0/databasefunctions"
	ehfuncs "github.com/CodeZipline/project-0/errorhandlerfunctions"
)

func main() {
	db, err := dbfuncs.OpenDatabase(c.DbArchitecture.DD)
	ehfuncs.Ehandler(err)
	defer db.Close()

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("key")
		value := dbfuncs.DbRead(db, key)
		//output, err := json.Marshal(products.inventory)
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
			fmt.Fprintln(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, key, value)
	})

	fmt.Println("Listening on ports 8080 (http) and 8081 (https)...")

	errorChan := make(chan error, 5)
	go func() {
		errorChan <- http.ListenAndServe(":8080", nil)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	for {
		select {
		case err := <-errorChan:
			if err != nil {
				log.Fatalln(err)
			}

		case sig := <-signalChan:
			fmt.Println("\nShutting down due to", sig)
			os.Exit(0)
		}
	}
}
