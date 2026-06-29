package parser

import (
	s "github.com/sunguru98/glox/scanner"
)

// Operators are arithmetic/Logical symbols which are already defined as 'Token's

// An expression is any of the following
// 1. Literal (String, Number, Boolean, nil)
// 2. Unary (single operand)
// 3. Binary (operands on both sides with operators in the middle)
// 4. Grouping (Expression inside parantheses)

type Expression interface {
	// We create an empty function
	// So that, any struct can be an Expression through implementing this method
	exp()
}

// Following structs are based on the Expression interface
// Mentioned expression could be anything from the above
// Each struct implementing Expression will be Pointer type and not Value
// Since each expression node can be recursive, copying will be expensive

//----------------------------------------------------------------------------------------------------

// 1. Binary expression is of the grammar
// expression operator expression
type Binary struct {
	Left     Expression
	Operator s.Token
	Right    Expression
}

func (*Binary) exp() {}

func CreateBinaryExpression(left Expression, operator s.Token, right Expression) *Binary {
	return &Binary{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

//----------------------------------------------------------------------------------------------------

// 2. Grouping expression is of the grammar
// "(" expression ")"
type Grouping struct {
	Expression Expression
}

func (*Grouping) exp() {}

func CreateGroupingExpression(expression Expression) *Grouping {
	return &Grouping{
		Expression: expression,
	}
}

//----------------------------------------------------------------------------------------------------

// 3. Literal expression is of the grammar
// NUMBER / STRING / true / false / nil
type Literal struct {
	Value any
}

func (*Literal) exp() {}

func CreateLiteralExpression(value any) *Literal {
	return &Literal{Value: value}
}

//----------------------------------------------------------------------------------------------------

// 4. Unary expression is of the grammar
// ("-" or "!") expression
type Unary struct {
	Operator s.Token
	Right    Expression
}

func (*Unary) exp() {}

func CreateUnaryExpression(operator s.Token, right Expression) *Unary {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

//----------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------

// 5. Variable is of the grammar
// var IDENTIFER = (expression)? ';'
type Variable struct {
	Name  s.Token
	Right Expression
}

func (*Variable) exp() {}

func CreateVariableExpression(name s.Token) *Variable {
	return &Variable{
		Name: name,
	}
}

//----------------------------------------------------------------------------------------------------

// 6. Assignment is of the grammar
// (call '.')? IDENTIFER '=' (assignment | logic_or)
// An assignment expression can either set a variable or a member field in an instance (class object)
type Assignment struct {
	Name  s.Token
	Value Expression
}

func (*Assignment) exp() {}

func CreateAssignmentExpression(name s.Token, value Expression) *Assignment {
	return &Assignment{
		Name:  name,
		Value: value,
	}
}

//----------------------------------------------------------------------------------------------------

// 7. Logical is of the grammar
// OR - logic_and ("or" logic_and)*
// AND - equality ("and" equality)*

type Logical struct {
	Left     Expression
	Operator s.Token
	Right    Expression
}

func (*Logical) exp() {}

func CreateLogicalExpression(left Expression, operator s.Token, right Expression) *Logical {
	return &Logical{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

//----------------------------------------------------------------------------------------------------

// 8. Call is of the grammar
// primary '('arguments?' | '.' IDENTIFIER)'*
// A zero argument call can have the arguments optional
// Call supports both function calls, and Class instance creation
// The '.' operator is used to call the instance's property/methods

type Call struct {
	Paren     s.Token // Closing paren token needed for line number reporting
	Callee    Expression
	Arguments []Expression
}

func (*Call) exp() {}

func CreateCallExpression(callee Expression, paren s.Token, arguments []Expression) *Call {
	return &Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

//----------------------------------------------------------------------------------------------------

// 9. Get expressions base on the alternative grammar of Call (above)

type Get struct {
	Object Expression
	Name   s.Token
}

func (*Get) exp() {}

func CreateGetExpression(object Expression, name s.Token) *Get {
	return &Get{
		Object: object,
		Name:   name,
	}
}

//----------------------------------------------------------------------------------------------------

// 10. Set expressions base on the alternative grammar of Assignment (above)

type Set struct {
	Object Expression
	Name   s.Token
	Value  Expression
}

func (*Set) exp() {}

func CreateSetExpression(object Expression, name s.Token, value Expression) *Set {
	return &Set{
		Object: object,
		Name:   name,
		Value:  value,
	}
}

//----------------------------------------------------------------------------------------------------
