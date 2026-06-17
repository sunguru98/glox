package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sunguru98/glox/lib"
	s "github.com/sunguru98/glox/scanner"
)

type Interpreter struct {
	Globals *Environment
	Env     *Environment
}

// ------------------------PRIMARY FUNCTIONS -------------------------

func InitInterpreter() *Interpreter {
	globals := InitEnvironment(nil)
	globals.Define("clock", &Clock{})

	return &Interpreter{
		Globals: globals,
		Env:     globals,
	}
}

func (i *Interpreter) Interpret(statements []Statement) {
	// We evaluate the expression, and check for errors
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			lib.RuntimeError(err)
		}
	}
}

// Execute handles all sorts of
// 1. Statements (Print, Expression, If, Block)
// 2. Variables

func (i *Interpreter) execute(st Statement) error {
	switch statement := st.(type) {
	case *PrintSt:
		// We evaluate the expression after print keyword
		value, err := i.evaluate(statement.Expr)
		if err != nil {
			return err
		}

		// And print the result
		fmt.Println(i.stringify(value))

	case *FunctionSt:
		// The function declaration statement is bunched as a function
		function := CreateFunction(statement)
		// And mapped with the interpreter with function name - node
		i.Env.Define(statement.Name.Lexeme, function)

	case *ExpressionSt:
		// Here since it's just an expression
		// No requirement to do anything other than evaluate
		_, err := i.evaluate(statement.Expr)
		if err != nil {
			return err
		}

	case *IfSt:
		// We check if the "if" conditional expression is truthy
		condition, err := i.evaluate(statement.Condition)
		if err != nil {
			return err
		}

		if i.isTruthy(condition) {
			// If yes, we evaluate the if/then block
			return i.execute(statement.ThenBranch)
		}

		// If no, we check if there exists an else branch
		if statement.ElseBranch != nil {
			// If there is an else branch we execute it
			return i.execute(statement.ElseBranch)
		}

	case *WhileSt:
		// Since a while statement is a loop, we loop
		// till the condition is no more truthy
		for {
			// Evaluating the statement condition
			statementCondition, err := i.evaluate(statement.Condition)
			if err != nil {
				return err
			}

			// Break if not truthy
			if !i.isTruthy(statementCondition) {
				break
			}

			// If not we execute the statements inside the block
			i.execute(statement.Body)
		}

	case *VariableSt:
		var variableValue any = nil

		// If the variable is initialized
		// We fetch the value
		if statement.Initializer != nil {
			value, err := i.evaluate(statement.Initializer)
			if err != nil {
				return err
			}

			variableValue = value
		}

		// The variableValue is then mapped with
		// The variable keyword/INITIALIZER
		i.Env.Define(statement.Name.Lexeme, variableValue)

	case *BlockSt:
		// A new sub-environment, linking the current env as parent
		// Is created (i.Env being parent, environment being the block env)
		environment := InitEnvironment(i.Env)
		err := i.executeBlock(statement.Statements, environment)
		if err != nil {
			return err
		}
	}

	return nil
}

// Evaluate handles all expression types
// 1. Literal
// 2. Grouping
// 3. Unary
// 4. Binary
// 5. Variable
// 6. Assignment
// 7. Logical (|| and &&)
// 8. Function call expression

func (i *Interpreter) evaluate(expression Expression) (any, error) {
	switch exp := expression.(type) {
	case *Literal:
		// A literal has already it's value present inside
		// Hence we can return just that.
		return exp.Value, nil

	case *Logical:
		// Similar to Binary, we evaluate the left expression
		left, err := i.evaluate(exp.Left)
		if err != nil {
			return nil, err
		}

		// For a Logical OR, if left is true, it's enough to evaluate the left
		// Due to short circuiting
		if exp.Operator.Type == s.Or && i.isTruthy(left) {
			return left, nil
		}

		// For a logical AND, if left is false, it's enough to evaluate the left
		// Due to short circuiting
		if exp.Operator.Type == s.And && !i.isTruthy(left) {
			return left, nil
		}

		// Else, we evaluate the right expression
		return i.evaluate(exp.Right)

	case *Grouping:
		// A parantheses can have multiple sub expressions inside
		// Hence we recurse the same function over and over
		// Till we evaluate all of it inside the paran
		return i.evaluate(exp.Expression)

	case *Unary:
		// Same goes with Unary.
		// Evaluate the expression first
		right, err := i.evaluate(exp.Right)
		if err != nil {
			return nil, err
		}

		// And then evaluate the operator before it
		// Which is either ! or -
		switch exp.Operator.Type {

		case s.Bang:
			return !i.isTruthy(right), nil
		case s.Minus:
			value, err := i.checkNumericOperand(exp.Operator, right)
			if err != nil {
				return nil, err
			}

			return -value, nil

		}

	case *Variable:
		// A variable expression needs the value
		// That is mapped to the environment
		return i.Env.Get(exp.Name)

	case *Assignment:
		// We first evaluate the value expression
		value, err := i.evaluate(exp.Value)
		if err != nil {
			return nil, err
		}

		// We then (re)assign the above evaluated value
		// With the variable name
		err = i.Env.Assign(exp.Name, value)
		if err != nil {
			return nil, err
		}

		return value, nil

	case *Call:
		// When a call expression is first met
		// FunctionName()
		// The interpreter wouldn't know what to do with that function
		// Hence we evaluate the FunctionSt Node to generate a function mapping with interpreter
		callee, err := i.evaluate(exp.Callee)
		if err != nil {
			return nil, err
		}

		// The arguments of each Call Statement is evaluated
		arguments := make([]any, 0)
		for _, argument := range exp.Arguments {
			evaluatedArg, err := i.evaluate(argument)
			if err != nil {
				return nil, err
			}

			arguments = append(arguments, evaluatedArg)
		}

		argumentsLen := len(arguments)
		// The callee expression is being checked
		// if it's of type Callable \
		function, ok := callee.(Callable)
		if !ok {
			return nil, fmt.Errorf("%s: Can only call functions and classes", exp.Paren)
		}

		// And the arguments received in the call statement
		// Is checked against the function declaration parameters size
		if argumentsLen != function.Arity() {
			return nil, fmt.Errorf("%s: Expected %d arguments, but got %d.", exp.Paren, function.Arity(), argumentsLen)
		}

		// The corresponding function is then called
		callResult := function.Call(i, arguments)

		// And fetches the result
		return callResult, nil

	case *Binary:
		// In unary we had just the right (since one)
		// Here we just repeat the same twice (two operands)
		left, err := i.evaluate(exp.Left)
		if err != nil {
			return nil, err
		}

		right, err := i.evaluate(exp.Right)
		if err != nil {
			return nil, err
		}

		switch exp.Operator.Type {
		// The plus symbol can be used for both
		case s.Plus:

			// 1. Numeric addition
			dValueL, okL := left.(float64)
			dValueR, okR := right.(float64)
			if okL && okR {
				return dValueL + dValueR, nil
			}

			// 2. String concatenation
			sValueL, okL := left.(string)
			sValueR, okR := right.(string)
			if okL && okR {
				return sValueL + sValueR, nil
			}

			// 3. Throw error otherwise
			return nil, fmt.Errorf("Operands must be numbers or two strings\n[line %d]", exp.Operator.LineNumber)

		// Every other math related operation
		// (-, *, /, <, <=, >, >=)
		// Has to check the if the operands are numbers
		// And then perform the evaluation accordingly
		case s.Minus:
			v, v1, err := i.checkNumericOperands(exp.Operator, left, right)
			if err != nil {
				return nil, err
			}

			return v - v1, nil

		case s.Slash:
			v, v1, err := i.checkNumericOperands(exp.Operator, left, right)
			if err != nil {
				return nil, err
			}

			return v / v1, nil

		case s.Star:
			v, v1, err := i.checkNumericOperands(exp.Operator, left, right)
			if err != nil {
				return nil, err
			}

			return v * v1, nil

		case s.Greater:
			v, v1, err := i.checkNumericOperands(exp.Operator, left, right)
			if err != nil {
				return nil, err
			}

			return v > v1, nil

		case s.GreaterEqual:
			v, v1, err := i.checkNumericOperands(exp.Operator, left, right)
			if err != nil {
				return nil, err
			}

			return v >= v1, nil

		case s.Less:
			v, v1, err := i.checkNumericOperands(exp.Operator, left, right)
			if err != nil {
				return nil, err
			}

			return v < v1, nil

		case s.LessEqual:
			v, v1, err := i.checkNumericOperands(exp.Operator, left, right)
			if err != nil {
				return nil, err
			}

			return v <= v1, nil

		// The rest is either ! or ==
		// Which is what we check if it's either equal or not
		case s.Bang:
			return !i.isEqual(left, right), nil
		case s.EqualEqual:
			return i.isEqual(left, right), nil
		}

	}

	// No other expression exists other than the above.
	return nil, nil
}

// -------------------------------------------------------------------

// ------------------------ UTILS ------------------------------------

func (i *Interpreter) executeBlock(statements []Statement, environment *Environment) error {
	// Keeping track of the parent environment
	previousEnvironment := i.Env
	// And switching temporarily the parent as the block env
	i.Env = environment

	// Once the environment is switched to block
	for _, st := range statements {
		// We execute the statements based on that block env
		err := i.execute(st)
		if err != nil {
			// We switch back env if the execution is stopped midway
			i.Env = previousEnvironment
			return err
		}
	}

	// And then switching back to the parent environment
	i.Env = previousEnvironment
	return nil
}

func (i *Interpreter) isTruthy(object any) bool {
	// We consider anything that's nil/false to be false
	if object == nil {
		return false
	}

	// Or, if the object is already of a bool type
	// We just return the value of it
	if value, ok := object.(bool); ok {
		return value
	}

	// Else, anything apart (empty string, 0, etc) from above two is true.
	return true
}

func (i *Interpreter) isEqual(aObj, bObj any) bool {
	// Two nils is going to be equal anyways
	if aObj == nil && bObj == nil {
		return true
	}

	// We check a being nil and b not being
	if aObj == nil {
		return false
	}

	// Else, we just try equalizing both and checking it
	return reflect.DeepEqual(aObj, bObj)
}

func (i *Interpreter) checkNumericOperand(operator s.Token, operand any) (float64, error) {
	// Check if it's actually a parsed number
	if value, ok := operand.(float64); ok {
		return value, nil
	}

	// Else return an error
	return 0, fmt.Errorf("Operand must be a number\n[line %d]", operator.LineNumber)
}

func (i *Interpreter) checkNumericOperands(operator s.Token, left, right any) (float64, float64, error) {
	// Check if both numbers are parsed ones (double/float64)
	v, ok := left.(float64)
	v1, ok2 := right.(float64)

	if ok && ok2 {
		return v, v1, nil
	}

	// Else return an error
	return 0, 0, fmt.Errorf("Operands must be numbers\n[line %d]", operator.LineNumber)
}

func (i *Interpreter) stringify(object any) string {
	// A nil value is simply nil in Lox
	if object == nil {
		return "nil"
	}

	// If the object is of a number
	if value, ok := object.(float64); ok {
		// We format it to string without losing the decimal precision
		text := strconv.FormatFloat(value, 'f', -1, 64)
		// And strip off the decimals
		if hasDecimal := strings.HasSuffix(text, ".0"); hasDecimal {
			text = text[0 : len(text)-2]
		}

		return text
	}

	// Else, we simply stringify the raw value
	return fmt.Sprintf("%v", object)
}

// -----------------------------------------------------------------------------
