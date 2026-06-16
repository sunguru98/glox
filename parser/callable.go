package parser

import (
	"fmt"
	"time"
)

// ----------------------------------------------------------------------------------------
type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) any
}

// ----------------------------------------------------------------------------------------

type Clock struct{}

func (*Clock) Arity() int {
	return 0
}

func (*Clock) Call(interpreter *Interpreter, arguments []any) any {
	return float64(time.Now().UnixMilli()) / 1000.0
}

func (*Clock) String() string {
	return "<native fn>"
}

// ----------------------------------------------------------------------------------------

type Function struct {
	Declaration *FunctionSt
}

func (f *Function) Arity() int {
	return len(f.Declaration.Params)
}

func (f *Function) Call(interpreter *Interpreter, arguments []any) any {
	environment := InitEnvironment(interpreter.Globals)
	for paramIndex, param := range f.Declaration.Params {
		argument := arguments[paramIndex]
		environment.Define(param.Lexeme, argument)
	}

	interpreter.executeBlock(f.Declaration.Body, environment)
	return nil
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s >", f.Declaration.Name.Lexeme)
}

func CreateFunction(declaration *FunctionSt) *Function {
	return &Function{
		Declaration: declaration,
	}
}

// ----------------------------------------------------------------------------------------
