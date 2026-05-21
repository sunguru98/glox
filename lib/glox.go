package lib

import (
	"fmt"
	"os"
)

var HadError = false

func Error(lineNumber int, message string) {
	Report(lineNumber, "", message)
}

func Report(lineNumber int, place, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", lineNumber, place, message)
	HadError = true
}
