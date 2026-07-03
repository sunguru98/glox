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
// 7. Return statement

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

// A function statement is of grammar
// fun IDENTIFIER '('IDENTIFIER,(IDENTIFIER)*')'? block
type FunctionSt struct {
	Name   s.Token
	Params []s.Token
	Body   []Statement
}

func (*FunctionSt) stmnt() {}

func CreateFunctionSt(name s.Token, params []s.Token, body []Statement) *FunctionSt {
	return &FunctionSt{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

// -------------------------------------------------------------------------------------------------------------------------

// A return statement is of grammar
// Here expression is optional as not all functions return
// 'return' (expression)? ';'
type ReturnSt struct {
	Keyword s.Token // The return keyword is preserved for line number error reporting
	Value   Expression
}

func (*ReturnSt) stmnt() {}

func CreateReturnSt(keyword s.Token, value Expression) *ReturnSt {
	return &ReturnSt{
		Keyword: keyword,
		Value:   value,
	}
}

// -------------------------------------------------------------------------------------------------------------------------

// A class statement is of grammar
// 'class' IDENTIFIER ('<' IDENTIFIER)? '{ (function)* }'
// The function* here is simply stating, there can be more than one method
// Based on FunctionSt grammar above

type ClassSt struct {
	Name       s.Token
	Methods    []*FunctionSt
	SuperClass *Variable
}

func (*ClassSt) stmnt() {}

func CreateClassSt(name s.Token, methods []*FunctionSt, superClass *Variable) *ClassSt {
	return &ClassSt{
		Name:       name,
		Methods:    methods,
		SuperClass: superClass,
	}
}

// -------------------------------------------------------------------------------------------------------------------------
