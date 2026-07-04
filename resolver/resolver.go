package resolver

import (
	"github.com/sunguru98/glox/lib"
	p "github.com/sunguru98/glox/parser"
	s "github.com/sunguru98/glox/scanner"
)

// -------------------------------------------- TYPES ------------------------------------------------------

type FunctionType = int

const (
	// By default a code snippet will be outside a function
	None FunctionType = iota
	// Setting this if inside a function block
	Function
	// Setting this if inside a constructor function
	Initializer
	// Setting this if inside a class method
	Method
)

type ClassType = int

const (
	CNone ClassType = iota
	Class
	SubClass
)

// ------------------------------ RESOLVER STRUCT ----------------------------------------------------------

type Resolver struct {
	Interpreter     *p.Interpreter
	Scopes          []map[string]bool
	CurrentFunction FunctionType
	CurrentClass    ClassType
}

// ----------------------------- PRIMARY FUNCTIONS ---------------------------------------------------------

func InitResolver(interpreter *p.Interpreter) *Resolver {
	return &Resolver{
		Interpreter:     interpreter,
		Scopes:          make([]map[string]bool, 0),
		CurrentFunction: None,
		CurrentClass:    CNone,
	}
}

func (r *Resolver) ResolveStatements(statements []p.Statement) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

// ----------------------------- RESOLVER FUNCTIONS -------------------------------------------------------

func (r *Resolver) resolveStatement(statement p.Statement) {
	switch st := statement.(type) {
	case *p.ClassSt:
		// Since by default, the resolver assumes there is no class to begin with
		// We assign the class type to be a Class
		enclosingClass := r.CurrentClass
		r.CurrentClass = Class

		// A class behaves the same as function
		// A block exists within, hence we declare and define
		r.declare(st.Name)
		r.define(st.Name)

		stLexeme := st.Name.Lexeme
		superClassLexeme := st.SuperClass.Name.Lexeme
		// The base and derived classes should be different
		if st.SuperClass != nil && stLexeme == superClassLexeme {
			superClass := st.SuperClass.Name
			lib.Report(superClass.LineNumber, " at '"+superClass.Lexeme+"'", "A class cannot inherit from itself")
		}

		// We check if the derives from a base class.
		// If yes, we resolve it too
		if st.SuperClass != nil {
			r.CurrentClass = SubClass
			r.resolveExpression(st.SuperClass)
		}

		// Incase we inherit from the base class
		// We allocate the scope for the super keyword
		if st.SuperClass != nil {
			r.beginScope()
			scopesLen := len(r.Scopes)
			scopeTop := r.Scopes[scopesLen-1]
			scopeTop["super"] = true
		}

		// A scope is created for the class's body
		// And the 'this' keyword is mentioned as declared and defined
		// So that when instances are created, the scope would be perfectly aligned
		r.beginScope()
		scopesLen := len(r.Scopes)
		scopeTop := r.Scopes[scopesLen-1]
		scopeTop["this"] = true

		// Every class has it's own member methods for their instances
		// We resolve each of their methods
		for _, method := range st.Methods {
			declaration := Method
			if method.Name.Lexeme == "init" {
				declaration = Initializer
			}

			r.resolveFunction(method, declaration)
		}

		// Finally the scope is destroyed
		r.endScope()

		// Followed by the super/base class's scope (if present)
		if st.SuperClass != nil {
			r.endScope()
		}

		// And the default class is reverted
		r.CurrentClass = enclosingClass

	case *p.BlockSt:
		// For a block statement, we create a new scope
		r.beginScope()
		// Resolve each statement inside the block
		r.ResolveStatements(st.Statements)
		// And delete/free the allocated scope
		r.endScope()

	case *p.VariableSt:
		// A variable is defined in the scopes stack (setting to false)
		r.declare(st.Name)
		if st.Initializer != nil {
			// And the initializer value is resolved
			// To check for same variable initialization
			r.resolveExpression(st.Initializer)
		}
		// Followed by it's definition (setting to true)
		r.define(st.Name)

	case *p.FunctionSt:
		// The function statement (declaration) is declared and defined
		// In the scope where it's declared, so that recursion would be available
		// Since marking it false would error out when it's value is seen again in the body
		r.declare(st.Name)
		r.define(st.Name)
		// Then the function body is resolved, with the FunctionType set to Function
		r.resolveFunction(st, Function)

	case *p.ExpressionSt:
		// An expression statement is simply a normal expression
		// Hence the resolveExpression would take care of the same
		r.resolveExpression(st.Expr)

	case *p.IfSt:
		// The if condition is resolved
		r.resolveExpression(st.Condition)
		// The branch when the condition is truthy is resolved
		r.resolveStatement(st.ThenBranch)
		// As well as the else statement
		// Since the interpreter can go eitherway
		if st.ElseBranch != nil {
			// And resolving both branches helps with the same
			r.resolveStatement(st.ElseBranch)
		}

	case *p.PrintSt:
		// Just like normal expression statements
		// Print statements hand over the resolution to resolveExpression
		r.resolveExpression(st.Expr)

	case *p.ReturnSt:
		if r.CurrentFunction == None {
			// A return statement should always be inside a function
			// And not as an empty standalone statement
			// Hence we mark it as a compiler error
			keyword := st.Keyword
			lib.Report(keyword.LineNumber, " at '"+keyword.Lexeme+"'", "Can't return from top-level code")
		}

		// If the return statement has a value
		if st.Value != nil {
			// And the function being an initializer
			// Then a return function isn't possible
			if r.CurrentFunction == Initializer {
				lib.Report(st.Keyword.LineNumber, " at '"+st.Keyword.Lexeme+"'", "Can't return a value from an intializer")
			}

			// We resolve the returned value
			r.resolveExpression(st.Value)
		}

	case *p.WhileSt:
		// Similar to if condition, the condition here is also resolved
		r.resolveExpression(st.Condition)
		// Followed by the body of it
		r.resolveStatement(st.Body)
	}
}

func (r *Resolver) resolveExpression(expression p.Expression) {
	switch expr := expression.(type) {
	case *p.Variable:
		// We fetch the last registered variable in scope
		scopesLen := len(r.Scopes)
		scopeTop := r.Scopes[scopesLen-1]

		// And see if the initializer expression's name matches with the variable name
		// In that case, it's a reinitialization error, and we mark it compile-time
		token := expr.Name
		scopeVal, ok := scopeTop[token.Lexeme]
		if scopesLen != 0 && (ok && scopeVal == false) {
			lib.Report(token.LineNumber, " at '"+token.Lexeme+"'", "Can't read local variable in own initializer")
		}

		// If not, we resolve the variable through finding the scope depth
		// Scope depth is discussed inside this function
		r.resolveLocal(expr, token)

	case *p.This:
		// The this keyword must be used only inside a class based setting
		token := expr.Keyword
		if r.CurrentClass == CNone {
			// Hence we early return if that is not the case
			lib.Report(token.LineNumber, " at '"+token.Lexeme+"'", "Can't use 'this' outside of a class")
			return
		}

		// The this keyword would act just like a variable,
		// as it's mentioning the current instance of that class
		r.resolveLocal(expr, token)

	case *p.Call:
		// Function calls are basically similar to declaration
		// Instead of defining/declaring the function, we just resolve it
		r.resolveExpression(expr.Callee)
		// And followed by the resolution of each argument inside it
		for _, argument := range expr.Arguments {
			r.resolveExpression(argument)
		}

	case *p.Get:
		// A get expression could be a member field/method
		// Hence we resolve the underlying expression (the field/method)
		r.resolveExpression(expr.Object)

	case *p.Set:
		// A set expression sets the value w.r.t the object (instance)
		// Hence we resolve both
		r.resolveExpression(expr.Value)
		r.resolveExpression(expr.Object)

	case *p.Super:
		// The super keyword cannot be used outside a class
		// Nor in a class that has not inherited a base class
		if r.CurrentClass == CNone {
			lib.Report(expr.Keyword.LineNumber, " at '"+expr.Keyword.Lexeme+"'", "Can't use 'super' outside of a class")
		} else if r.CurrentClass != SubClass {
			lib.Report(expr.Keyword.LineNumber, " at '"+expr.Keyword.Lexeme+"'", "Can't use 'super' in a class with no superclass")
		}

		// Resolving it just like a variable
		r.resolveLocal(expr, expr.Keyword)

	case *p.Assignment:
		// The assignment operation simply means
		// To resolve the value it's being assigned to
		r.resolveExpression(expr.Value)
		// Followed by the variable the assigned value holds
		r.resolveLocal(expr, expr.Name)

	case *p.Binary:
		// Similar to Unary, instead of one operand
		// Binary has 2 operands
		r.resolveExpression(expr.Left)
		r.resolveExpression(expr.Right)

	case *p.Logical:
		r.resolveExpression(expr.Left)
		r.resolveExpression(expr.Right)

	case *p.Unary:
		// A unary operand only has an expression to the right
		// Hence we resolve just that
		r.resolveExpression(expr.Right)

	case *p.Grouping:
		// The expression whatever is in the parentheses is resolved
		r.resolveExpression(expr.Expression)

	case *p.Literal:
		// Since literal values have no variables/subexpressions
		// We just return
		return
	}
}

func (r *Resolver) resolveLocal(expr p.Expression, name s.Token) {
	scopesLen := len(r.Scopes) - 1
	for index := scopesLen; index >= 0; index -= 1 {
		// We traverse the stack from top (last element in array)
		// to bottom (first element in array)
		// And see if the scope has the token
		scope := r.Scopes[index]
		_, ok := scope[name.Lexeme]

		// If yes, we calculate the depth
		// Depth - Distance between the current scope the resolution happens
		// Versus the actual scope the variable is being found
		if ok {
			depth := scopesLen - index
			// And the interpreter is being told to record the same
			r.Interpreter.Resolve(expr, depth)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function *p.FunctionSt, fType FunctionType) {
	// Assigning a temp variable with the currentFunction value, to change after
	enclosingFunction := r.CurrentFunction
	// Temporarily assigning the function type being Function, as the resolution is inside the function
	r.CurrentFunction = fType

	// The function scope is being attached for the body
	r.beginScope()
	// And each parameter is being declared and marked as defined
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}

	// The function statements inside are all resolved
	r.ResolveStatements(function.Body)
	// Followed by the detachment of the same scope.
	r.endScope()

	// Switching back to whatever r.CurrentFunction was before the switch
	r.CurrentFunction = enclosingFunction
}

// --------------------------------------- UTILS -----------------------------------------------------------

func (r *Resolver) beginScope() {
	// A scope is created specifically for a block whenever it's starting up
	// Each scope will have the token/variable name followed by a boolean to indicate
	// Whether it has been defined or declared
	// A stack like structure is maintained
	r.Scopes = append(r.Scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	// Ending a scope simply means removing the scope data out of the list
	// Popping a stack in technical means
	scopesLen := len(r.Scopes)
	r.Scopes = r.Scopes[:scopesLen-1]
}

func (r *Resolver) declare(name s.Token) {
	// Check if there exists a scope.
	// If the scope stack is empty, then there is no variable to attach with
	scopesLen := len(r.Scopes)
	if scopesLen == 0 {
		return
	}

	// Else, we peek (Fetch the top element of the stack)
	scopeTop := r.Scopes[scopesLen-1]
	// And check if there exsits already a variable declared
	_, ok := scopeTop[name.Lexeme]
	if ok {
		// Since same variable reinitialization is marked as a compiler error
		lib.Report(name.LineNumber, " at '"+name.Lexeme+"'", "Already a variable with this name in this scope")
	}

	// Since we haven't checked the initializer value yet, we mark it as not-ready
	// A.k.a this variable is not ready for resolving scope, hence marked as false
	scopeTop[name.Lexeme] = false
}

func (r *Resolver) define(name s.Token) {
	// Empty stack check
	scopesLen := len(r.Scopes)
	if scopesLen == 0 {
		return
	}

	// Defining is called after declaration
	// And initializer value being resolved
	// Hence we mark it as true
	// As in it's available for scope resolution
	scopeTop := r.Scopes[scopesLen-1]
	scopeTop[name.Lexeme] = true

}

// --------------------------------------------------------------------------------------------------------
