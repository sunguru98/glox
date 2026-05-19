package main

import (
	"fmt"
	"os"
)

func RunFile(filePath string) {
}

func RunPrompt() {
}

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		RunFile(args[0])
	} else {
		RunPrompt()
	}
}
