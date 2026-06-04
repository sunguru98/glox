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
}

func InitInterpreter() *Interpreter {
	return &Interpreter{}
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

	// Else, anything apart from above two is true.
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

func (i *Interpreter) evaluate(expression Expression) (any, error) {
	switch exp := expression.(type) {
	case *Literal:
		// A literal has already it's value present inside
		// Hence we can return just that.
		return exp.Value, nil

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

		// No other unary operator exists other than above mentioned
		return nil, nil

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

		// Unreachable
		return nil, nil
	}

	// No other expression exists other than the above.
	return nil, nil
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

func (i *Interpreter) Interpret(expression Expression) {
	// We evaluate the expression, and check for errors
	value, err := i.evaluate(expression)
	if err != nil {
		// If there is one, we report it as a runtime error
		lib.RuntimeError(err)
	}

	// Else we print out the evaluated expression's value
	fmt.Println(i.stringify(value))
}
