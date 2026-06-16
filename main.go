package main

import (
	"bufio"
	"fmt"
	"os"

	l "github.com/sunguru98/glox/lib"
	p "github.com/sunguru98/glox/parser"
	s "github.com/sunguru98/glox/scanner"
)

var interpreter = p.InitInterpreter()

func Run(source string) {
	// Scanner
	scanner := s.InitScanner(source)
	// Tokens
	tokens := scanner.ScanTokens()

	// Parser
	parser := p.InitParser(tokens)
	// Parse Expression
	statements := parser.Parse()

	// Return if parsing reported an error
	if l.HadError {
		return
	}

	// Else evaluate the expression
	interpreter.Interpret(statements)
}

func RunFile(filePath string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	Run(string(bytes))

	if l.HadError {
		os.Exit(65)
	}

	if l.HadRuntimeError {
		os.Exit(70)
	}

	return nil
}

func RunPrompt() {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		Run(text)
		l.HadError = false
		fmt.Print("> ")
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
