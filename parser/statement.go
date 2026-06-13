package parser

import s "github.com/sunguru98/glox/scanner"

// -------------------------------------------------------------------------------------------------------------------------
// A statement can be of either
// 1. Print statement
// 2. Generic expression statement
// 3. Variable statement
// 4. Block statements
// 5. Conditional if statement
// 6. While/For statement
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

func CreateExpressionSt(expression Expression) *ExpressionSt {
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

func CreatePrintSt(expression Expression) *PrintSt {
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

func CreateVariableSt(name s.Token, init Expression) *VariableSt {
	return &VariableSt{
		Name:        name,
		Initializer: init,
	}
}

// -------------------------------------------------------------------------------------------------------------------------

// A block statement is of grammar
// '{' declaration '}'

type BlockSt struct {
	Statements []Statement
}

func (*BlockSt) stmnt() {}

func CreateBlockSt(statements []Statement) *BlockSt {
	return &BlockSt{
		Statements: statements,
	}
}

// -------------------------------------------------------------------------------------------------------------------------

// An if statement is of grammar
// 'if' '(' conditional_expression ')' then_statement
// ('else' else_statement)?
// ()? means optional
type IfSt struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (*IfSt) stmnt() {}

func CreateIfSt(condition Expression, thenB, elseB Statement) *IfSt {
	return &IfSt{
		Condition:  condition,
		ThenBranch: thenB,
		ElseBranch: elseB,
	}
}

// -------------------------------------------------------------------------------------------------------------------------

// A while statement is of grammar
// 'while' '(' expression ')' statement

// Since a for loop is just syntactic sugar for while
// Its grammar alone changes
// 'for' '(' (variableDeclaration | expression | ';') expression?';' expression? ')'
// statement;

type WhileSt struct {
	Condition Expression
	Body      Statement
}

func (*WhileSt) stmnt() {}

func CreateWhileSt(condition Expression, body Statement) *WhileSt {
	return &WhileSt{
		Condition: condition,
		Body:      body,
	}
}

// -------------------------------------------------------------------------------------------------------------------------
