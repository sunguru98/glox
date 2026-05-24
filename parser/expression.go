package parser

import s "github.com/sunguru98/glox/scanner"

// An expression is either a
// 1. Literal (String, Number, Boolean, nil)
// 2. Unary (single operand)
// 3. Binary (operands on both sides with operators in the middle)
// 4. Operators (Arithmetic/Logical symbols)
// 5. Grouping (Expression inside parantheses)
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
// NUMBER or STRING or true or false or nil
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
