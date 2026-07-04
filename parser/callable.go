package parser

import (
	"errors"
	"fmt"
	"time"

	s "github.com/sunguru98/glox/scanner"
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
	IsInitializer bool
	Declaration   *FunctionSt
	ClosureEnv    *Environment // The environment outside the function declaration
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
			// If the function is a constructor function, we fetch the value of 'this'
			if f.IsInitializer {
				return f.ClosureEnv.GetAt(0, "this"), nil
			}

			// If yes, fetch the value and return
			return returnValue.Value, nil
		}

		// Else, it's a regular non-return value error
		return nil, err
	}

	// If the function calling is an initializer function
	// The closure environment surrounding it is the immediate one with 'this' keyword
	if f.IsInitializer {
		return f.ClosureEnv.GetAt(0, "this"), nil
	}

	// Else, by default all functions are of 'void' type
	// Meaning they return nil unless specified
	return nil, nil
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s >", f.Declaration.Name.Lexeme)
}

func (f *Function) Bind(instance *Instance) *Function {
	// A new environment is created keeping the closure as the parent
	env := InitEnvironment(f.ClosureEnv)
	// The environment then assigns the this keyword to be locked within that scope
	env.Define("this", instance)

	// And the resultant function/method is created through that environment
	return CreateFunction(f.Declaration, env, f.IsInitializer)
}

func CreateFunction(declaration *FunctionSt, closureEnv *Environment, initializer bool) *Function {
	return &Function{
		ClosureEnv:    closureEnv,
		Declaration:   declaration,
		IsInitializer: initializer,
	}
}

// ----------------------------------------------------------------------------------------

// A class type in the programming language
type Class struct {
	Name       string
	Methods    map[string]*Function
	SuperClass *Class
}

func (c *Class) Arity() int {
	// Either the constructor function doesn't exist
	// which means 0 arguments obviously
	initializer := c.FindMethod("init")
	if initializer == nil {
		return 0
	}

	// Or we return the number of arguments of the constructor function
	return initializer.Arity()
}

func (c *Class) Call(interpreter *Interpreter, arguments []any) (any, error) {
	// When the class's constructor is called, an instance is created
	instance := CreateInstance(c)
	initializer := c.FindMethod("init")

	// If the initializer function/constructor exists, we call it
	// Via the current instance being bound with the constructor arguments
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}

	return instance, nil
}

func (c *Class) String() string {
	return c.Name
}

func (c *Class) FindMethod(name string) *Function {
	// If the method exists in subclass first, we invoke that
	if method, ok := c.Methods[name]; ok {
		return method
	}

	// Or if the method exists in the base class, we invoke that
	if c.SuperClass != nil {
		return c.SuperClass.FindMethod(name)
	}

	return nil
}

func CreateClass(name string, methods map[string]*Function, superClass *Class) *Class {
	return &Class{
		Name:       name,
		Methods:    methods,
		SuperClass: superClass,
	}
}

// ------------------------------------------------------------------------------------------

// An instance type in the programming language (created through Class.Call())
type Instance struct {
	class  *Class
	fields map[string]any
}

func CreateInstance(class *Class) *Instance {
	return &Instance{
		class:  class,
		fields: make(map[string]any),
	}
}

func (i *Instance) String() string {
	return i.class.Name + " instance"
}

func (i *Instance) Get(name s.Token) (any, error) {
	// If the name is a member field, we return that
	field, ok := i.fields[name.Lexeme]
	if ok {
		return field, nil
	}

	// Or if it's a member method, we return that instead
	method := i.class.FindMethod(name.Lexeme)
	if method != nil {
		// We not just return the method
		// But also bind the 'this' keyword for that method
		// So that, each method is exclusive for their own objects/instances
		return method.Bind(i), nil
	}

	// Else, there exists no such name
	return nil, fmt.Errorf("Undefined property %s", name.Lexeme)
}

func (i *Instance) Set(name s.Token, value any) {
	i.fields[name.Lexeme] = value
}

// ------------------------------------------------------------------------------------------
