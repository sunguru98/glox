package main

import (
	"fmt"
	"os"

	"github.com/sunguru98/glox/lib"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		err := lib.RunFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(-1)
		}

	} else {
		lib.RunPrompt()
	}
}
