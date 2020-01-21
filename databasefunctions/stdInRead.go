package databasefunctions

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	ehfuncs "github.com/CodeZipline/project-0/errorhandlerfunctions"
)

// -ReadHelper is a helperfunction to perform the read of stdio and make sure the result is not empty
func readHelper() string {
	s := ""
	for s == "" {
		reader := bufio.NewReader(os.Stdin)
		readLine, err := reader.ReadString('\n')
		ehfuncs.Ehandler(err)
		s = strings.TrimSuffix(readLine, "\n")
		if s == "" {
			fmt.Println("Enter a non-empty string: ")
		}
	}
	return s
}

// -StdInCommand switch statement to capture the user response to indicate which database function to execute
func stdInCommand(mode string) string {
	var response string
	switch mode {
	case "commandMode":
		fmt.Print("Enter Database Command Key: ")
		response = readHelper()
		return response
	default:
		fmt.Print("NO MODE SELECTED")
		return response
	}
}

// -StdInRead switch statement to perform the indicated database function
func stdInRead(mode string) (string, string) {
	var k, v string
	switch mode {
	case "readMode":
		fmt.Print("Enter The Key to be read: ")
		k = readHelper()
		return k, v
	case "writeMode":
		fmt.Print("Enter The Key to be written: ")
		k = readHelper()
		fmt.Print("Enter The Value to be written: ")
		v = readHelper()
		return k, v
	case "deleteMode":
		fmt.Print("Enter The Key to be deleted: ")
		k = readHelper()
		return k, v

	default:
		fmt.Print("NO MODE SELECTED")
		return k, v
	}

}
