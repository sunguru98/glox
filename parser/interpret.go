package parser

import s "github.com/sunguru98/glox/scanner"

func isTruthy(object any) bool {
	// We consider anything that's nil/false to be false
	if object == nil {
		return false
	}

	if value, ok := object.(bool); ok {
		return value
	}

	// Rest being true
	return true
}

func isEqual(aObj, bObj any) bool {
	if aObj == nil && bObj == nil {
		return true
	}

	if aObj == nil {
		return false
	}

	return aObj == bObj
}

func Interpret(expression Expression) any {
	switch exp := expression.(type) {
	case *Literal:
		return exp.Value

	case *Grouping:
		return Interpret(exp.Expression)

	case *Unary:
		right := Interpret(exp.Right)
		switch exp.Operator.Type {
		case s.Minus:
			return -right.(float64)
		case s.Bang:
			return !isTruthy(right)
		}

		// Unreachable
		return nil

	case *Binary:
		left := Interpret(exp.Left)
		right := Interpret(exp.Right)

		switch exp.Operator.Type {
		case s.Plus:
			// The plus symbol can be used for both

			// 1. Numeric addition
			dValueL, okL := left.(float64)
			dValueR, okR := right.(float64)
			if okL && okR {
				return dValueL + dValueR
			}

			// 2. String concatenation
			sValueL, okL := left.(string)
			sValueR, okR := right.(string)
			if okL && okR {
				return sValueL + sValueR
			}

		case s.Minus:
			return left.(float64) - right.(float64)
		case s.Slash:
			return left.(float64) / right.(float64)
		case s.Star:
			return left.(float64) * right.(float64)
		case s.Greater:
			return left.(float64) > right.(float64)
		case s.GreaterEqual:
			return left.(float64) >= right.(float64)
		case s.Less:
			return left.(float64) < right.(float64)
		case s.LessEqual:
			return left.(float64) <= right.(float64)
		case s.Bang:
			return !isEqual(left, right)
		case s.EqualEqual:
			return isEqual(left, right)
		}

		// Unreachable
		return nil
	}

	// Unreachable
	return nil
}
