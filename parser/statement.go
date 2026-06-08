package parser

import s "github.com/sunguru98/glox/scanner"

// -------------------------------------------------------------------------------------------------------------------------
// A statement can be of either
// 1. Print statement
// 2. Generic expression statement
// 3. Variable statement
type Statement interface {
	stmnt()
}

// -------------------------------------------------------------------------------------------------------------------------

// A generic expression statement is of grammar
// expression ';'
// Basically any expression that ends with a semicolon

type ExpressionSt struct {
	Expr Expression
}

func (*ExpressionSt) stmnt() {}

func CreateNewExpressionSt(expression Expression) *ExpressionSt {
	return &ExpressionSt{
		Expr: expression,
	}
}

// -------------------------------------------------------------------------------------------------------------------------

// A print statement is of grammar
// print expression ';'

type PrintSt struct {
	Expr Expression
}

func (*PrintSt) stmnt() {}

func CreateNewPrintSt(expression Expression) *PrintSt {
	return &PrintSt{
		Expr: expression,
	}
}

// -------------------------------------------------------------------------------------------------------------------------

// A variable statememnt is of grammar
// "var" IDENTIFIER/NAME ('=' expression)? ';'
// The parantheses with a question mark basically means optional
// In other words, a variable can be declared without initialization

type VariableSt struct {
	Name        s.Token
	Initializer Expression
}

func (*VariableSt) stmnt() {}

func CreateNewVariableSt(name s.Token, init Expression) *VariableSt {
	return &VariableSt{
		Name:        name,
		Initializer: init,
	}
}

// -------------------------------------------------------------------------------------------------------------------------
