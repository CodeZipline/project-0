package main

import (
	c "github.com/CodeZipline/project-0/configurations"
	dbfuncs "github.com/CodeZipline/project-0/databasefunctions"
	ehfuncs "github.com/CodeZipline/project-0/errorhandlerfunctions"
)

func main() {
	db, err := dbfuncs.OpenDatabase(c.DbArchitecture.DD)
	ehfuncs.Ehandler(err)
	//It's largely used to clean up connections, which is what you're doing with this code and http response bodies.
	// db.Close() closes the underlying db connection and you're telling the code to execute it when the function exits."
	defer db.Close()

	dbfuncs.DbMethods()
	for {
		dbfuncs.DbMenu(db)
	}

}
