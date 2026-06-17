package parser

import (
	"fmt"

	s "github.com/sunguru98/glox/scanner"
)

// ------------------------------------------- ENVIRONMENT ----------------------------------------------
// The place where variable's values are being stored
type Environment struct {
	// Enclosing environments are the ones
	// That are outside the current scope
	Enclosing *Environment
	Map       map[string]any
}

func (e *Environment) Define(name string, value any) {
	e.Map[name] = value
}

func (e *Environment) Get(name s.Token) (any, error) {
	lexeme := name.Lexeme
	value, ok := e.Map[lexeme]

	if ok {
		return value, nil
	}

	// If there exists an outer scope/environment pick that
	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	return nil, fmt.Errorf("Undefined variable %s.", lexeme)

}

func (e *Environment) Assign(name s.Token, value any) error {
	lexeme := name.Lexeme
	_, ok := e.Map[lexeme]

	if ok {
		e.Map[lexeme] = value
		return nil
	}

	// If there exists an outer scope/environment pick that
	if e.Enclosing != nil {
		err := e.Enclosing.Assign(name, value)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Undefined variable %s.", lexeme)
}

// -----------------------------------------------------------------------------------------------------

func InitEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		Map:       make(map[string]any),
	}
}
