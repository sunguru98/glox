package lib

import (
	"bufio"
	"fmt"
	"os"
)

var hadError = false

func Error(lineNumber int, message string) {
	Report(lineNumber, "", message)
}

func Report(lineNumber int, place, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", lineNumber, place, message)
	hadError = true
}

func Run(input string) {
	// Scanner
	// Tokens
	// Print each token
}

func RunFile(filePath string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	Run(string(bytes))

	if hadError {
		os.Exit(65)
	}

	return nil
}

func RunPrompt() {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Print("> ")
		text := scanner.Text()
		if text == "" {
			break
		}

		Run(text)
		hadError = false
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading line: %v\n", err)
		os.Exit(-1)
	}
}
