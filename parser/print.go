package parser

import (
	"fmt"
	"strings"
)

// A Pretty printer function to showcase the evaluation of an expression
func Print(expression Expression) string {
	// We decide on the type of expression node and print accordingly
	switch exp := expression.(type) {
	case *Binary:
		// Operator symbol, Left operand, Right operand
		return parenthesize(exp.Operator.Lexeme, exp.Left, exp.Right)
	case *Unary:
		// Operator symbol (! or -), Right operand
		return parenthesize(exp.Operator.Lexeme, exp.Right)
	case *Grouping:
		// Since the printer already as "()", to differentiate it, the word "group" is added
		return parenthesize("group", exp.Expression)
	case *Literal:
		// If the underlying literal value is nothing, then it's mentioned as nil
		if exp.Value == nil {
			return "nil"
		}
		// Else the underlying literal value is stringified
		return fmt.Sprint(exp.Value)
	default:
		// AST should not reach here
		panic(fmt.Sprintf("Undefined AST Expression type %T", expression))
	}
}

func parenthesize(name string, expressions ...Expression) string {
	var stringBuilder strings.Builder
	// Prints the type of lexeme with open parentheses
	fmt.Fprintf(&stringBuilder, "(%s", name)

	// A recursive function where it prints for every sub expression
	for _, expression := range expressions {
		stringBuilder.WriteString(" ")
		// Which is then appended
		stringBuilder.WriteString(Print(expression))
	}

	// Ends with a closing parentheses and returns as a full string
	fmt.Fprintf(&stringBuilder, ")")
	return stringBuilder.String()
}
