package parser

import (
	"errors"
	"fmt"
	"time"
)

// ----------------------------------------------------------------------------------------
// Custom error as Return value
type ReturnValue struct {
	Value any
}

func (r *ReturnValue) Error() string {
	return fmt.Sprintf("Return value: %v", r.Value)
}

// ----------------------------------------------------------------------------------------
// All functions in the programming language implement this interface

type Callable interface {
	Arity() int                                                  // Represents number of parameters
	Call(interpreter *Interpreter, arguments []any) (any, error) // The function body invocation logic
}

// ----------------------------------------------------------------------------------------

type Clock struct{}

func (*Clock) Arity() int {
	return 0
}

func (*Clock) Call(interpreter *Interpreter, arguments []any) (any, error) {
	// Default Native function to fetch the time taken
	return float64(time.Now().UnixMilli()) / 1000.0, nil
}

func (*Clock) String() string {
	return "<native fn>"
}

// ----------------------------------------------------------------------------------------

// A function type in the programming language
type Function struct {
	Declaration *FunctionSt
	ClosureEnv  *Environment // The environment outside the function declaration
}

func (f *Function) Arity() int {
	// We return the function declaration's parameter count mentioned
	return len(f.Declaration.Params)
}

func (f *Function) Call(interpreter *Interpreter, arguments []any) (any, error) {
	// We create a sub-environment with the attached closure environment being the parent
	// You could think like a "stack frame" created with a function call
	environment := InitEnvironment(f.ClosureEnv)

	// For whatever argument is passed in the function call,
	// We map with the function declaration parameter
	// func(a, b) -> func(1, 2) -> a = 1, b = 2
	for paramIndex, param := range f.Declaration.Params {
		argument := arguments[paramIndex]
		environment.Define(param.Lexeme, argument)
	}

	// With the values that's mapped in this sub environment
	// The function block (in the declaration) is run statement by statement
	err := interpreter.executeBlock(f.Declaration.Body, environment)
	if err != nil {
		// Check if the error is of type ReturnValue
		if returnValue, ok := errors.AsType[*ReturnValue](err); ok {
			// If yes, fetch the value and return
			return returnValue.Value, nil
		}

		// Else, it's a regular non-return value error
		return nil, err
	}

	// By default all functions are of 'void' type
	// Meaning they return nil unless specified
	return nil, nil
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s >", f.Declaration.Name.Lexeme)
}

func CreateFunction(declaration *FunctionSt, closureEnv *Environment) *Function {
	return &Function{
		ClosureEnv:  closureEnv,
		Declaration: declaration,
	}
}

// ----------------------------------------------------------------------------------------
