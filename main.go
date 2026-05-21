package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sunguru98/glox/lib"
	"github.com/sunguru98/glox/scanner"
)

func Run(source string) {
	// Scanner
	scanner := scanner.InitScanner(source)
	// Tokens
	tokens := scanner.ScanTokens()
	// Print each token
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func RunFile(filePath string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	Run(string(bytes))

	if lib.HadError {
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
		lib.HadError = false
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading line: %v\n", err)
		os.Exit(-1)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		err := RunFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(-1)
		}

	} else {
		RunPrompt()
	}
}
