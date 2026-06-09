package parser

import (
	"fmt"

	s "github.com/sunguru98/glox/scanner"
)

// ------------------------------------------- ENVIRONMENT ----------------------------------------------
// The place where variable's values are being stored
type Environment struct {
	Map map[string]any
}

func (e *Environment) Define(name string, value any) {
	e.Map[name] = value
}

func (e *Environment) Get(name s.Token) (any, error) {
	lexeme := name.Lexeme
	value, ok := e.Map[lexeme]

	if !ok {
		return nil, fmt.Errorf("Undefined variable %s.", lexeme)
	}

	return value, nil
}

// -----------------------------------------------------------------------------------------------------

func InitEnvironment() *Environment {
	return &Environment{
		Map: make(map[string]any),
	}
}
